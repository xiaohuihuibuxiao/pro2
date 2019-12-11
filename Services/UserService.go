package Services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		logger.Log("Login err:", err.Error())
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

	err_find := col_user.FindOne(context.Background(), bson.M{"userid": r.Userid}).Decode(&finduser)
	if err_find != nil && err_find.Error() != "mongo: no documents in result" {
		commonresponse.Code = Baseinfo.CONST_FIND_FAIL
		commonresponse.Msg = err_find.Error()
		logger.Log("Create_User_Err:", "find user by userid in mongo fails", "mongo error:", err_find)
		return commonresponse
	}
	if finduser != nil { //没有报错且找到了该用户 说明已经被注册了
		commonresponse.Code = Baseinfo.CONST_USER_OCCUPIED
		commonresponse.Msg = "this userid has been occupied,pls use another one!"
		logger.Log("Create_User_Err:", "this userid has been occupied,pls use another one!")
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
	ins_result, err_insert := col_user.InsertOne(context.Background(), &newuser)
	if err_insert != nil {
		commonresponse.Code = Baseinfo.CONST_INSERT_FAIL
		commonresponse.Msg = err_insert.Error()
		logger.Log("Create_User_Err:", err_insert.Error())
		return commonresponse
	}
	var newuserinfo *Baseinfo.User
	col_user.FindOne(context.Background(), bson.D{{"_id", ins_result.InsertedID.(primitive.ObjectID)}}).Decode(&newuserinfo)
	commonresponse.Code = Baseinfo.Success
	commonresponse.Data = newuserinfo
	return commonresponse
}
