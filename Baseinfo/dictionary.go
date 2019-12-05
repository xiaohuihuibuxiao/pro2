package Baseinfo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dictionary struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Code       string             //6位十进制
	Name       string
	Parentid   string
	Shortname  string
	Leveltype  int //0-china 1-省（直辖市，北京） 2-市（直辖市自己，北京市） 3-区（西湖区，滨江区）
	Citycode   string
	Zipcode    string
	Mergername string
	Ing        float64 //经度
	Lat        float64 //纬度
	Pinnyin    string
}
