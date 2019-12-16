package Baseinfo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

//设备信息
type Device struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Sid       primitive.ObjectID //objectid
	Userid    string
	Gatewayid string //产品所属网关的编号
	Deviceid  string //产品本code，第三方设备传感器 主要是面向运维或者非技术人员查看的编号
	//产品所属网关的编号 如果设备本身即有联网功能 则自己就是自己的网关
	Isnode  bool  //显示设备本身为节点还是网关
	Devtype int64 //不同类型的设备（传感器），有不同的devtype 为方便查看，可以直接使用二进制数的不同位来表示
	//如 0001表示摄像头，0010表示人体运动传感器，0100表示五合一传感器 但是在数据库中用十进制的结
	// 果表示，如1，2，4，8...
	Title string //TODO 表示设备的详细信息，如具体的位置
	Addr  string //表示设备的地址 xx省xx市xx区xx商圈（xx小区，xx园区）+（title内容（可以在存储时把title的信
	// 息也一起补充进来））
	Spacecode string      //地址的编码 按照层级设计过 每个代码对应唯一地址
	Expand    interface{} //存储传感器最近上报的一次数据
	External  interface{} //需要时用来补充信息
}

//解绑设备方法
func UnboundDeviceBydeviceid(dev *Device, kind string, sessionContext mongo.SessionContext, col_dev, col_space *mongo.Collection) (errcode int64, errmsg string, data interface{}) {
	filter := bson.D{{"deviceid", dev.Deviceid}}
	var update bsonx.Doc

	switch kind {
	case "0":
		update = bsonx.Doc{{"$set", bsonx.Document(
			bsonx.Doc{{"gatewayid", bsonx.String("")},
				{"userid", bsonx.String("")},
				{"sid", bsonx.ObjectID(primitive.NilObjectID)},
				{"addr", bsonx.String("")},
				{"spacecode", bsonx.String("")},
			})}}
		if dev.Sid != primitive.NilObjectID {
			RemoveDev(dev.Sid, dev.Id, sessionContext, col_space)
		}

	case "1": //解绑账户关联
		update = bsonx.Doc{{"$set", bsonx.Document(
			bsonx.Doc{{"userid", bsonx.String("")}})}}
	case "2": //解绑网关管理
		update = bsonx.Doc{{"$set", bsonx.Document(
			bsonx.Doc{{"gatewayid", bsonx.String("")}})}}
	case "3": //解绑地址
		update = bsonx.Doc{{"$set", bsonx.Document(
			bsonx.Doc{
				{"sid", bsonx.ObjectID(primitive.NilObjectID)},
				{"addr", bsonx.String("")},
				{"spacecode", bsonx.String("")},
			})}}
		if dev.Sid != primitive.NilObjectID {
			RemoveDev(dev.Sid, dev.Id, sessionContext, col_space)
		}
	default:
		return CONST_PARAM_ERROR, "operate type is out of range", nil
	}
	_, err := col_dev.UpdateOne(sessionContext, filter, update)
	if err != nil {
		return Fail, err.Error(), nil
	}
	var device *Device
	_ = col_dev.FindOne(context.Background(), bson.D{{"deviceid", dev.Deviceid}}).Decode(&device)
	return Success, "", device
}

func UnboundDeviceByid(dev *Device, kind string, sessionContext mongo.SessionContext, col_dev, col_space *mongo.Collection) (errcode int64, errmsg string, data interface{}) {
	filter := bson.D{{"_id", dev.Id}}
	var update bsonx.Doc

	switch kind {
	case "0":
		update = bsonx.Doc{{"$set", bsonx.Document(
			bsonx.Doc{{"gatewayid", bsonx.String("")},
				{"userid", bsonx.String("")},
				{"sid", bsonx.ObjectID(primitive.NilObjectID)},
				{"addr", bsonx.String("")},
				{"spacecode", bsonx.String("")},
			})}}
		if dev.Sid != primitive.NilObjectID {
			RemoveDev(dev.Sid, dev.Id, sessionContext, col_space)
		}
	case "1":
		update = bsonx.Doc{{"$set", bsonx.Document(
			bsonx.Doc{{"userid", bsonx.String("")}})}}
	case "2":
		update = bsonx.Doc{{"$set", bsonx.Document(
			bsonx.Doc{{"gatewayid", bsonx.String("")}})}}
	case "3":
		update = bsonx.Doc{{"$set", bsonx.Document(
			bsonx.Doc{
				{"sid", bsonx.ObjectID(primitive.NilObjectID)},
				{"addr", bsonx.String("")},
				{"spacecode", bsonx.String("")},
			})}}
		if dev.Sid != primitive.NilObjectID {
			RemoveDev(dev.Sid, dev.Id, sessionContext, col_space)
		}
	default:
		return CONST_PARAM_ERROR, "operatie type is out of range", nil
	}
	_, err := col_dev.UpdateOne(sessionContext, filter, update)
	if err != nil {
		return Fail, err.Error(), nil
	}
	var device *Device
	_ = col_dev.FindOne(sessionContext, bson.D{{"deviceid", dev.Deviceid}}).Decode(&device)
	return Success, "", dev
}

func Unbounduser() {

}

func Unboundaddr() {

}

func Unboundccode() {

}
