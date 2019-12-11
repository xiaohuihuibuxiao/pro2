package Baseinfo

import (
	"context"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Logger struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Time     string
	Url      string
	User     string
	Method   string
	External interface{} //扩展字段
	Expand   interface{}
}

//操作日志不记录get方法的请求,其余每个接口调用后均调用此函数往logger表存储该操作记录
func RecordOperation(url, method, user string) {
	if method != "GET" {
		col_log := Client.Database("test").Collection("logger")
		log := Logger{
			Id:     primitive.NewObjectIDFromTimestamp(time.Now()),
			Time:   time.Now().Format("2006-01-02 15:04:05"),
			Url:    url,
			User:   user,
			Method: method,
		}
		insertone, _ := col_log.InsertOne(context.Background(), log)
		logs.Info("detecting one operating log", insertone.InsertedID, method, user, url)
	}

}
