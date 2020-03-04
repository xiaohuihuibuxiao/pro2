package Baseinfo

import (
	"context"
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var Client *mongo.Client

func init() {
	var path string
	flag.StringVar(&path, "c", "D:/wlh/project/pro2/pro2/Baseinfo/config.json", "conf file")
	flag.Parse()

	LoadConf(path)
	opts := options.Client().ApplyURI(Conf.Mgoaddr)
	opts.SetAuth(options.Credential{
		AuthMechanism:           "SCRAM-SHA-1",
		AuthMechanismProperties: nil,
		AuthSource:              Conf.authsource,
		Username:                Conf.Username,
		Password:                Conf.Password,
		PasswordSet:             false,
	})

	Client, _ = mongo.Connect(context.Background(), opts)

	if err := Client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatal("fail to connect to mongo: ", err)
		return
	}
	log.Println("connect successfully to mongodb ...")
}

var Conf *Config

// Config 配置参数
type Config struct {
	Mgoaddr    string
	authsource string
	Username   string
	Password   string
}

func read(v *viper.Viper, path string) {
	v.SetConfigFile(path)
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalln("viper ReadInConfig error: ", err)
	}

	config := &Config{}

	config.Mgoaddr = v.GetString("mgoaddr")
	config.authsource = v.GetString("authsource")
	config.Username = v.GetString("username")
	config.Password = v.GetString("password")

	Conf = config
}

func LoadConf(path string) {
	v := viper.New()
	read(v, path)

	// monitor file change
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config file changed:", e.Name)
		read(v, path)
		log.Println("config file updated:", Conf)
	})
	v.WatchConfig()
}
