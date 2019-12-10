package Baseinfo

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
)

type Config struct {
	Addr  string
	Mongo struct {
		Authsource string
		Username   string
		Password   string
	}
	Expiredtime int
}

var Client *mongo.Client
var Expiredtime int //单位为秒

func init() {
	v := &Config{}
	//data,err:=ioutil.ReadFile("./config.json")//放在服务器上时 这么写路径
	//data, err := ioutil.ReadFile("D:/code/xatt/pro2/Baseinfo/config.json")//家里的路径
	data, err := ioutil.ReadFile("./Baseinfo/config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	err_j := json.Unmarshal(data, &v)
	if err_j != nil {
		log.Fatal(err_j)
		return
	}
	Expiredtime = v.Expiredtime
	opts := options.Client().ApplyURI("mongodb://47.100.44.103:27017")
	opts.SetAuth(options.Credential{
		AuthMechanism:           "SCRAM-SHA-1",
		AuthMechanismProperties: nil,
		AuthSource:              v.Mongo.Authsource,
		Username:                v.Mongo.Username,
		Password:                v.Mongo.Password,
		//	PasswordSet:             false,
	})

	Client, _ = mongo.Connect(context.Background(), opts)
	//Client, err := mongo.Connect(context.Background(), opts)
	//if err != nil {
	//	return
	//}

	if err := Client.Ping(context.Background(), readpref.Primary()); err != nil {
		return
	}
	log.Println("connect successfully to mongodb ...")

}
