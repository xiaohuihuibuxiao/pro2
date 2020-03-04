package Services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"pro2/Baseinfo"
	"pro2/util"
)

//----common--------
type CommonResponse struct {
	Code   int64       `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Expand interface{} `json:"expand"`
}

type CommonRequest struct {
	Token  string      `json:"token"`
	Method string      `json:"method"`
	Url    string      `json:"url"`
	Msg    interface{} `json:"msg"`
}

//--------登陆-----------
type UserLoginRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
	Method   string `json:"method"`
	Url      string `json:"url"`
}

func UserLoginEndpoint(userloginService WUserLoginService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*UserLoginRequest)
		result := userloginService.Login(r.UserId, r.Password)
		return result, nil
	}
}

//--------创建新用户----------
type UserCreateRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Title    string `json:"title"`
	Nickname string `json:"nickname"`
}

func UserCreateEndpoint(userCreateService WUserCreateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*CommonRequest)
		result := userCreateService.NewAccount(r.Msg.(*UserCreateRequest))
		return result, nil
	}
}

//---------------------middleware---------------------
//加入限流功能的 中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, util.NewMyError(429, "too many requests")
			}
			return next(ctx, request)
		}
	}
}

type CommonLogger struct {
	Token  string `json:"token"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

//登陆日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(*UserLoginRequest)
			if r.Method != "GET" {
				Baseinfo.RecordOperation(r.Url, r.Method, r.UserId)
			}

			_ = logger.Log("method", r.Method, "url", r.Url)
			return next(ctx, request)
		}
	}
}

//token验证中间件
func CheckTokenMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(*CommonRequest)
			err, user0 := Baseinfo.LoginTokenAuth(r.Token)
			if err != nil {
				return nil, util.NewMyError(403, "error token")

			}
			newCtx := context.WithValue(ctx, "LoginUser", user0)
			return next(newCtx, request)
		}
	}
}

//--------获取用户列表----------
type UserListRequest struct {
	UserId string `json:"userId"`
}

func UserListEndpoint(userListService WUserListService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*CommonRequest)
		result := userListService.ObtainUserList(r.Msg.(*UserListRequest).UserId)
		return result, nil
	}
}

//--------编辑用户----------
type UserEditRequest struct {
	UserId   string `json:"userId"`
	Phone    string `json:"phone"`
	Title    string `json:"title"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func UserEditEndpoint(userEditService WUserEditService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*CommonRequest)
		result := userEditService.UserEdit(r.Msg.(*UserEditRequest), r.Token)
		return result, nil
	}
}

//--------删除用户----------
type UserDelRequest struct {
	UserId string `json:"userId"`
}

func UserDelEndpoint(userDelService WUserDelService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*CommonRequest)
		result := userDelService.UserDel(r.Msg.(*UserDelRequest).UserId, r.Token)
		return result, nil
	}
}

//--------修改用户密码----------
type UserResetRequest struct {
	UserId      string `json:"userId"`
	NewPassword string `json:"newPassword"`
}

func UserResetEndpoint(userResetService WUserResetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*CommonRequest)
		result := userResetService.UserReset(r.Msg.(*UserResetRequest), r.Token)
		return result, nil
	}
}

//--------用户注销----------
type UserLogoutRequest struct {
	Userid string `json:"Userid"`
}

func UserLogoutEndpoint(userLogoutService WUserLogoutService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*UserLogoutRequest)
		result := userLogoutService.UserLogout(r.Userid)
		return result, nil
	}
}
