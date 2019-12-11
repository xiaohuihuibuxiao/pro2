package Baseinfo

import (
	"context"
	"encoding/json"
	"flag"
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
	Expiredtime int
}

var Client *mongo.Client
var Expiredtime int //单位为秒

func init() {
	v := &Config{}

	road := flag.String("conf", "./config.json", "config road")

	flag.Parse()
	fmt.Println("port:", *road)

	var configdata []byte
	data0, err0 := ioutil.ReadFile(*road)
	if data0 == nil || err0 != nil {
		data1, err1 := ioutil.ReadFile("./config.json")
		if err1 != nil || data1 == nil {
			log.Fatal("read default config err:", err1)
			return
		}
		configdata = data1
	} else {
		configdata = data0
	}

	err_j := json.Unmarshal(configdata, &v)
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
