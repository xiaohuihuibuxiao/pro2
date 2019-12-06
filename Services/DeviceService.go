package Services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pro2/Baseinfo"
	"time"
)

//--新建设备--
type WDeviceCreateService interface {
	NewDevice(r *DeviceCreateRequest) *CommonResponse
}

type DeviceCreateService struct{}

func (this DeviceCreateService) NewDevice(r *DeviceCreateRequest) *CommonResponse {
	col_device := Baseinfo.Client.Database("test").Collection("device")
	response := &CommonResponse{}

	//TODO how to get token
	token := r.Token
	err_checktoken, _ := Baseinfo.Logintokenauth(token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}

	newdeviceinfo := &Baseinfo.Device{
		Id:        primitive.NewObjectIDFromTimestamp(time.Now()),
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

//--查询设备--

type WDeviceQueryService interface {
	QUeryDevice(r *DeviceQueryRequest) *CommonResponse
}
type DeviceQUeryService struct{}

func (this DeviceQUeryService) QUeryDevice(r *DeviceQueryRequest) *CommonResponse {
	fmt.Println("进入查询设备")
	response := &CommonResponse{}
	col_device := Baseinfo.Client.Database("test").Collection("device")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}
	var result *Baseinfo.Device
	deviceid := r.Deviceid
	Id, err_obj := primitive.ObjectIDFromHex(deviceid)
	if err_obj == nil { //传入的是id
		err_find := col_device.FindOne(context.Background(), bson.M{"_id": Id}).Decode(&result)
		if result == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
	} else { //c传入的是deviceid
		err_find := col_device.FindOne(context.Background(), bson.M{"deviceid": deviceid}).Decode(&result)
		if result == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
	}
	if result.Userid == "" {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = " can't query unbound devices!"
		return response
	}
	if tokenuser != result.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = " not allowed for another user' devices!"
		return response
	}
	response.Code = Baseinfo.Success
	response.Data = result
	return response
}

//--删除设备--
type WDeviceDeleteService interface {
	DeleteDevice(r *DeviceDeleteRequest) *CommonResponse
}
type DeviceDeleteService struct{}

//不考虑网关的删除和相关影响 TODO 清除sid信息暂未验证
func (this DeviceDeleteService) DeleteDevice(r *DeviceDeleteRequest) *CommonResponse {
	var deletecount int64
	response := &CommonResponse{}
	col_device := Baseinfo.Client.Database("test").Collection("device")
	col_space := Baseinfo.Client.Database("test").Collection("space")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}
	var dev *Baseinfo.Device
	id, err_obj := primitive.ObjectIDFromHex(r.Deviceid)
	if err_obj == nil {
		//传入的时id
		err_find := col_device.FindOne(context.Background(), bson.M{"_id": id}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}

		if tokenuser != dev.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "no authority to delete another user' device !"
			return response
		}
		count, err_del := col_device.DeleteOne(context.Background(), bson.M{"_id": id})
		if err_del != nil {
			response.Code = Baseinfo.CONST_DELETE_FAIL
			response.Msg = err_del.Error()
			return response
		}
		deletecount = count.DeletedCount
	} else { //传入的是deviceid
		err_find := col_device.FindOne(context.Background(), bson.M{"deviceid": r.Deviceid}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
		if tokenuser != dev.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "no authority to delete another user' device !"
			return response
		}
		count, err_del := col_device.DeleteOne(context.Background(), bson.M{"deviceid": r.Deviceid})
		if err_del != nil {
			response.Code = Baseinfo.CONST_DELETE_FAIL
			response.Msg = err_del.Error()
			return response
		}
		deletecount = count.DeletedCount
	}

	//--清除sid信息
	if dev.Sid != primitive.NilObjectID {
		var devids []primitive.ObjectID
		var space *Baseinfo.Space
		err_find := col_space.FindOne(context.Background(), bson.D{{"_id", dev.Sid}}).Decode(&space)
		if space == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
		for _, v := range space.Devids {
			if v != dev.Id {
				devids = append(devids, v)
			}
		}
		_, err_upd := col_space.UpdateOne(context.Background(), bson.D{{"_id", dev.Sid}}, bson.D{{"$set", bson.D{{"devids", devids}}}})
		if err_upd != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = err_find.Error()
			return response
		}
	}

	response.Code = Baseinfo.Success
	response.Data = deletecount
	return response
}

//--修改设备--
type WDeviceReviseService interface {
	ReviseDevice(r *DeviceReviseRequest) *CommonResponse
}
type DeviceReviseService struct{}

func (this DeviceReviseService) ReviseDevice(r *DeviceReviseRequest) *CommonResponse {
	fmt.Println("进入修改设备")
	response := &CommonResponse{}
	col_device := Baseinfo.Client.Database("test").Collection("device")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}
	var device *Baseinfo.Device
	id, err_obj := primitive.ObjectIDFromHex(r.Deviceid)
	if err_obj == nil {
		//传入的是id
		err_find := col_device.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&device)
		if device == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
		if tokenuser != device.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "no authority to delete another user' device !"
			return response
		}
		_, err_upd := col_device.UpdateOne(context.Background(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{
			{"title", r.Title},
			{"expand", r.Expand},
		}}})
		if err_upd != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = err_upd.Error()
			return response
		}
	} else {
		//传入的是devid
		err_find := col_device.FindOne(context.Background(), bson.D{{"deviceid", r.Deviceid}}).Decode(&device)
		if device == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
		if tokenuser != device.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "no authority to delete another user' device !"
			return response
		}
		_, err_upd := col_device.UpdateOne(context.Background(), bson.D{{"deviceid", r.Deviceid}}, bson.D{{"$set", bson.D{
			{"title", r.Title},
			{"expand", r.Expand},
		}}})
		if err_upd != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = err_upd.Error()
			return response
		}
	}
	response.Code = Baseinfo.Success
	return response
}
