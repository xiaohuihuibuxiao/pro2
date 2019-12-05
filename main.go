package main

import (
	"github.com/astaxie/beego"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	. "pro2/Services"
)

func main() {

	user := UserService{}
	endp := GenUserEndpoint(user)

	serverHanlder := httptransport.NewServer(endp, DecodeUserRequest, EncodeUserResponse)
	app := beego.Handler("/user/login", serverHanlder)

	http.ListenAndServe(":8080", app.Handlers)
}
