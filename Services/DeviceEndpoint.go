package Services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

//--新建设备--
type DeviceCreateRequest struct {
	Deviceid string `json:"deviceid"`
	Isnode   bool   `json:"isnode"`
	Devtype  int64  `json:"devtype"`
	Title    string `json:"title"`
}

func DeviceCreateEndpoint(userCreateService WDeviceCreateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) { //TODO requesr数据内容哪来的
		r := request.(*DeviceCteateRequest)
		result := userCreateService.NewDevice(r)
		return result, nil
	}
}
