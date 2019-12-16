package Services

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		_ = logger.Log("Create_Device_Err:", err_checktoken.Error())
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

	ctx := context.Background()
	var newdevice *Baseinfo.Device
	SessionErr := Baseinfo.Client.Database("test").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result_insert_, err_insert := col_device.InsertOne(sessionContext, newdeviceinfo)
		if err_insert != nil {
			response.Code = Baseinfo.CONST_INSERT_FAIL
			response.Msg = err_insert.Error()
			_ = logger.Log("Create_Device_Err:", err_insert.Error())
			return err_insert
		}
		err_f := col_device.FindOne(sessionContext, bson.D{{"_id", result_insert_.InsertedID}}).Decode(&newdevice)
		if err_f != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			_ = logger.Log("Create_User_Err:", "can't find recently created user!"+err_f.Error())
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_f.Error()
			return err_f
		} else {
			_ = sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})
	if SessionErr != nil {
		_ = logger.Log("Create_Device_Err:", SessionErr)
		return response
	}
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
		_ = logger.Log("Query_Device_Err:", err_checktoken.Error())
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
			_ = logger.Log("Query_Device_Err:", err_find.Error())
			return response
		}
	} else { //c传入的是deviceid
		err_find := col_device.FindOne(context.Background(), bson.M{"deviceid": deviceid}).Decode(&result)
		if result == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			_ = logger.Log("Query_Device_Err:", err_find.Error())
			return response
		}
	}
	if result.Userid == "" {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = " can't query unbound devices!"
		_ = logger.Log("Query_Device_Err:", " can't query unbound devices!")
		return response
	}
	if tokenuser != result.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = " not allowed for another user' devices!"
		_ = logger.Log("Query_Device_Err:", " not allowed for another user' devices!")
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

func (this DeviceDeleteService) DeleteDevice(r *DeviceDeleteRequest) *CommonResponse {
	var deletecount int64
	response := &CommonResponse{}
	col_device := Baseinfo.Client.Database("test").Collection("device")
	col_space := Baseinfo.Client.Database("test").Collection("space")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		_ = logger.Log("Delete_Device_Err:", err_checktoken.Error())
		return response
	}
	var dev *Baseinfo.Device
	var spaceid primitive.ObjectID
	id, err_obj := primitive.ObjectIDFromHex(r.Deviceid)

	ctx := context.Background()
	SessionErr := Baseinfo.Client.Database("test").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		if err_obj == nil {
			//传入的时id
			err_find := col_device.FindOne(sessionContext, bson.M{"_id": id}).Decode(&dev)
			if dev == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = err_find.Error()
				_ = logger.Log("Delete_Device_Err:", err_find.Error())
				return err_find
			}
			spaceid = dev.Sid
			if tokenuser != dev.Userid {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "no authority to delete another user' device !"
				_ = logger.Log("Delete_Device_Err:", "no authority to delete another user' device !")
				return errors.New("no authority to delete another user' device !")
			}
			count, err_del := col_device.DeleteOne(sessionContext, bson.M{"_id": id})
			if err_del != nil {
				response.Code = Baseinfo.CONST_DELETE_FAIL
				response.Msg = err_del.Error()
				_ = logger.Log("Delete_Device_Err:", err_del.Error())
				return err_del
			}
			deletecount = count.DeletedCount
		} else { //传入的是deviceid
			err_find := col_device.FindOne(sessionContext, bson.M{"deviceid": r.Deviceid}).Decode(&dev)
			if dev == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = err_find.Error()
				_ = logger.Log("Delete_Device_Err:", err_find.Error())
				return err_find
			}
			spaceid = dev.Sid
			if tokenuser != dev.Userid {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "no authority to delete another user' device !"
				_ = logger.Log("Delete_Device_Err:", "no authority to delete another user' device !")
				return errors.New("no authority to delete another user' device !")
			}
			count, err_del := col_device.DeleteOne(sessionContext, bson.M{"deviceid": r.Deviceid})
			if err_del != nil {
				response.Code = Baseinfo.CONST_DELETE_FAIL
				response.Msg = err_del.Error()
				_ = logger.Log("Delete_Device_Err:", err_del.Error())
				return err_del
			}
			deletecount = count.DeletedCount
		}

		//--清除sid信息
		if spaceid != primitive.NilObjectID {
			var devids []primitive.ObjectID
			var space *Baseinfo.Space
			_ = col_space.FindOne(sessionContext, bson.D{{"_id", dev.Sid}}).Decode(&space)
			if space != nil {
				for _, v := range space.Devids {
					if v != dev.Id {
						devids = append(devids, v)
					}
				}
				_, err_upd := col_space.UpdateOne(sessionContext, bson.D{{"_id", dev.Sid}}, bson.D{{"$set", bson.D{{"devids", devids}}}})
				if err_upd != nil {
					response.Code = Baseinfo.CONST_UPDATE_FAIL
					response.Msg = "fail to update coresponding space! "
					_ = sessionContext.AbortTransaction(sessionContext)
					_ = logger.Log("Delete_Device_Err: (fail to update coresponding sapce)", err_upd.Error())
					return err_upd
				}
			}
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})

	if SessionErr != nil {
		_ = logger.Log("Delete_Device_Err:", SessionErr)
		return response
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
		_ = logger.Log("Revise_Device_Err:", err_checktoken.Error())
		return response
	}
	var device *Baseinfo.Device
	var revieddev *Baseinfo.Device
	id, err_obj := primitive.ObjectIDFromHex(r.Deviceid)

	ctx := context.Background()
	SessionErr := Baseinfo.Client.Database("test").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		if err_obj == nil {
			//传入的是id
			err_find := col_device.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&device)
			if device == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = err_find.Error()
				_ = logger.Log("Revise_Device_Err:", err_find.Error())
				return err_find
			}
			if tokenuser != device.Userid {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "no authority to delete another user' device !"
				_ = logger.Log("Revise_Device_Err:", "no authority to delete another user' device !")
				return errors.New("no authority to delete another user' device")
			}

			_, err_upd := col_device.UpdateOne(context.Background(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{
				{"title", r.Title},
				{"external", r.External},
			}}})
			if err_upd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = "fail to update"
				_ = logger.Log("Revise_Device_Err:", err_upd.Error())
				return errors.New("fail to update ")
			}
			err_f := col_device.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&revieddev)
			if err_f != nil {
				_ = sessionContext.AbortTransaction(sessionContext)
				response.Msg = err_f.Error()
				response.Msg = "can't find recently revised device"
				return errors.New("can't find recently revised device")
			}
		} else {
			//传入的是devid
			err_find := col_device.FindOne(context.Background(), bson.D{{"deviceid", r.Deviceid}}).Decode(&device)
			if device == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = "find no device to revise "
				_ = logger.Log("Revise_Device_Err:", err_find.Error())
				return errors.New("find no device to revise ")
			}
			if tokenuser != device.Userid {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "no authority to delete another user' device !"
				_ = logger.Log("Revise_Device_Err:", "no authority to delete another user' device !")
				return errors.New("no authority to delete another user' device")
			}
			_, err_upd := col_device.UpdateOne(context.Background(), bson.D{{"deviceid", r.Deviceid}}, bson.D{{"$set", bson.D{
				{"title", r.Title},
				{"external", r.External},
			}}})
			if err_upd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = "fail to update device"
				_ = logger.Log("Revise_Device_Err:", err_upd.Error())
				return errors.New("fail to update device")
			}
			_ = col_device.FindOne(context.Background(), bson.D{{"deviceid", r.Deviceid}}).Decode(&revieddev)
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})
	if SessionErr != nil {
		_ = logger.Log("Revise_Device_Err:", SessionErr)
		return response
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
		_ = logger.Log("Bind_Device_Err:", err_checktoken)
		return response
	}
	deviceid := r.Deviceid
	sid := r.Sid
	gatewayid := r.Gatewayid
	userid := r.Userid

	if deviceid == "" {
		response.Code = Baseinfo.CONST_PARAM_LACK
		response.Msg = "deviceid can't be nil!"
		_ = logger.Log("Bind_Device_Err:", "deviceid can't be nil!")
		return response
	}

	if tokenuser != userid && userid != "" {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "logining user and binding user is unmathced!"
		_ = logger.Log("Bind_Device_Err:", "logining user and binding user is unmathced!")
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
			_ = logger.Log("Bind_Device_Err:", err_find)
			return response
		}
	} else {
		isid = false
		err_find := col_device.FindOne(context.Background(), bson.M{"deviceid": deviceid}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = err_find.Error()
			_ = logger.Log("Bind_Device_Err:", err_find)
			return response
		}
	}

	//绑定账号
	if userid != "" {
		if dev.Userid != "" {
			response.Code = Baseinfo.CONST_PARAM_ERROR
			response.Msg = "device has been bound!"
			_ = logger.Log("Bind_Device_Err:", "device has been bound!")
			return response
		}
		if isid {
			_, err_upd := col_device.UpdateOne(context.Background(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"userid", userid}}}})
			if err_upd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err_upd.Error()
				_ = logger.Log("Bind_Device_Err:", err_upd.Error())
				return response
			}
		} else {
			_, err_upd := col_device.UpdateOne(context.Background(), bson.D{{"deviceid", deviceid}}, bson.D{{"$set", bson.D{{"userid", userid}}}})
			if err_upd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err_upd.Error()
				_ = logger.Log("Bind_Device_Err:", err_upd.Error())
				return response
			}
		}
	}
	//sid为真实有效的id时，绑定房源
	sid_obj, err_sid := primitive.ObjectIDFromHex(sid)
	if err_sid == nil {
		var devids []primitive.ObjectID
		var space *Baseinfo.Space
		_ = col_space.FindOne(context.Background(), bson.D{{"_id", sid_obj}}).Decode(&space)
		if space == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "find no space by sid!" //包含房源不存在的情况
			_ = logger.Log("Bind_Device_Err:", "find no space by sid!")
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
			_ = logger.Log("Bind_Device_Err:", err_upd)
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
				_ = logger.Log("Bind_Device_Err:", err_update)
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
				_ = logger.Log("Bind_Device_Err:", err1)
				return response
			}
		}
	} else {
		response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
		response.Msg = "invalid spaceid"
		_ = logger.Log("Bind_Device_Err:", "invalid spaceid")
		return response
	}
	if gatewayid != "" { //gatewayid只允许输入编号，不是_id
		var gateway *Baseinfo.Device
		_ = col_device.FindOne(context.Background(), bson.D{{"deviceid", gatewayid}}).Decode(&gateway)
		if gateway == nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = "can't find gateway"
			_ = logger.Log("Bind_Device_Err:", "can't find gateway")
			return response
		}
		if gateway.Userid != tokenuser {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "gateway has been bound in another user or not bound"
			_ = logger.Log("Bind_Device_Err:", "gateway has been bound in another user or not bound")
			return response
		}
		if isid {
			_, err := col_device.UpdateOne(context.Background(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"gatewayid", gatewayid}}}})
			if err != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err.Error()
				_ = logger.Log("Bind_Device_Err:", err.Error())
				return response
			}
		} else {
			_, err := col_device.UpdateOne(context.Background(), bson.D{{"deviceid", deviceid}}, bson.D{{"$set", bson.D{{"gatewayid", gatewayid}}}})
			if err != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = err.Error()
				_ = logger.Log("Bind_Device_Err:", err.Error())
				return response
			}
		}
	}

	var newdev *Baseinfo.Device
	if isid {
		_ = col_device.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&newdev)
	} else {
		_ = col_device.FindOne(context.Background(), bson.D{{"deviceid", deviceid}}).Decode(&newdev)
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
		_ = logger.Log("unbound device err:", err_checktoken)
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
		_ = col_device.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "find no device!"
			_ = logger.Log("unbound device err:", "find no device!")
			return response
		}
		if tokenuser != dev.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't operate another user' device!"
			_ = logger.Log("unbound device err:", "cant't operate another user' device!")
			return response
		}
		var node *Baseinfo.Device
		_ = col_device.FindOne(context.Background(), bson.D{{"gatewayid", dev.Deviceid}, {"isnode", true}}).Decode(&node)
		if node != nil {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't unbound gateway with nodes under it!"
			_ = logger.Log("unbound device err:", "cant't unbound gateway with nodes under it!")
			return response
		}
		errcode, errmsg, data = Baseinfo.UnboundDeviceByid(dev, kind, col_device, col_space)

	} else {
		var dev *Baseinfo.Device
		_ = col_device.FindOne(context.Background(), bson.D{{"deviceid", deviceid}}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "find no device!"
			_ = logger.Log("unbound device err:", "find no device!")
			return response
		}
		if tokenuser != dev.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't operate another user' device!"
			_ = logger.Log("unbound device err:", "cant't operate another user' device!")
			return response
		}
		var node *Baseinfo.Device
		_ = col_device.FindOne(context.Background(), bson.D{{"gatewayid", deviceid}, {"isnode", true}}).Decode(&node)
		if node != nil {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't unbound gateway with nodes under it!"
			_ = logger.Log("unbound device err:", "cant't unbound gateway with nodes under it!")
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
		response.Code = Baseinfo.CONST_INSERT_FAIL
		response.Msg = err_ins.Error()
		_ = logger.Log("upload data err:", err_ins)
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
			_ = logger.Log("upload data err:", updateresult)
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
			_ = logger.Log("upload data err:", updateresult.Err().Error())
			return response
		}
	}
	response.Code = Baseinfo.Success
	return response
}
