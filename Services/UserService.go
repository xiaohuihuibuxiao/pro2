package Services

import (
	"context"
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
	fmt.Println("登陆service")
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
	UserEdit(r *UserEditRequest) *CommonResponse
}

type UserEditService struct{}

func (this UserEditService) UserEdit(r *UserEditRequest) *CommonResponse {
	fmt.Println("进入编辑用户接口", r)
	return &CommonResponse{
		Code:   200,
		Msg:    "",
		Data:   nil,
		Expand: nil,
	}
}

//------------------删除用户------------------
type WUserDelService interface {
	UserDel(userid string) *CommonResponse
}

type UserDelService struct{}

func (this UserDelService) UserDel(userid string) *CommonResponse {
	fmt.Println("进入删除接口", userid)
	return &CommonResponse{
		Code:   200,
		Msg:    "",
		Data:   nil,
		Expand: nil,
	}
}

//------------------修改用户密码------------------
type WUserResetService interface {
	UserReset(r *UserResetRequest) *CommonResponse
}

type UserResetService struct{}

func (this UserResetService) UserReset(r *UserResetRequest) *CommonResponse {
	fmt.Println("进入修改密码接口", r)
	return &CommonResponse{
		Code:   200,
		Msg:    "",
		Data:   nil,
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
