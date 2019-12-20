package Services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

//--创建空间--
type SpaceCreateRequest struct {
	Token    string `json:"token"`
	Level    int    `json:"level"`
	Province string `json:"province"`
	City     string `json:"city"`
	Area     string `json:"area"`
	District string `json:"distinct"`
	Building string `json:"building"`
	Storey   string `json:"storey"`
	Room     string `json:"room"`
	Place    string `json:"place"`
	Title    string `json:"title"`
}

//TODO 强转雷西不过错误的情况处理
func SpaceCreateEndpoint(spacecreateservice WSpaceCreateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*SpaceCreateRequest)
		result := spacecreateservice.NewSpace(r)
		return result, nil
	}
}

//--查询空间--
type SpaceQueryRequest struct {
	Token string `json:"token"`
	Sid   string `json:"deviceid"`
}

func SpaceQueryEndpoint(spacequeryservice WSpaceQueryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*SpaceQueryRequest)
		result := spacequeryservice.QuerySpace(r)
		return result, nil
	}
}

//--修改空间--
type SpaceReviseRequest struct {
	Title string `json:"title"`
	Token string `json:"token"`
	Sid   string `json:"sid"`
}

func SpaceReviseEndpoint(spacereviseservice WSpaceReviseService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*SpaceReviseRequest)
		result := spacereviseservice.ReviseSapce(r)
		return result, nil
	}
}

//--删除空间--
type SpaceDelRequest struct {
	Token string `json:"token"`
	Sid   string `json:"sid"`
}

func SpaceDelEndpoint(spacereviseservice WSpaceDelService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*SpaceDelRequest)
		result := spacereviseservice.DelSapce(r)
		return result, nil
	}
}

//--复制空间--
type SpaceCloneRequest struct {
	Token string `json:"token"`
	Sid   string `json:"sid"`
}

func SpaceCloneEndpoint(spacecloneservice WSpaceCloneService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*SpaceCloneRequest)
		result := spacecloneservice.CloneSpace(r)
		return result, nil
	}
}

//--
