package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"net/http"
	. "pro2/Services"
)

func main() {
	r := mymux.NewRouter()
	//--登陆--
	userlogin := UserLoginService{}
	endp_user := UserLoginEndpoint(userlogin)
	userlogin_handler := httptransport.NewServer(endp_user, DecodeUserLoginRequest, EncodeuUserLoginResponse)
	//	r.Handle(`/user/{uid:\d+}`,serverHanlder)
	r.Methods("POST").Path(`/user/login/{userid}`).Handler(userlogin_handler)

	//--创建新用户--
	usercreats := UserCreateService{}
	endp_usercreate := UserCreateEndpoint(usercreats)
	usercreate_handler := httptransport.NewServer(endp_usercreate, DecodeUserCreateRequest, EncodeuUserCreateResponse)
	r.Methods("POST").Path(`/user/register`).Handler(usercreate_handler)

	//---新建设备---
	devicecreate := DeviceCreateService{}
	ep_devicecreats := DeviceCreateEndpoint(devicecreate)
	devicecreate_handler := httptransport.NewServer(ep_devicecreats, DecodeDeviceCreateRequest, EncodeDeviceCreateReponse)
	r.Methods("GET").Path(`/device/{deviceid}`).Handler(devicecreate_handler)
	//---查询设备--

	//--删除设备--

	//--修改设备--

	//--绑定设备--

	//--解绑设备--

	//--设备上报数据--

	//--创建空间--

	//--查询空间--

	//--删除空间--

	//--修改空间--

	//--给空间添加设备--

	//--清除空间信息--

	//--复制空间---

	http.ListenAndServe(":8080", r)

}
