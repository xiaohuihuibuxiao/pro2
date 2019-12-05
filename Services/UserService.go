package Services

type IUserService interface {
	GetName(userid int) string
}
type UserService struct{}

func (this UserService) GetName(userid int) string {
	if userid == 101 {
		return "shenyi"
	}
	return "guest"
}

//------------------------------------

type WUserLoginService interface {
	Login(userid string, pwd string) *CommonResponse
}

type UserLoginService struct{}

func (this UserLoginService) Login(userid string, pwd string) *CommonResponse {
	if userid == pwd {
		return &CommonResponse{
			Code:   200,
			Msg:    "",
			Data:   "1111",
			Expand: nil,
		}
	}
	return &CommonResponse{
		Code:   500,
		Msg:    "err",
		Data:   nil,
		Expand: nil,
	}
}
