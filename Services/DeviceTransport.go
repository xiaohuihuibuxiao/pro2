package Services

import (
	"context"
	"encoding/json"
	mymux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"pro2/Baseinfo"
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
	ok, _, userid := Baseinfo.TokenCheck_asymmetricalkey(token)
	if ok {
		Baseinfo.RecordOperation(r.URL.String(), r.Method, userid)
	}
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
	ok, _, userid := Baseinfo.TokenCheck_asymmetricalkey(r.Header.Get("token"))
	if ok {
		Baseinfo.RecordOperation(r.URL.String(), r.Method, userid)
	}
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
	ok, _, userid := Baseinfo.TokenCheck_asymmetricalkey(r.Header.Get("token"))
	if ok {
		Baseinfo.RecordOperation(r.URL.String(), r.Method, userid)
	}
	return &DeviceDeleteRequest{
		Token:    r.Header.Get("token"),
		Deviceid: mymux.Vars(r)["deviceid"],
	}, nil
}

func EncodeDeviceDeleteReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--修改设备--
func DecodeDeviceReviseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var bodyinfo struct {
		Title    string                 `json:"title"`
		External map[string]interface{} `json:"external"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &bodyinfo)
	ok, _, userid := Baseinfo.TokenCheck_asymmetricalkey(r.Header.Get("token"))
	if ok {
		Baseinfo.RecordOperation(r.URL.String(), r.Method, userid)
	}
	return &DeviceReviseRequest{
		Token:    r.Header.Get("token"),
		Deviceid: mymux.Vars(r)["deviceid"],
		External: bodyinfo.External,
		Title:    bodyinfo.Title,
	}, nil
}

func EncodeDeviceReviseReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--绑定设备--
func DecodeDeviceBindRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	ok, _, userid := Baseinfo.TokenCheck_asymmetricalkey(r.Header.Get("token"))
	if ok {
		Baseinfo.RecordOperation(r.URL.String(), r.Method, userid)
	}
	return &DeviceBindRequest{
		Token:     r.Header.Get("token"),
		Deviceid:  mymux.Vars(r)["deviceid"],
		Gatewayid: mymux.Vars(r)["gatewayid"],
		Sid:       mymux.Vars(r)["sid"],
		Userid:    mymux.Vars(r)["userid"],
	}, nil
}

func EncodeDeviceBindReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--解绑设备--
func DecodeDeviceUnboundRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	ok, _, userid := Baseinfo.TokenCheck_asymmetricalkey(r.Header.Get("token"))
	if ok {
		Baseinfo.RecordOperation(r.URL.String(), r.Method, userid)
	}
	//t, _ := strconv.Atoi(mymux.Vars(r)["type"])
	return &DeviceUnboundRequest{
		Token:    r.Header.Get("token"),
		Deviceid: mymux.Vars(r)["deviceid"],
		Type:     mymux.Vars(r)["type"],
	}, nil
}

func EncodeDeviceUnboundReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

//--上传数据--
func DecodeDeviceUploadRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var bodyinfo struct {
		Deviceid string      `json:"deviceid"`
		T        string      `json:"t"`
		Userid   string      `json:"userid"`
		Devtype  int         `json:"devtype"`
		Data     interface{} `json:"data"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &bodyinfo)
	ok, _, userid := Baseinfo.TokenCheck_asymmetricalkey(r.Header.Get("token"))
	if ok {
		Baseinfo.RecordOperation(r.URL.String(), r.Method, userid)
	}
	return &DeviceUploadRequest{
		Deviceid: mymux.Vars(r)["deviceid"],
		T:        bodyinfo.T,
		Userid:   bodyinfo.Userid,
		Devtype:  bodyinfo.Devtype,
		Data:     bodyinfo.Data,
	}, nil
}

func EncodeDeviceUploadReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
