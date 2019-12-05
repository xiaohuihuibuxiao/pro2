package Services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

//---------------------------------------------------

func DecodeUserLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userid := r.URL.Query().Get("userid") //TODO 这个是从？后面去读参数，即接口的param内容
	password := r.URL.Query().Get("password")
	fmt.Println("userid,pwd", userid, password)
	vars := mymux.Vars(r)
	if name, ok := vars["name"]; ok {
		fmt.Println("name", name)
	}
	bodu, _ := ioutil.ReadAll(r.Body)
	var aa struct {
		Name string `json:"name"`
	}
	json.Unmarshal(bodu, &aa)
	fmt.Println("body", bodu)

	return &UserLoginRequest{
		Userid:   userid,
		Password: password,
	}, nil
}

func EncodeuUserLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
