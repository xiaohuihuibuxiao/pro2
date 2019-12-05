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
}

func UserLoginEndpoint(userloginService WUserLoginService) endpoint.Endpoint {
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
