package Services

import (
	"context"
	"encoding/json"
	"fmt"
	mymux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"pro2/util"
)

//------------------登陆---------------------------------
func DecodeUserLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	fmt.Println("ccc")
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
		Method:   r.Method,
		Url:      r.URL.String(),
	}
	return a, nil
}

func EncodeuUserLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Println("ddd")
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
