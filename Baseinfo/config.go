package Baseinfo

import (
	"context"
	"encoding/json"
	"fmt"
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
}

var Client *mongo.Client

func init() {
	v := &Config{}
	//data,err:=ioutil.ReadFile("./config.json")//放在服务器上时 这么写路径
	data, err := ioutil.ReadFile("D:/common_pro/src/pro2/Baseinfo/config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	err_j := json.Unmarshal(data, &v)
	if err_j != nil {
		log.Fatal(err_j)
		return
	}

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
	//	fmt.Println("aaaa", err)
	//	return
	//}

	if err := Client.Ping(context.Background(), readpref.Primary()); err != nil {
		fmt.Println("", err)
		return
	}
	log.Println("connect successfully to mongodb ...")

}