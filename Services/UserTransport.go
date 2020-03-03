package Services

import (
	"context"
	"encoding/json"
	mymux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"pro2/util"
)

//------------------登陆---------------------------------
func DecodeUserLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	body, _ := ioutil.ReadAll(r.Body)
	var bodyinfo struct {
		Userid   string `json:"Userid"`
		Password string `json:"Password"`
	}
	json.Unmarshal(body, &bodyinfo)
	password := bodyinfo.Password
	return &UserLoginRequest{
		Userid:   bodyinfo.Userid,
		Password: password,
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
		Userid   string `json:"userid"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Title    string `json:"title"`
		Nickname string `json:"nickname"`
	}
	json.Unmarshal(body, &newuser)
	a := &UserCreateRequest{
		Userid:   newuser.Userid,
		Password: newuser.Password,
		Phone:    newuser.Phone,
		Email:    newuser.Email,
		Title:    newuser.Title,
		Nickname: newuser.Nickname,
	}
	return a, nil
}

func EncodeuUserCreateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
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
//TODO 函数的实现
func DecodeUserListRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	//var userid string = ""
	//vars := mymux.Vars(r)
	//if user, ok := vars["userid"]; ok {
	//	userid = user
	//}
	//body, _ := ioutil.ReadAll(r.Body)
	//var bodyinfo struct {
	//	Password string `json:"password"`
	//}
	//json.Unmarshal(body, &bodyinfo)
	//password := bodyinfo.Password
	//a := &UserLoginRequest{
	//	Userid:   userid,
	//	Password: password,
	//	Method:   r.Method,
	//	Url:      r.URL.String(),
	//}
	return nil, nil
}

//------------------编辑用户-------------------------
func DecodeUserEditRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var userid string = ""
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
	return &UserEditRequest{
		Userid:   userid,
		Phone:    bodyinfo.Phone,
		Title:    bodyinfo.Title,
		Nickname: bodyinfo.Nickname,
		Email:    bodyinfo.Email,
	}, nil
}

//------------------删除用户-------------------------
func DecodeUserDelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var userid string = ""
	vars := mymux.Vars(r)
	if user, ok := vars["userid"]; ok {
		userid = user
	}
	return &UserDelRequest{
		Userid: userid,
	}, nil
}

//------------------修改用户密码-------------------------
func DecodeUserResetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	body, _ := ioutil.ReadAll(r.Body)
	var bodyinfo struct {
		Userid           string `json:"Userid"`
		Originalpassword string `json:"Originalpassword"`
		Newpassword      string `json:"Newpassword"`
	}
	_ = json.Unmarshal(body, &bodyinfo)
	return &UserResetRequest{
		Userid:           bodyinfo.Userid,
		Originalpassword: bodyinfo.Originalpassword,
		Newpassword:      bodyinfo.Newpassword,
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
