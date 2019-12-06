package Services

import (
	"context"
	"encoding/json"
	mymux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//--新建设备--
func DecodeDeviceCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mymux.Vars(r)
	var deviceid string

	token := r.Header.Get("token")

	if id, ok := vars["deviceid"]; ok {
		deviceid = id
	}
	body, _ := ioutil.ReadAll(r.Body)
	var newdevice struct {
		Isnode  bool   `json:"isnode"`
		Devtype int64  `json:"devtype"`
		Title   string `json:"title"`
	}
	json.Unmarshal(body, &newdevice)
	return &DeviceCreateRequest{
		Token:    token,
		Deviceid: deviceid,
		Isnode:   newdevice.Isnode,
		Devtype:  newdevice.Devtype,
		Title:    newdevice.Title,
	}, nil

}

func EncodeDeviceCreateReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--查询设备--

func DecodeDeviceQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &DeviceQueryRequest{
		Token:    r.Header.Get("token"),
		Deviceid: mymux.Vars(r)["deviceid"],
	}, nil
}

func EncodeDeviceQueryReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--删除设备--
func DecodeDeviceDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &DeviceQueryRequest{
		Token:    r.Header.Get("token"),
		Deviceid: mymux.Vars(r)["deviceid"],
	}, nil
}

func EncodeDeviceDeleteReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
