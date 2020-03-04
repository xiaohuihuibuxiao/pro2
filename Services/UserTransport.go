package Services

import (
	"context"
	"encoding/json"
	mymux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"pro2/Baseinfo"
	"pro2/util"
)

//------------------登陆---------------------------------
func DecodeUserLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var info struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
	}
	_ = json.Unmarshal(body, &info)
	return &UserLoginRequest{
		UserId:   info.UserId,
		Password: info.Password,
		Method:   r.Method,
		Url:      r.URL.String(),
	}, nil
}

func EncodeuUserLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//----------------新建用户----------
func DecodeUserCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var newuser struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Title    string `json:"title"`
		Nickname string `json:"nickname"`
	}
	_ = json.Unmarshal(body, &newuser)
	return &CommonRequest{
		Token:  r.Header.Get("token"),
		Method: r.Method,
		Url:    r.URL.String(),
		Msg: &UserCreateRequest{
			UserId:   newuser.UserId,
			Password: newuser.Password,
			Phone:    newuser.Phone,
			Email:    newuser.Email,
			Title:    newuser.Title,
			Nickname: newuser.Nickname,
		},
	}, nil
}

//------------------邪恶的分割线-----------------
func MyErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("content-type", contentType)
	if myerr, ok := err.(*util.MyError); ok {
		w.WriteHeader(myerr.Code)
		w.Write(body)
	} else {
		w.WriteHeader(500)
		w.Write(body)
	}
}

//------------------获取用户列表-------------------------
func DecodeUserListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	token := r.Header.Get("token")
	ok, err, user := Baseinfo.TokenCheckAsymmetricalKey(token)
	if !ok {
		return nil, err
	}
	return &CommonRequest{
		Token:  r.Header.Get("token"),
		Method: r.Method,
		Url:    r.URL.String(),
		Msg: &UserListRequest{
			UserId: user,
		},
	}, nil
}

//------------------编辑用户-------------------------
func DecodeUserEditRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var userid = ""
	vars := mymux.Vars(r)
	if user, ok := vars["userid"]; ok {
		userid = user
	}
	body, _ := ioutil.ReadAll(r.Body)
	var bodyinfo struct {
		Phone    string `json:"Phone"`
		Title    string `json:"title"`
		Nickname string `json:"nickname"`
		Email    string `json:"Email"`
	}
	_ = json.Unmarshal(body, &bodyinfo)
	return &CommonRequest{
		Token:  r.Header.Get("token"),
		Method: r.Method,
		Url:    r.URL.String(),
		Msg: &UserEditRequest{
			UserId:   userid,
			Phone:    bodyinfo.Phone,
			Title:    bodyinfo.Title,
			Nickname: bodyinfo.Nickname,
			Email:    bodyinfo.Email,
		},
	}, nil
}

//------------------删除用户-------------------------
func DecodeUserDelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var userid string = ""
	vars := mymux.Vars(r)
	if user, ok := vars["userid"]; ok {
		userid = user
	}
	return &CommonRequest{
		Token:  r.Header.Get("token"),
		Method: r.Method,
		Url:    r.URL.String(),
		Msg: &UserDelRequest{
			UserId: userid,
		},
	}, nil
}

//------------------修改用户密码-------------------------
func DecodeUserResetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var info struct {
		UserId           string `json:"userId"`
		OriginalPassword string `json:"originalPassword"`
		NewPassword      string `json:"newPassword"`
	}
	_ = json.Unmarshal(body, &info)
	return &CommonRequest{
		Token:  r.Header.Get("token"),
		Method: r.Method,
		Url:    r.URL.String(),
		Msg: &UserResetRequest{
			UserId:      info.UserId,
			NewPassword: info.NewPassword,
		},
	}, nil
}

//------------------用户注销-------------------------
func DecodeUserLogoutRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var bodyinfo struct {
		Userid string `json:"Userid"`
	}
	_ = json.Unmarshal(body, &bodyinfo)
	return &UserLogoutRequest{
		Userid: bodyinfo.Userid,
	}, nil
}

func EncodeuUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
