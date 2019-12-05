package Services

import (
	"context"
	"pro2/Baseinfo"
)

//--新建设备--
type WDeviceCreateService interface {
	NewDevice(r *DeviceCteateRequest) *CommonResponse
}

type DeviceCreateService struct{}

func (this DeviceCreateService) NewDevice(r *DeviceCteateRequest) *CommonResponse {
	col_device := Baseinfo.Client.Database("test").Collection("device")
	response := &CommonResponse{}

	//TODO how to get token
	//token:=

	newdeviceinfo := &Baseinfo.Device{
		Userid:    "",
		Gatewayid: "",
		Deviceid:  r.Deviceid,
		Isnode:    r.Isnode,
		Devtype:   r.Devtype,
		Title:     r.Title,
		Addr:      "",
		Housecode: "",
		Expand:    nil,
		External:  nil,
	}
	result_insert_, err_insert := col_device.InsertOne(context.Background(), newdeviceinfo)
	if err_insert != nil {
		response.Code = Baseinfo.CONST_INSERT_FAIL
		response.Msg = err_insert.Error()
		return response
	}
	response.Code = Baseinfo.Success
	response.Data = result_insert_.InsertedID
	return response
}
