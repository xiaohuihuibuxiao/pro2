package main

import (
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"os"
	. "pro2/Services"
)

func Init() *mymux.Router {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "mykit", "1.0")
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	///user := UserService{} //用户服务
	//limit := rate.NewLimiter(1, 5)
	//	endp := RateLimit(limit)((CheckTokenMiddleware()(GenUserEndpoint(user))))
	//	endp:=RateLimit(limit)(UserServiceLogMiddleware(logger)(CheckTokenMiddleware()(GenUserEndpoint(user))))

	///options := []httptransport.ServerOption{
	//	httptransport.ServerErrorEncoder(MyErrorEncoder), //????
	//}
	//serverHanlder := httptransport.NewServer(endp, DecodeUserRequest, EncodeUserResponse, options...)

	r := mymux.NewRouter()

	//---用户管理----
	limit := rate.NewLimiter(1, 5) // 限制频繁登陆操作

	//httptransport.NewServer(UserServiceLogMiddleware(logger)(UserLoginEndpoint(UserLoginService{})), DecodeDeviceCreateRequest, EncodeDeviceCreateReponse)

	// /isms.v1
	userListHandler := httptransport.NewServer(UserListEndpoint(UserListService{}), DecodeUserListRequest, EncodeuUserResponse)
	userCreateHandler := httptransport.NewServer(UserCreateEndpoint(UserCreateService{}), DecodeUserCreateRequest, EncodeuUserResponse)
	userEditHandler := httptransport.NewServer(UserEditEndpoint(UserEditService{}), DecodeUserEditRequest, EncodeuUserResponse)
	userDelHandler := httptransport.NewServer(UserDelEndpoint(UserDelService{}), DecodeUserDelRequest, EncodeuUserResponse)
	userResetHandler := httptransport.NewServer(UserResetEndpoint(UserResetService{}), DecodeUserResetRequest, EncodeuUserResponse)
	userLoginHandler := httptransport.NewServer(RateLimit(limit)(UserServiceLogMiddleware(logger)(UserLoginEndpoint(UserLoginService{}))), DecodeUserLoginRequest, EncodeuUserLoginResponse)
	userLogoutHandler := httptransport.NewServer(UserLogoutEndpoint(UserLogoutService{}), DecodeUserLogoutRequest, EncodeuUserResponse)

	r.Methods("GET").Path(`/isms/v1/user`).Handler(userListHandler)              //--3.1.1获取用户列表
	r.Methods("POST").Path(`/isms/v1/user`).Handler(userCreateHandler)           //--3.1.2创建用户
	r.Methods("PUT").Path(`/isms/v1/user/{userid}`).Handler(userEditHandler)     //--3.1.3编辑用户
	r.Methods("DELETE").Path(`/isms/v1/user/{userid}`).Handler(userDelHandler)   //--3.1.4删除用户
	r.Methods("DELETE").Path(`/isms/v1/user/resetpwd`).Handler(userResetHandler) //--3.1.5 更改用户密码
	r.Methods("POST").Path(`/isms/v1//user/login`).Handler(userLoginHandler)     //--3.1.6 用户登陆
	r.Methods("POST").Path(`/isms/v1//user/logout`).Handler(userLogoutHandler)   //--3.1.6 用户注销

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
	r.Methods("DELETE").Path(`/space/{sid}`).Handler(spacedel_handler) //--删除空间--ok--一般不要删除空间，需要的话新建一个空间并绑定即可
	// 直接解绑设备然后新建空间再绑定就好
	r.Methods("POST").Path(`/space/clone/{sid}`).Handler(spaceclone_handler) //--复制空间--ok

	//将某一级以及以下的设备全都迁移到另一个空间下面（层级关系肯定得不能乱）

	return r
}
