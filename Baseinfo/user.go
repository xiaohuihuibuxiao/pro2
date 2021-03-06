package Baseinfo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Userid   string
	Level    int64 //0-普通 1 管理员 2-super管理员/
	Password string
	Phone    string
	Title    string
	Nickname string
	Email    string
	External interface{}
}

func LoginAuth(user, pwd string) (string, int64, error) {
	if user == "" {
		return "", CONST_PARAM_LACK, errors.New("userId can't be blank")
	}
	colUser := Client.Database("isms").Collection("user")
	var userInfo *User
	err := colUser.FindOne(context.Background(), bson.M{"userid": user}).Decode(&userInfo)
	if err != nil {
		return "", CONST_USER_NOTEXIST, err
	}
	if userInfo.Password != pwd {
		return "", CONST_USERPWD_UNMATCH, errors.New("please check your userid or password !")
	}
	//TODO 需要用新的生成topken函数
	//token, err_token := Gentoken(user)
	token, err_token := TokenGen_asymmetricalkey(user)
	if err_token != nil {
		return "", CONST_TOEKN_ERROR, err_token
	}
	return token, Success, errors.New("")
}
