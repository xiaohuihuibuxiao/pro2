package Services

import (
	"context"
	"encoding/json"
	mymux "github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func DecodeDeviceCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mymux.Vars(r)
	var deviceid string

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
	return &DeviceCteateRequest{
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
