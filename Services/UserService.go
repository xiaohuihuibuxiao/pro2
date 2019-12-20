package Services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"pro2/Baseinfo"
	"time"
)

//------------------登陆------------------

type WUserLoginService interface {
	Login(userid string, pwd string) *CommonResponse
}

type UserLoginService struct{}

func (this UserLoginService) Login(userid string, pwd string) *CommonResponse {
	token, errcode, err := Baseinfo.Loginauth(userid, pwd)
	if err != nil {
		_ = logger.Log("Login_Err:", err)
	}

	return &CommonResponse{
		Code:   errcode,
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
	var finduser *Baseinfo.User
	commonresponse := &CommonResponse{}
	col_user := Baseinfo.Client.Database("test").Collection("user")

	ErrFind := col_user.FindOne(context.Background(), bson.M{"userid": r.Userid}).Decode(&finduser)
	if ErrFind != nil && ErrFind.Error() != "mongo: no documents in result" {
		commonresponse.Code = Baseinfo.CONST_FIND_FAIL
		commonresponse.Msg = ErrFind.Error()
		_ = logger.Log("Create_User_Err:", "find user by userid in mongo fails", "mongo error:", ErrFind)
		return commonresponse
	}
	if finduser != nil { //没有报错且找到了该用户 说明已经被注册了
		commonresponse.Code = Baseinfo.CONST_USER_OCCUPIED
		commonresponse.Msg = "this userid has been occupied,pls use another one!"
		_ = logger.Log("Create_User_Err:", "this userid has been occupied,pls use another one!")
		return commonresponse
	}

	newuser := &Baseinfo.User{
		Id:       primitive.NewObjectIDFromTimestamp(time.Now()),
		Userid:   r.Userid,
		Level:    0, //目前暂时不在乎权限。以后默认新的用户均为0权限-即普通用户
		Password: r.Password,
		Phone:    r.Phone,
		Title:    r.Title,
		Nickname: r.Nickname,
		Email:    r.Email,
		External: nil,
	}
	ctx := context.Background()

	var newuserinfo *Baseinfo.User
	SessionErr := Baseinfo.Client.Database("test").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		InsResult, ErrInsert := col_user.InsertOne(sessionContext, &newuser)
		if ErrInsert != nil {
			commonresponse.Code = Baseinfo.CONST_INSERT_FAIL
			commonresponse.Msg = ErrInsert.Error()
			_ = logger.Log("CreateUserErr:", ErrInsert.Error())
			return ErrInsert
		}

		ErrF := col_user.FindOne(sessionContext, bson.D{{"_id", InsResult.InsertedID.(primitive.ObjectID)}}).Decode(&newuserinfo)
		if ErrF != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			_ = logger.Log("CreateUserErr:", "can't find recently created user!"+ErrF.Error())
			commonresponse.Code = Baseinfo.CONST_FIND_FAIL
			commonresponse.Msg = ErrF.Error()
			return ErrF
		} else {
			_ = sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})

	if SessionErr != nil {
		_ = logger.Log("Create_User_Err:", SessionErr)
		return commonresponse
	}
	commonresponse.Code = Baseinfo.Success
	commonresponse.Data = newuserinfo
	return commonresponse
}
