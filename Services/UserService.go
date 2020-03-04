package Services

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"pro2/Baseinfo"
	"time"
)

//------------------登陆------------------

type WUserLoginService interface {
	Login(userid string, pwd string) *CommonResponse
}

type UserLoginService struct{}

func (this UserLoginService) Login(userId string, pwd string) *CommonResponse {
	token, errCode, err := Baseinfo.LoginAuth(userId, pwd)
	if err != nil {
		_ = logger.Log("LoginError:", err)
	}
	return &CommonResponse{
		Code:   errCode,
		Msg:    err.Error(),
		Data:   token,
		Expand: nil,
	}
}

//----创建新用户----

type WUserCreateService interface {
	NewAccount(r *UserCreateRequest) *CommonResponse
}

type UserCreateService struct{}

func (this UserCreateService) NewAccount(r *UserCreateRequest) *CommonResponse {
	var findUser *Baseinfo.User
	response := &CommonResponse{}
	colUser := Baseinfo.Client.Database("isms").Collection("user")

	ErrFind := colUser.FindOne(context.Background(), bson.M{"userid": r.UserId}).Decode(&findUser)
	if ErrFind != nil && ErrFind.Error() != "mongo: no documents in result" {
		response.Code = Baseinfo.CONST_FIND_FAIL
		response.Msg = ErrFind.Error()
		_ = logger.Log("Create_User_Err:", "find user by userid in mongo fails", "mongo error:", ErrFind)
		return response
	}
	if findUser != nil { //没有报错且找到了该用户 说明已经被注册
		response.Code = Baseinfo.CONST_USER_OCCUPIED
		response.Msg = "this userId has been occupied,pls use another one!"
		_ = logger.Log("CreateUserError:", "this userId has been occupied,pls use another one!")
		return response
	}

	newUser := &Baseinfo.User{
		Id:       primitive.NewObjectIDFromTimestamp(time.Now()),
		Userid:   r.UserId,
		Level:    0, //目前暂时不在乎权限。以后默认新的用户均为0权限-即普通用户
		Password: r.Password,
		Phone:    r.Phone,
		Title:    r.Title,
		Nickname: r.Nickname,
		Email:    r.Email,
		External: nil,
	}
	ctx := context.Background()

	var newUserInfo *Baseinfo.User
	SessionErr := Baseinfo.Client.Database("isms").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		InsResult, ErrInsert := colUser.InsertOne(sessionContext, &newUser)
		if ErrInsert != nil {
			response.Code = Baseinfo.CONST_INSERT_FAIL
			response.Msg = ErrInsert.Error()
			_ = logger.Log("CreateUserError:", ErrInsert.Error())
			return ErrInsert
		}

		ErrF := colUser.FindOne(sessionContext, bson.D{{"_id", InsResult.InsertedID.(primitive.ObjectID)}}).Decode(&newUserInfo)
		if ErrF != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			_ = logger.Log("CreateUserError:", "can't find recently created user!"+ErrF.Error())
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = ErrF.Error()
			return ErrF
		} else {
			_ = sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})

	if SessionErr != nil {
		_ = logger.Log("Create_User_Err:", SessionErr)
		return response
	}
	response.Code = Baseinfo.Success
	response.Data = newUserInfo
	return response
}

//------------------获取用户列表------------------
type WUserListService interface {
	ObtainUserList(user string) *CommonResponse
}

type UserListService struct{}

func (this UserListService) ObtainUserList(user string) *CommonResponse {
	response := &CommonResponse{}
	colUser := Baseinfo.Client.Database("isms").Collection("user")

	var user0 *Baseinfo.User
	err := colUser.FindOne(context.Background(), bson.M{"userid": user}).Decode(&user0)
	if user0 == nil {
		response.Code = Baseinfo.CONST_FIND_FAIL
		response.Msg = "unauthorized user"
		return response
	}
	if user0.Level != 0 {
		//非管理员账户
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "only admin user can obtain user lists"
		return response
	}

	cur, err := colUser.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		response.Code = Baseinfo.CONST_FIND_FAIL
		response.Msg = "can't find data"
		return response
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		response.Code = Baseinfo.CONST_FIND_FAIL
		response.Msg = "can't find data"
		return response
	}
	var all []*Baseinfo.User
	err = cur.All(context.Background(), &all)
	if err != nil {
		log.Println(err)
		response.Code = 500
		response.Msg = "can't find data"
		return response
	}
	_ = cur.Close(context.Background())

	//log.Println("collection.Find curl.All: ", all)
	//for _, one := range all {
	//	log.Println(one)
	//}
	response.Code = 200
	response.Data = all

	return response
}

//------------------编辑用户------------------
type WUserEditService interface {
	UserEdit(r *UserEditRequest, token string) *CommonResponse
}

type UserEditService struct{}

func (this UserEditService) UserEdit(r *UserEditRequest, token string) *CommonResponse {
	response := &CommonResponse{}
	colUser := Baseinfo.Client.Database("isms").Collection("user")

	ok, err, user := Baseinfo.TokenCheckAsymmetricalKey(token)
	if !ok || err != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err.Error()
		return response
	}

	var user0 *Baseinfo.User
	err = colUser.FindOne(context.Background(), bson.M{"userid": user}).Decode(&user0)
	if user0 == nil {
		response.Code = Baseinfo.CONST_USER_NOTEXIST
		response.Msg = "登陆用户不存在"
		return response
	}
	if user0.Level != 0 && user != r.UserId {
		//不是管理员用户，但是要修改信息的账户不是自己的账户--不允许
		response.Code = Baseinfo.CONST_ACTION_UNALLOWED
		response.Msg = "不允许修改其他用户信息"
		return response
	}
	ctx := context.Background()
	var user1 *Baseinfo.User
	var revisedUser1 *Baseinfo.User
	sessionErr := Baseinfo.Client.Database("isms").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err = colUser.FindOne(sessionContext, bson.D{{"userid", r.UserId}}).Decode(&user1)
		if user1 == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "user not exists"
			_ = logger.Log("EdituserError:", err.Error())
			return errors.New("user not exists")
		}

		_, err = colUser.UpdateOne(sessionContext, bson.D{{"userid", r.UserId}}, bson.D{{"$set", bson.D{
			{"phone", r.Phone},
			{"title", r.Title},
			{"nickname", r.Nickname},
			{"email", r.Email},
		}}})
		if err != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = "fail to update user"
			_ = logger.Log("ReviseuserError:", err.Error())
			return errors.New("fail to update user")
		}
		err = colUser.FindOne(sessionContext, bson.D{{"userid", r.UserId}}).Decode(&revisedUser1)
		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "can't find recently revised device"
			return errors.New("can't find recently revised device")
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})
	if sessionErr != nil {
		_ = logger.Log("ReviseUserError:", sessionErr)
		return response
	}
	return &CommonResponse{
		Code:   200,
		Msg:    "",
		Data:   revisedUser1,
		Expand: nil,
	}
}

//------------------删除用户------------------
type WUserDelService interface {
	UserDel(userid, token string) *CommonResponse
}

type UserDelService struct{}

func (this UserDelService) UserDel(userid, token string) *CommonResponse {
	response := &CommonResponse{}
	colUser := Baseinfo.Client.Database("isms").Collection("user")

	ok, err, user := Baseinfo.TokenCheckAsymmetricalKey(token)
	if !ok || err != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err.Error()
		return response
	}

	var user0 *Baseinfo.User
	err = colUser.FindOne(context.Background(), bson.M{"userid": user}).Decode(&user0)
	if user0 == nil {
		response.Code = Baseinfo.CONST_USER_NOTEXIST
		response.Msg = "登陆用户不存在"
		return response
	}
	if user0.Level != 0 {
		//非管理员账户  ---无权删除用户
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "normal user is not allowed to delete account"
		return response
	}
	var user1 *Baseinfo.User
	err = colUser.FindOne(context.Background(), bson.M{"userid": userid}).Decode(&user1)
	if user1 == nil {
		response.Code = Baseinfo.CONST_USER_NOTEXIST
		response.Msg = "account to delete not exist"
		return response
	}
	if user1.Level == 0 {
		//要删除的账户是管理员账户  ---不能删除管理员账户
		response.Code = Baseinfo.CONST_ACTION_UNALLOWED
		response.Msg = "admin account mustn't be deleted"
		return response
	}
	_, err = colUser.DeleteOne(context.Background(), bson.M{"userid": userid})
	if err != nil {
		response.Code = Baseinfo.CONST_DELETE_FAIL
		response.Msg = "删除用户失败"
		_ = logger.Log("DelUserError:", err)
		return response
	}

	return &CommonResponse{
		Code:   200,
		Msg:    "",
		Data:   nil,
		Expand: nil,
	}
}

//------------------修改用户密码------------------
type WUserResetService interface {
	UserReset(r *UserResetRequest, token string) *CommonResponse
}

type UserResetService struct{}

func (this UserResetService) UserReset(r *UserResetRequest, token string) *CommonResponse {
	response := &CommonResponse{}
	colUser := Baseinfo.Client.Database("isms").Collection("user")

	ok, err, user := Baseinfo.TokenCheckAsymmetricalKey(token)
	if !ok || err != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err.Error()
		return response
	}
	if user != r.UserId {
		//登陆用户和要修改密码的用户不一致
		response.Code = Baseinfo.CONST_ACTION_UNALLOWED
		response.Msg = "can't revise other user's password"
		return response
	}

	var user0 *Baseinfo.User
	var revisedUser0 *Baseinfo.User
	sessionErr := Baseinfo.Client.Database("isms").Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			_ = logger.Log("EditUserPwdError:", err.Error())
			return err
		}
		err = colUser.FindOne(sessionContext, bson.D{{"userid", r.UserId}}).Decode(&user0)
		if user0 == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "user not exists"
			_ = logger.Log("EditUserError:", err.Error())
			return errors.New("user not exists")
		}

		if user0.Password == r.NewPassword {
			return errors.New("new password can't be the same to the previous one")
		}

		_, err = colUser.UpdateOne(sessionContext, bson.D{{"userid", r.UserId}}, bson.D{{"$set", bson.D{
			{"password", r.NewPassword},
		}}})
		if err != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = "fail to update password"
			_ = logger.Log("ReviseUserError:", err.Error())
			return errors.New("fail to update password")
		}
		err = colUser.FindOne(sessionContext, bson.D{{"userid", r.UserId}}).Decode(&revisedUser0)
		if err != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "can't find recently revised user"
			return errors.New("can't find recently revised user")
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})
	if sessionErr != nil {
		_ = logger.Log("ReviseUserPwdError:", sessionErr)
		return response
	}
	return &CommonResponse{
		Code:   200,
		Msg:    "",
		Data:   revisedUser0,
		Expand: nil,
	}
}

//------------------用户注销------------------
type WUserLogoutService interface {
	UserLogout(userid string) *CommonResponse
}

type UserLogoutService struct{}

func (this UserLogoutService) UserLogout(userid string) *CommonResponse {
	fmt.Println("进入用户注销接口", userid)
	return &CommonResponse{
		Code:   200,
		Msg:    "",
		Data:   nil,
		Expand: nil,
	}
}
