package Services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	Uid int `json:"uid"`
}
type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := userService.GetName(r.Uid)
		return UserResponse{Result: result}, nil
	}
}

//-------------------登陆---------------------
type UserLoginRequest struct {
	Userid   string `json:"userid"`
	Password string `json:"password"`
}

type CommonResponse struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Expand interface{} `json:"expand"`
}

func UserLoginEndpoint(userloginService WUserLoginService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) { //TODO requesr数据内容哪来的
		r := request.(*UserLoginRequest)
		result := userloginService.Login(r.Userid, r.Password)
		return result, nil
	}
}
