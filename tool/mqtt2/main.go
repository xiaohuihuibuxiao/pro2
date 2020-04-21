package main

import (
	"encoding/json"
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"

	//"log"
	"os"
	"os/signal"
	"syscall"
)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	var a struct {
		DeviceId string `json:"deviceId"`
		Status   string `json:"status"`
		Service  string `json:"service"`
		T        int64  `json:"t"`
		Num      int    `json:"num"`
	}
	_ = json.Unmarshal(message.Payload(), &a)
	fmt.Printf("%+v\n", a)
}

func main() {
	//MQTT.DEBUG = log.New(os.Stdout, "", 0)
	//MQTT.ERROR = log.New(os.Stdout, "", 0)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	//hostname, _ := os.Hostname()

	server := flag.String("server", "tcp://47.100.44.103:1883", "The full url of the MQTT server to connect to ex: tcp://xxxx:1883")
	topic := flag.String("topic", "test", "Topic to subscribe to")
	qos := flag.Int("qos", 2, "The QoS to subscribe to messages at")
	clientid := flag.String("clientid", "MQTTClientSub", "A clientid for the connection")
	//clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	//username := flag.String("username", "wlh", "A username to authenticate to the MQTT server")
	//	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	connOpts.SetKeepAlive(2 * time.Second)

	connOpts.SetUsername("wlh")
	//tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	//connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		fmt.Println("onconnect0000000000000000000")
		if token := c.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", *server)
	}
	<-c
}
