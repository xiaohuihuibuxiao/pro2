package Services

import (
	"context"
	"encoding/json"
	"errors"
	mymux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	// http://localhost:xxx/?uid=101
	if r.URL.Query().Get("uid") != "" {
		uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
		return UserRequest{
			Uid: uid,
		}, nil
	}
	return nil, errors.New("参数错误")

}
func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//------------------登陆---------------------------------

func DecodeUserLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var userid string = ""
	vars := mymux.Vars(r)
	if user, ok := vars["userid"]; ok {
		userid = user
	}
	body, _ := ioutil.ReadAll(r.Body)
	var bodyinfo struct {
		Password string `json:"password"`
	}
	json.Unmarshal(body, &bodyinfo)
	password := bodyinfo.Password
	a := &UserLoginRequest{
		Userid:   userid,
		Password: password,
	}
	return a, nil
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
