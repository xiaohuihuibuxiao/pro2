package main

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"os"
	. "pro2/Services"

	"golang.org/x/time/rate"
)

func Init() *mymux.Router {

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "mykit", "1.0")
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	//---其实就是用于用户登陆的 可以改写下 利用其限流器
	user := UserService{} //用户服务
	limit := rate.NewLimiter(1, 5)
	endp := RateLimit(limit)((CheckTokenMiddleware()(GenUserEndpoint(user)))) //这个可以改写用户登陆接口 限制登陆次数
	//endp:=RateLimit(limit)(UserServiceLogMiddleware(logger)(CheckTokenMiddleware()(GenUserEndpoint(user))))

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(MyErrorEncoder), //????
	}
	serverHanlder := httptransport.NewServer(endp, DecodeUserRequest, EncodeUserResponse, options...)
	fmt.Println(serverHanlder)

	////增加handler 用于获取用户token
	//accessService:=&AccessService{}
	//accessServiceEndpoint:=AccessEndpoint(accessService)
	//accessHandler:=httptransport.NewServer(accessServiceEndpoint,DecodeAccessRequest,EncodeAccessResponse,options...)//只是用于生成token，已有该功能
	//

	r := mymux.NewRouter()

	//---用户相关----
	//userlogin := UserLoginService{}
	//usercreats := UserCreateService{}
	endp_user := UserLoginEndpoint(UserLoginService{})
	endp_usercreate := UserCreateEndpoint(UserCreateService{})
	usercreate_handler := httptransport.NewServer(endp_usercreate, DecodeUserCreateRequest, EncodeuUserCreateResponse)
	userlogin_handler := httptransport.NewServer(endp_user, DecodeUserLoginRequest, EncodeuUserLoginResponse)

	//	r.Handle(`/user/{uid:\d+}`,serverHanlder)
	r.Methods("POST").Path(`/user/login/{userid}`).Handler(userlogin_handler) //--登陆--ok
	r.Methods("POST").Path(`/user/register`).Handler(usercreate_handler)      //--创建新用户--ok

	//-----设备相关----

	//devicecreate := DeviceCreateService{}
	//devicedelete := DeviceDeleteService{}
	//deviceQuery := DeviceQUeryService{}
	//devicerevise := DeviceReviseService{}
	//devicebind := DeviceBindService{}
	//deviceunbound := DeviceUnboundService{}
	//deviceupload := DeviceUploadService{}

	ep_devicecreats := DeviceCreateEndpoint(DeviceCreateService{})
	ep_devicedelete := DeviceDeleteEndpoint(DeviceDeleteService{})
	ep_devicequery := DeviceQueryEndpoint(DeviceQUeryService{})
	ep_devicerevise := DeviceReviseEndpoint(DeviceReviseService{})
	ep_devicebind := DeviceBindEndpoint(DeviceBindService{})
	ep_deviceunbound := DeviceUnboundEndpoint(DeviceUnboundService{})
	ep_deviceupload := DeviceUploadEndpoint(DeviceUploadService{})

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

	//spacecreate := SpaceCreateService{}
	//spacequery := SpaceQueryService{}
	//spacerevise := SpaceReviseService{}
	//spacedel := SpaceDelService{}
	//spaceclone := SpaceCloneService{}

	ep_spacescreate := SpaceCreateEndpoint(SpaceCreateService{})
	ep_spacequery := SpaceQueryEndpoint(SpaceQueryService{})
	ep_spacerevise := SpaceReviseEndpoint(SpaceReviseService{})
	ep_spacedel := SpaceDelEndpoint(SpaceDelService{})
	ep_spaceclone := SpaceCloneEndpoint(SpaceCloneService{})

	spacecreate_handler := httptransport.NewServer(ep_spacescreate, DecodeSpaceCreateRequest, EncodeSpaceCreateReponse)
	spacequery_handler := httptransport.NewServer(ep_spacequery, DecodeSpaceQUeryRequest, EncodeSpaceQueryReponse)
	spacerevise_handler := httptransport.NewServer(ep_spacerevise, DecodeSpaceReviseRequest, EncodeSpaceReviseReponse)
	spacedel_handler := httptransport.NewServer(ep_spacedel, DecodeSpaceDelRequest, EncodeSpaceDelReponse)
	spaceclone_handler := httptransport.NewServer(ep_spaceclone, DecodeCloneRequest, EncodeSpaceCloneReponse)

	r.Methods("POST").Path(`/space`).Handler(spacecreate_handler)            //--创建空间--ok
	r.Methods("GET").Path(`/space/{sid}`).Handler(spacequery_handler)        //--查询空间--ok
	r.Methods("PUT").Path(`/space/{sid}`).Handler(spacerevise_handler)       //--修改空间--ok
	r.Methods("DELETE").Path(`/space/{sid}`).Handler(spacedel_handler)       //--删除空间--ok
	r.Methods("POST").Path(`/space/clone/{sid}`).Handler(spaceclone_handler) //--复制空间--ok
	//--给空间添加设备--
	//--清除空间信息--

	return r
}
