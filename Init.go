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

	r := mymux.NewRouter()
	r.Use()

	//---用户相关----
	usercreate_handler := httptransport.NewServer(UserCreateEndpoint(UserCreateService{}), DecodeUserCreateRequest, EncodeuUserCreateResponse)
	userlogin_handler := httptransport.NewServer(UserLoginEndpoint(UserLoginService{}), DecodeUserLoginRequest, EncodeuUserLoginResponse)
	//	r.Handle(`/user/{uid:\d+}`,serverHanlder)
	r.Methods("POST").Path(`/user/login/{userid}`).Handler(userlogin_handler) //--登陆--ok
	r.Methods("POST").Path(`/user/register`).Handler(usercreate_handler)      //--创建新用户--ok

	//-----设备相关----
	devicecreate_handler := httptransport.NewServer(DeviceCreateEndpoint(DeviceCreateService{}), DecodeDeviceCreateRequest, EncodeDeviceCreateReponse)
	devicedelete_handler := httptransport.NewServer(DeviceDeleteEndpoint(DeviceDeleteService{}), DecodeDeviceDeleteRequest, EncodeDeviceDeleteReponse)
	devicequery_handler := httptransport.NewServer(DeviceQueryEndpoint(DeviceQUeryService{}), DecodeDeviceQueryRequest, EncodeDeviceQueryReponse)
	devicerevise_handler := httptransport.NewServer(DeviceReviseEndpoint(DeviceReviseService{}), DecodeDeviceReviseRequest, EncodeDeviceReviseReponse)
	devicebind_handler := httptransport.NewServer(DeviceBindEndpoint(DeviceBindService{}), DecodeDeviceBindRequest, EncodeDeviceBindReponse)
	deviceunbound_handler := httptransport.NewServer(DeviceUnboundEndpoint(DeviceUnboundService{}), DecodeDeviceUnboundRequest, EncodeDeviceUnboundReponse)
	deviceupload_handler := httptransport.NewServer(DeviceUploadEndpoint(DeviceUploadService{}), DecodeDeviceUploadRequest, EncodeDeviceUploadReponse)

	r.Methods("POST").Path(`/device/{deviceid}`).Handler(devicecreate_handler)                              //---新建设备---ok
	r.Methods("DELETE").Path(`/device/{deviceid}`).Handler(devicedelete_handler)                            //--删除设备--ok
	r.Methods("GET").Path(`/device/{deviceid}`).Handler(devicequery_handler)                                //---查询设备--ok
	r.Methods("PUT").Path(`/device/{deviceid}`).Handler(devicerevise_handler)                               //--修改设备--ok
	r.Methods("PUT").Path(`/device/bind/{deviceid}/{userid}/{sid}/{gatewayid}`).Handler(devicebind_handler) //--绑定--ok
	r.Methods("PUT").Path(`/device/unbound/{deviceid}/{type}`).Handler(deviceunbound_handler)               //--解绑--ok
	r.Methods("POST").Path(`/device/upload/{deviceid}`).Handler(deviceupload_handler)                       //--上报数据--ok

	//--空间相关--
	spacecreate_handler := httptransport.NewServer(SpaceCreateEndpoint(SpaceCreateService{}), DecodeSpaceCreateRequest, EncodeSpaceCreateReponse)
	spacequery_handler := httptransport.NewServer(SpaceQueryEndpoint(SpaceQueryService{}), DecodeSpaceQUeryRequest, EncodeSpaceQueryReponse)
	spacerevise_handler := httptransport.NewServer(SpaceReviseEndpoint(SpaceReviseService{}), DecodeSpaceReviseRequest, EncodeSpaceReviseReponse)
	spacedel_handler := httptransport.NewServer(SpaceDelEndpoint(SpaceDelService{}), DecodeSpaceDelRequest, EncodeSpaceDelReponse)
	spaceclone_handler := httptransport.NewServer(SpaceCloneEndpoint(SpaceCloneService{}), DecodeCloneRequest, EncodeSpaceCloneReponse)

	r.Methods("POST").Path(`/space`).Handler(spacecreate_handler)      //--创建空间--ok
	r.Methods("GET").Path(`/space/{sid}`).Handler(spacequery_handler)  //--查询空间--ok
	r.Methods("PUT").Path(`/space/{sid}`).Handler(spacerevise_handler) //--修改空间--ok
	r.Methods("DELETE").Path(`/space/{sid}`).Handler(spacedel_handler) //--删除空间--ok--一般不要删除空间，需要换个空间的话
	// 直接解绑设备然后新建空间再绑定就好
	r.Methods("POST").Path(`/space/clone/{sid}`).Handler(spaceclone_handler) //--复制空间--ok TODO

	return r
}
