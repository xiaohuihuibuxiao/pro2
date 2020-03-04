package Services

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

//--------登陆-----------
type UserLoginRequest struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
	Method   string `json:"method"`
	Url      string `json:"url"`
}

func UserLoginEndpoint(userloginService WUserLoginService) endpoint.Endpoint {
	fmt.Println("登陆endpoint")
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fmt.Println("登陆1111")
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
		r := request.(*UserCreateRequest)
		result := userCreateService.NewAccount(r)
		return result, nil
	}
}

//---------------------midware---------------------
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
	fmt.Println("进入登陆日志中间件")
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		fmt.Println("登陆中间件00000")
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			fmt.Println("00000")
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
			r := request.(CommonLogger)
			uc := UserClaim{}
			getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (i interface{}, e error) {
				return []byte(secKey), nil
			})
			if getToken != nil && getToken.Valid { //验证通过
				newCtx := context.WithValue(ctx, "LoginUser", getToken.Claims.(*UserClaim).Uname)
				return next(newCtx, request)
			} else {
				return nil, util.NewMyError(403, "error token")
			}
		}
	}
}

//--------获取用户列表----------
type UserListRequest struct {
	UserId string `json:"userId"`
}

func UserListEndpoint(userListService WUserListService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*UserListRequest)
		result := userListService.ObtainUserList(r.UserId)
		return result, nil
	}
}

//--------编辑用户----------
type UserEditRequest struct {
	Userid   string `json:"userid"`
	Phone    string `json:"Phone"`
	Title    string `json:"Title"`
	Nickname string `json:"Nickname"`
	Email    string `json:"Email"`
}

func UserEditEndpoint(userEditService WUserEditService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*UserEditRequest)
		result := userEditService.UserEdit(r)
		return result, nil
	}
}

//--------删除用户----------
type UserDelRequest struct {
	Userid string `json:"userid"`
}

func UserDelEndpoint(userDelService WUserDelService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*UserDelRequest)
		result := userDelService.UserDel(r.Userid)
		return result, nil
	}
}

//--------修改用户密码----------
type UserResetRequest struct {
	Userid           string `json:"Userid"`
	Originalpassword string `json:"Originalpassword"`
	Newpassword      string `json:"Newpassword"`
}

func UserResetEndpoint(userResetService WUserResetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*UserResetRequest)
		result := userResetService.UserReset(r)
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
