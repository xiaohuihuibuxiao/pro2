package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

type msg struct {
	Messagetype int
	Data        []byte
}

//var addr = flag.String("addr", "47.100.44.103:14141", "http service address")
//var addr = flag.String("addr", "192.168.9.63:9090", "http service address")
//var addr = flag.String("addr", "localhost:9090", "http service address")
//var addr = flag.String("addr", "172.18.0.216:9090", "http service address")
var addr = flag.String("addr", "wufangjun.51vip.biz:8000", "http service address")

//var addr = flag.String("addr", "47.100.44.103:9090", "http service address")

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   *addr,
		//Path:   "/ismc/v1/pad/ws",
		Path: "/ismc/v1/homepage/ws",
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("aaaaaaa")
		fmt.Println("err:", err.Error())
		return
	}

	fmt.Println("成功建立连接")
	//go keepAlive(conn)
	//go Send(conn)
	go Send2(conn) //发送心跳
	Read(conn)
}

func keepAlive(ws *websocket.Conn) {
	time.Sleep(2 * time.Second)
	data, _ := json.Marshal("ping")
	for {
		err := ws.WriteMessage(2, data)
		if err != nil {
			log.Println("发送Ping错误", err.Error())
			break
		}
		fmt.Println("发送ping消息", string(data))
		<-time.NewTicker(20 * time.Second).C
	}
	fmt.Println("结束了")
}

func Send(ws *websocket.Conn) {
	i := 0
	for {
		i++
		fmt.Println("send i", i)
		time.Sleep(3 * time.Second)
		data := &struct {
			T string
		}{
			T: time.Now().Format("2006-01-02 15:04:05"),
		}

		d, _ := json.Marshal(data)
		fmt.Println("----发送消息", data)
		if err := ws.WriteMessage(2, d); err != nil {
			log.Println("发送消息给客户端发生错误", err.Error())
			// 切断服务
			_ = ws.Close()
			return
		}
	}
}

type Message struct {
	MsgType string      `json:"msg_type"` //消息类型/image/video/...
	Data    interface{} `json:"data"`     //推送消息数组
	Time    time.Time   `json:"time"`     //消息创建时间
}

func Send2(ws *websocket.Conn) {
	i := 0
	for {
		i++
		fmt.Println("send i", i)
		time.Sleep(20 * time.Second)
		data := Message{
			MsgType: "ping",
			Data:    "heartbeat",
			Time:    time.Now(),
		}

		d, _ := json.Marshal(data)
		fmt.Println("----发送消息", data)
		//if err := ws.WriteJSON(d); err != nil {
		if err := ws.WriteMessage(2, d); err != nil {
			log.Println("发送消息给客户端发生错误", err.Error())
			// 切断服务
			_ = ws.Close()
			return
		}
	}
}

func Read(ws *websocket.Conn) {
	i := 0
	for {
		i++
		t, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("读取服务端消息err", err.Error())
			break
		}
		fmt.Println("read i", i, t, string(data), time.Now().Format("2006-01-02 15:04:05"))
		switch t {
		case websocket.PingMessage:
			//收到服务端的ping  --准备回复Pong
			fmt.Println("收到服务端的ping消息")
			data, _ := json.Marshal("pong")
			SendOneMsg(websocket.PongMessage, data, ws)
		default:
			continue
		}
	}
	fmt.Println("读取出错，退出读部分")
}

func SendOneMsg(t int, data []byte, ws *websocket.Conn) {
	fmt.Println("发送消息", t, string(data))
	if err := ws.WriteMessage(t, data); err != nil {
		log.Println("发送消息错误", err.Error())
		// 切断服务
		_ = ws.Close()
		return
	}
}
