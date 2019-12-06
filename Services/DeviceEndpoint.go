package Services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

//--新建设备--
type DeviceCreateRequest struct {
	Token    string `json:"token"`
	Deviceid string `json:"deviceid"`
	Isnode   bool   `json:"isnode"`
	Devtype  int64  `json:"devtype"`
	Title    string `json:"title"`
}

func DeviceCreateEndpoint(userCreateService WDeviceCreateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) { //TODO request数据内容哪来的
		r := request.(*DeviceCreateRequest)
		result := userCreateService.NewDevice(r)
		return result, nil
	}
}

//--查询设备--

type DeviceQueryRequest struct {
	Token    string `json:"token"`
	Deviceid string `json:"deviceid"`
}

func DeviceQueryEndpoint(deviceQueryService WDeviceQueryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*DeviceQueryRequest)
		result := deviceQueryService.QUeryDevice(r)
		return result, nil
	}
}

//--删除设备--
type DeviceDeleteRequest struct {
	Token    string `json:"token"`
	Deviceid string `json:"devieid"`
}

func DeviceDeleteEndpoint(devicedeleteservice WDeviceDeleteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) { //TODO request数据内容哪来的
		r := request.(*DeviceDeleteRequest)
		result := devicedeleteservice.DeleteDevice(r)
		return result, nil
	}
}

//--修改设备--
type DeviceReviseRequest struct {
	Token    string                 `json:"token"`
	Title    string                 `json:"title"`
	Deviceid string                 `json:"deviceid"`
	Expand   map[string]interface{} `json:"expand"`
}

func DeviceReviseEndpoint(devicereviseservice WDeviceReviseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*DeviceReviseRequest)
		result := devicereviseservice.ReviseDevice(r)
		return result, nil
	}
}
