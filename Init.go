package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	. "pro2/Services"
)

func Init() *mymux.Router {
	r := mymux.NewRouter()

	//---用户相关----
	userlogin := UserLoginService{}
	usercreats := UserCreateService{}
	endp_user := UserLoginEndpoint(userlogin)
	endp_usercreate := UserCreateEndpoint(usercreats)
	usercreate_handler := httptransport.NewServer(endp_usercreate, DecodeUserCreateRequest, EncodeuUserCreateResponse)
	userlogin_handler := httptransport.NewServer(endp_user, DecodeUserLoginRequest, EncodeuUserLoginResponse)

	//	r.Handle(`/user/{uid:\d+}`,serverHanlder)
	r.Methods("POST").Path(`/user/login/{userid}`).Handler(userlogin_handler) //--登陆--ok
	r.Methods("POST").Path(`/user/register`).Handler(usercreate_handler)      //--创建新用户--ok

	//-----设备相关----
	devicecreate := DeviceCreateService{}
	devicedelete := DeviceDeleteService{}
	deviceQuery := DeviceQUeryService{}
	devicerevise := DeviceReviseService{}
	devicebind := DeviceBindService{}
	deviceunbound := DeviceUnboundService{}
	deviceupload := DeviceUploadService{}

	ep_devicecreats := DeviceCreateEndpoint(devicecreate)
	ep_devicedelete := DeviceDeleteEndpoint(devicedelete)
	ep_devicequery := DeviceQueryEndpoint(deviceQuery)
	ep_devicerevise := DeviceReviseEndpoint(devicerevise)
	ep_devicebind := DeviceBindEndpoint(devicebind)
	ep_deviceunbound := DeviceUnboundEndpoint(deviceunbound)
	ep_deviceupload := DeviceUploadEndpoint(deviceupload)

	devicecreate_handler := httptransport.NewServer(ep_devicecreats, DecodeDeviceCreateRequest, EncodeDeviceCreateReponse)
	devicedelete_handler := httptransport.NewServer(ep_devicedelete, DecodeDeviceDeleteRequest, EncodeDeviceDeleteReponse)
	devicequery_handler := httptransport.NewServer(ep_devicequery, DecodeDeviceQueryRequest, EncodeDeviceQueryReponse)
	devicerevise_handler := httptransport.NewServer(ep_devicerevise, DecodeDeviceReviseRequest, EncodeDeviceReviseReponse)
	devicebind_handler := httptransport.NewServer(ep_devicebind, DecodeDeviceBindRequest, EncodeDeviceBindReponse)
	deviceunbound_handler := httptransport.NewServer(ep_deviceunbound, DecodeDeviceUnboundRequest, EncodeDeviceUnboundReponse)
	deviceupload_handler := httptransport.NewServer(ep_deviceupload, DecodeDeviceUploadRequest, EncodeDeviceUploadReponse)

	r.Methods("POST").Path(`/device/{deviceid}`).Handler(devicecreate_handler)                              //---新建设备---ok
	r.Methods("DELETE").Path(`/device/{deviceid}`).Handler(devicedelete_handler)                            //--删除设备--ok 差修改space信息
	r.Methods("GET").Path(`/device/{deviceid}`).Handler(devicequery_handler)                                //---查询设备--ok
	r.Methods("PUT").Path(`/device/{deviceid}`).Handler(devicerevise_handler)                               //--修改设备--ok
	r.Methods("PUT").Path(`/device/bind/{deviceid}/{userid}/{sid}/{gatewayid}`).Handler(devicebind_handler) //--绑定--TODO 未测试
	r.Methods("PUT").Path(`/device/unbound/{deviceid}/{type}`).Handler(deviceunbound_handler)               //--解绑--TODO 未测试
	r.Methods("POST").Path(`/device/upload/{deviceid}`).Handler(deviceupload_handler)                       //--上报数据--ok

	spacecreate := SpaceCreateService{}

	ep_spacescreate := SpaceCreateEndpoint(spacecreate)

	spacecreate_handler := httptransport.NewServer(ep_spacescreate, DecodeSpaceCreateRequest, EncodeSpaceCreateReponse)

	r.Methods("POST").Path(`/space`).Handler(spacecreate_handler) //--创建空间--ok

	//--创建空间--

	//--查询空间--

	//--删除空间--

	//--修改空间--

	//--给空间添加设备--

	//--清除空间信息--

	//--复制空间---

	return r
}
