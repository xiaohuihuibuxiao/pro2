package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"net/http"
	. "pro2/Services"
)

func main() {

	//user := UserService{}
	//endp := GenUserEndpoint(user)
	//serverHanlder := httptransport.NewServer(endp, DecodeUserRequest, EncodeUserResponse)

	user := UserLoginService{}
	endp_user := UserLoginEndpoint(user)
	serverHanlder := httptransport.NewServer(endp_user, DecodeUserLoginRequest, EncodeuUserLoginResponse)
	r := mymux.NewRouter()
	//	r.Handle(`/user/{uid:\d+}`,serverHanlder)
	r.Methods("GET").Path(`/user/login/{name}`).Handler(serverHanlder)

	http.ListenAndServe(":8080", r)

}
