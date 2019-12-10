package Services

import (
	"context"
	"github.com/astaxie/beego/logs"
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
		Spacecode: "",
		Expand:    nil,
		External:  nil,
	}
	result_insert_, err_insert := col_device.InsertOne(context.Background(), newdeviceinfo)
	if err_insert != nil {
		response.Code = Baseinfo.CONST_INSERT_FAIL
		response.Msg = err_insert.Error()
		return response
	}
	var newdevice *Baseinfo.Device
	col_device.FindOne(context.Background(), bson.D{{"_id", result_insert_.InsertedID}}).Decode(&newdevice)
	response.Code = Baseinfo.Success
	response.Data = newdevice
	return response
}

//--查询设备--

type WDeviceQueryService interface {
	QUeryDevice(r *DeviceQueryRequest) *CommonResponse
}
type DeviceQUeryService struct{}

func (this DeviceQUeryService) QUeryDevice(r *DeviceQueryRequest) *CommonResponse {
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

//不考虑网关的删除和相关影响
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
	var spaceid primitive.ObjectID
	id, err_obj := primitive.ObjectIDFromHex(r.Deviceid)
	if err_obj == nil {
		//传入的时id
		err_find := col_device.FindOne(context.Background(), bson.M{"_id": id}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
		spaceid = dev.Sid
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
		spaceid = dev.Sid
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
	if spaceid != primitive.NilObjectID {
		var devids []primitive.ObjectID
		var space *Baseinfo.Space
		err_find := col_space.FindOne(context.Background(), bson.D{{"_id", dev.Sid}}).Decode(&space)
		if space != nil {
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
	response := &CommonResponse{}
	col_device := Baseinfo.Client.Database("test").Collection("device")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}
	var device *Baseinfo.Device
	var revieddev *Baseinfo.Device
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
			{"external", r.External},
		}}})
		if err_upd != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = err_upd.Error()
			return response
		}
		col_device.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&revieddev)
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
			{"external", r.External},
		}}})
		if err_upd != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = err_upd.Error()
			return response
		}
		col_device.FindOne(context.Background(), bson.D{{"deviceid", r.Deviceid}}).Decode(&revieddev)
	}
	response.Code = Baseinfo.Success
	response.Data = revieddev
	return response
}

//--绑定设备--

type WDeviceBindService interface {
	BindDevice(r *DeviceBindRequest) *CommonResponse
}
type DeviceBindService struct{}

func (this DeviceBindService) BindDevice(r *DeviceBindRequest) *CommonResponse {
	response := &CommonResponse{}
	col_device := Baseinfo.Client.Database("test").Collection("device")
	col_space := Baseinfo.Client.Database("test").Collection("space")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}
	deviceid := r.Deviceid
	sid := r.Sid
	gatewayid := r.Gatewayid
	userid := r.Userid

	if deviceid == "" {
		response.Code = Baseinfo.CONST_PARAM_LACK
		response.Msg = "deviceid can't be nil!"
		return response
	}

	if tokenuser != userid && userid != "" {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "logining user and binding user is unmathced!"
		return response
	}

	var dev *Baseinfo.Device
	id, err_obj := primitive.ObjectIDFromHex(deviceid)
	var isid bool
	if err_obj == nil {
		isid = true
		err_find := col_device.FindOne(context.Background(), bson.M{"_id": id}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
	} else {
		isid = false
		err_find := col_device.FindOne(context.Background(), bson.M{"deviceid": deviceid}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			return response
		}
	}

	//绑定账号
	if userid != "" {
		if dev.Userid != "" {
			response.Code = Baseinfo.CONST_PARAM_ERROR
			response.Msg = "device has been bound!"
			return response
		}
		if isid {
			_, err_upd := col_device.UpdateOne(context.Background(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"userid", userid}}}})
			if err_upd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err_upd.Error()
				return response
			}
		} else {
			_, err_upd := col_device.UpdateOne(context.Background(), bson.D{{"deviceid", deviceid}}, bson.D{{"$set", bson.D{{"userid", userid}}}})
			if err_upd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err_upd.Error()
				return response
			}
		}
	}
	//sid为真实有效的id时，绑定房源
	sid_obj, err_sid := primitive.ObjectIDFromHex(sid)
	if err_sid == nil {
		var devids []primitive.ObjectID
		var space *Baseinfo.Space
		col_space.FindOne(context.Background(), bson.D{{"_id", sid_obj}}).Decode(&space)
		if space == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "find no space by sid!" //包含房源不存在的情况
			return response
		}
		devids = space.Devids
		if isid {
			devids = append(devids, id)
		} else {
			devids = append(devids, dev.Id)
		}
		//更新space表
		_, err_upd := col_space.UpdateOne(context.Background(), bson.D{{"_id", sid_obj}}, bson.D{{"$set", bson.D{{"devids", devids}}}})
		if err_upd != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = err_upd.Error()
			return response
		}
		//跟新device表
		if isid {
			_, err_update := col_device.UpdateOne(context.Background(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{
				{"sid", sid_obj},
				{"addr", space.Addr},
				{"spacecode", space.Spacecode},
			}}})
			if err_update != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err_update.Error()
				return response
			}
		} else {
			_, err1 := col_device.UpdateOne(context.Background(), bson.D{{"deviceid", deviceid}}, bson.D{{"$set", bson.D{
				{"sid", sid_obj},
				{"addr", space.Addr},
				{"spacecode", space.Spacecode},
			}}})
			if err1 != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err1.Error()
				return response
			}
		}
	} else {
		response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
		response.Msg = "invalid spaceid"
		return response
	}
	if gatewayid != "" { //gatewayid只允许输入编号，不是_id
		var gateway *Baseinfo.Device
		col_device.FindOne(context.Background(), bson.D{{"deviceid", gatewayid}}).Decode(&gateway)
		if gateway == nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = "can't find gateway"
			return response
		}
		if gateway.Userid != tokenuser {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "gateway has been bound in another user or not bound"
			return response
		}
		if isid {
			_, err := col_device.UpdateOne(context.Background(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"gatewayid", gatewayid}}}})
			if err != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err.Error()
				return response
			}
		} else {
			_, err := col_device.UpdateOne(context.Background(), bson.D{{"deviceid", deviceid}}, bson.D{{"$set", bson.D{{"gatewayid", gatewayid}}}})
			if err != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err.Error()
				return response
			}
		}
	}

	var newdev *Baseinfo.Device
	if isid {
		col_device.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&newdev)
	} else {
		col_device.FindOne(context.Background(), bson.D{{"deviceid", deviceid}}).Decode(&newdev)
	}
	response.Code = Baseinfo.Success
	response.Data = newdev
	return response
}

//--解绑设备--
type WDeviceUnboundService interface {
	UnboundDevice(r *DeviceUnboundRequest) *CommonResponse
}
type DeviceUnboundService struct{}

func (this DeviceUnboundService) UnboundDevice(r *DeviceUnboundRequest) *CommonResponse {
	response := &CommonResponse{}
	col_device := Baseinfo.Client.Database("test").Collection("device")
	col_space := Baseinfo.Client.Database("test").Collection("space")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}
	deviceid := r.Deviceid
	kind := r.Type

	var errcode int64
	var errmsg string
	var data interface{}
	id, err_obj := primitive.ObjectIDFromHex(deviceid)
	if err_obj == nil {
		var dev *Baseinfo.Device
		col_device.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "find no device!"
			return response
		}
		if tokenuser != dev.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't operate another user' device!"
			return response
		}
		var node *Baseinfo.Device
		col_device.FindOne(context.Background(), bson.D{{"gatewayid", dev.Deviceid}, {"isnode", true}}).Decode(&node)
		if node != nil {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't unbound gateway with nodes under it!"
			return response
		}
		errcode, errmsg, data = Baseinfo.UnboundDeviceByid(dev, kind, col_device, col_space)

	} else {
		var dev *Baseinfo.Device
		col_device.FindOne(context.Background(), bson.D{{"deviceid", deviceid}}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "find no device!"
			return response
		}
		if tokenuser != dev.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't operate another user' device!"
			return response
		}
		var node *Baseinfo.Device
		col_device.FindOne(context.Background(), bson.D{{"gatewayid", deviceid}, {"isnode", true}}).Decode(&node)
		if node != nil {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't unbound gateway with nodes under it!"
			return response
		}
		errcode, errmsg, data = Baseinfo.UnboundDeviceBydeviceid(dev, kind, col_device, col_space)
	}
	response.Code = errcode
	response.Msg = errmsg
	response.Data = data
	return response
}

//--上报数据--
type WDeviceUploadService interface {
	UploadData(r *DeviceUploadRequest) *CommonResponse
}
type DeviceUploadService struct{}

func (this DeviceUploadService) UploadData(r *DeviceUploadRequest) *CommonResponse {
	response := &CommonResponse{}
	col_device := Baseinfo.Client.Database("test").Collection("device")
	col_history := Baseinfo.Client.Database("test").Collection("history")

	deviceid := r.Deviceid

	history := &Baseinfo.Sensorhistory{
		Id:       primitive.NewObjectIDFromTimestamp(time.Now()),
		Userid:   r.Userid,
		Deviceid: r.Deviceid,
		Devtype:  r.Devtype,
		//UPloadtime:    time.Now().Format("2006-01-02 15:04:05"),
		UPloadtime: r.T,
		Expand:     r.Data,
		External:   nil,
	}
	//记录至history表
	_, err_ins := col_history.InsertOne(context.Background(), history)
	if err_ins != nil {
		logs.Info("fail to inset to history ")
		response.Code = Baseinfo.CONST_INSERT_FAIL
		response.Msg = err_ins.Error()
		return response
	}

	//更新device表的expand
	id, err_obj := primitive.ObjectIDFromHex(deviceid)
	if err_obj == nil {
		filter := bson.D{{"_id", id}}
		update := bson.D{{"$set", bson.D{{"expand", &struct {
			T    string
			Data interface{}
		}{
			T:    r.T,
			Data: r.Data,
		},
		}}}}
		updateresult := col_device.FindOneAndUpdate(context.Background(), filter, update)
		if updateresult.Err() != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = updateresult.Err().Error()
			return response
		}
	} else {
		filter := bson.D{{"deviceid", deviceid}}
		update := bson.D{{"$set", bson.D{{"expand", &struct {
			T    string
			Data interface{}
		}{
			T:    r.T,
			Data: r.Data,
		},
		}}}}
		updateresult := col_device.FindOneAndUpdate(context.Background(), filter, update)
		if updateresult.Err() != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = updateresult.Err().Error()
			return response
		}
	}
	response.Code = Baseinfo.Success
	return response
}
