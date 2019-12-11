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
	Userid   string `json:"userid"`
	Password string `json:"password"`
	Method   string `json:"method"`
	Url      string `json:"url"`
}

func UserLoginEndpoint(userloginService WUserLoginService) endpoint.Endpoint {
	fmt.Println("bbb")
	return func(ctx context.Context, request interface{}) (response interface{}, err error) { //TODO requesr数据内容哪来的
		r := request.(*UserLoginRequest)
		result := userloginService.Login(r.Userid, r.Password)
		return result, nil
	}
}

//--------创建新用户----------
type UserCreateRequest struct {
	Userid   string `json:"userid"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Title    string `json:"title"`
	Nickname string `json:"nickname"`
}

func UserCreateEndpoint(userCreateService WUserCreateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) { //TODO requesr数据内容哪来的
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
	fmt.Println("aa")
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//	r := request.(UserRequest)
			r := request.(*UserLoginRequest)

			if r.Method != "GET" {
				Baseinfo.RecordOperation(r.Url, r.Method, r.Userid)
			}

			logger.Log("method", r.Method, "url", r.Url)
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
