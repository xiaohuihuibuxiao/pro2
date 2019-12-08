package Services

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//--创建空间--
func DecodeSpaceCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var bodyinfo struct {
		Level    int    `json:"level"`
		Province string `json:"province"`
		City     string `json:"city"`
		Area     string `json:"area"`
		District string `json:"district"`
		Building string `json:"building"`
		Storey   string `json:"storey"`
		Room     string `json:"room"`
		Place    string `json:"place"`
		Title    string `json:"title"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &bodyinfo)
	if err != nil {
		return nil, err
	}
	return &SpaceCreateRequest{
		Token:    r.Header.Get("token"),
		Level:    bodyinfo.Level,
		Province: bodyinfo.Province,
		City:     bodyinfo.City,
		Area:     bodyinfo.Area,
		District: bodyinfo.District,
		Building: bodyinfo.Building,
		Storey:   bodyinfo.Storey,
		Room:     bodyinfo.Room,
		Place:    bodyinfo.Place,
		Title:    bodyinfo.Title,
	}, nil
}

func EncodeSpaceCreateReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--查询空间--
func DecodeSpaceQUeryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &SpaceQueryRequest{
		Sid:   mux.Vars(r)["sid"],
		Token: r.Header.Get("token"),
	}, nil
}

func EncodeSpaceQueryReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--修改空间--
func DecodeSpaceReviseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var bodyinfo struct {
		Title string `json:"title"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &bodyinfo)
	if err != nil {
		return nil, err
	}
	return &SpaceReviseRequest{
		Title: bodyinfo.Title,
		Token: r.Header.Get("token"),
		Sid:   mux.Vars(r)["sid"],
	}, nil
}

func EncodeSpaceReviseReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--删除空间--
func DecodeSpaceDelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &SpaceReviseRequest{
		Token: r.Header.Get("token"),
		Sid:   mux.Vars(r)["sid"],
	}, nil
}

func EncodeSpaceDelReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--复制空间--
func DecodeCloneRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &SpaceReviseRequest{
		Token: r.Header.Get("token"),
		Sid:   mux.Vars(r)["sid"],
	}, nil
}

func EncodeSpaceCloneReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
