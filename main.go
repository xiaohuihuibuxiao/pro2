package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	. "service.jtthink.com/Services"
)

func main()  {

	user:= UserService{}
	endp:= GenUserEndpoint(user)

	serverHanlder:=httptransport.NewServer(endp,DecodeUserRequest,EncodeUserResponse)


	http.ListenAndServe(":8080",serverHanlder)

}