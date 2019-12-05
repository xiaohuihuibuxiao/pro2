package Baseinfo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Sensorhistory struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Userid   string
	Pcode    string
	T        string
	Devtype  int
	Expand   interface{} //存储不同类型的传感器的数据
	External interface{} //扩展字段

}
