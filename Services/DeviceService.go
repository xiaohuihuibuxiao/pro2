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
	errChecktoken, _ := Baseinfo.Logintokenauth(token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("CreateDeviceErr:", errChecktoken.Error())
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
	sessionErr := Baseinfo.Client.Database("test").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		resultInsert, errInsert := col_device.InsertOne(sessionContext, newdeviceinfo)
		if errInsert != nil {
			response.Code = Baseinfo.CONST_INSERT_FAIL
			response.Msg = errInsert.Error()
			_ = logger.Log("CreateDeviceErr:", errInsert.Error())
			return errInsert
		}
		ErrF := col_device.FindOne(sessionContext, bson.D{{"_id", resultInsert.InsertedID}}).Decode(&newdevice)
		if ErrF != nil {
			_ = sessionContext.AbortTransaction(sessionContext)
			_ = logger.Log("CreateUserErr:", "can't find recently created user!"+ErrF.Error())
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = ErrF.Error()
			return ErrF
		} else {
			_ = sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})
	if sessionErr != nil {
		_ = logger.Log("CreateDeviceErr:", sessionErr)
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
	colDevice := Baseinfo.Client.Database("test").Collection("device")

	errChecktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("QueryDeviceErr:", errChecktoken.Error())
		return response
	}
	var result *Baseinfo.Device
	deviceid := r.Deviceid
	Id, errObj := primitive.ObjectIDFromHex(deviceid)
	if errObj == nil { //传入的是id
		errFind := colDevice.FindOne(context.Background(), bson.M{"_id": Id}).Decode(&result)
		if result == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = errFind.Error()
			_ = logger.Log("Query_Device_Err:", errFind.Error())
			return response
		}
	} else { //c传入的是deviceid
		errFind := colDevice.FindOne(context.Background(), bson.M{"deviceid": deviceid}).Decode(&result)
		if result == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = errFind.Error()
			_ = logger.Log("Query_Device_Err:", errFind.Error())
			return response
		}
	}
	if result.Userid == "" {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = " can't query unbound devices!"
		_ = logger.Log("QueryDeviceErr:", " can't query unbound devices!")
		return response
	}
	if tokenuser != result.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = " not allowed for another user' devices!"
		_ = logger.Log("QueryDeviceErr:", " not allowed for another user' devices!")
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
	colDevice := Baseinfo.Client.Database("test").Collection("device")
	colSpace := Baseinfo.Client.Database("test").Collection("space")

	errChecktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("DeleteDeviceErr:", errChecktoken.Error())
		return response
	}
	var dev *Baseinfo.Device
	var spaceid primitive.ObjectID
	id, errObj := primitive.ObjectIDFromHex(r.Deviceid)

	ctx := context.Background()
	sessionErr := Baseinfo.Client.Database("test").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		if errObj == nil {
			//传入的时id
			errFind := colDevice.FindOne(sessionContext, bson.M{"_id": id}).Decode(&dev)
			if dev == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = errFind.Error()
				_ = logger.Log("DeleteDeviceErr:", errFind.Error())
				return errFind
			}
			spaceid = dev.Sid
			if tokenuser != dev.Userid {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "no authority to delete another user' device !"
				_ = logger.Log("DeleteDeviceErr:", "no authority to delete another user' device !")
				return errors.New("no authority to delete another user' device")
			}
			count, errDel := colDevice.DeleteOne(sessionContext, bson.M{"_id": id})
			if errDel != nil {
				response.Code = Baseinfo.CONST_DELETE_FAIL
				response.Msg = errDel.Error()
				_ = logger.Log("Delete_Device_Err:", errDel.Error())
				return errDel
			}
			deletecount = count.DeletedCount
		} else { //传入的是deviceid
			errFind := colDevice.FindOne(sessionContext, bson.M{"deviceid": r.Deviceid}).Decode(&dev)
			if dev == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = errFind.Error()
				_ = logger.Log("DeleteDeviceErr:", errFind.Error())
				return errFind
			}
			spaceid = dev.Sid
			if tokenuser != dev.Userid {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "no authority to delete another user' device !"
				_ = logger.Log("Delete_Device_Err:", "no authority to delete another user' device !")
				return errors.New("no authority to delete another user' device")
			}
			count, errDel := colDevice.DeleteOne(sessionContext, bson.M{"deviceid": r.Deviceid})
			if errDel != nil {
				response.Code = Baseinfo.CONST_DELETE_FAIL
				response.Msg = errDel.Error()
				_ = logger.Log("DeleteDeviceErr:", errDel.Error())
				return errDel
			}
			deletecount = count.DeletedCount
		}

		//--清除sid信息
		if spaceid != primitive.NilObjectID {
			var devids []primitive.ObjectID
			var space *Baseinfo.Space
			_ = colSpace.FindOne(sessionContext, bson.D{{"_id", dev.Sid}}).Decode(&space)
			if space != nil {
				for _, v := range space.Devids {
					if v != dev.Id {
						devids = append(devids, v)
					}
				}
				_, errUpd := colSpace.UpdateOne(sessionContext, bson.D{{"_id", dev.Sid}}, bson.D{{"$set", bson.D{{"devids", devids}}}})
				if errUpd != nil {
					response.Code = Baseinfo.CONST_UPDATE_FAIL
					response.Msg = "fail to update corresponding space! "
					_ = sessionContext.AbortTransaction(sessionContext)
					_ = logger.Log("Delete_Device_Err: (fail to update corresponding sapce)", errUpd.Error())
					return errUpd
				}
			}
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})

	if sessionErr != nil {
		_ = logger.Log("DeleteDeviceErr:", sessionErr)
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
	colDevice := Baseinfo.Client.Database("test").Collection("device")

	errChecktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("Revise_Device_Err:", errChecktoken.Error())
		return response
	}
	var device *Baseinfo.Device
	var revieddev *Baseinfo.Device
	id, errObj := primitive.ObjectIDFromHex(r.Deviceid)

	ctx := context.Background()
	sessionErr := Baseinfo.Client.Database("test").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		if errObj == nil {
			//传入的是id
			errFind := colDevice.FindOne(sessionContext, bson.D{{"_id", id}}).Decode(&device)
			if device == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = errFind.Error()
				_ = logger.Log("ReviseDeviceErr:", errFind.Error())
				return errFind
			}
			if tokenuser != device.Userid {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "no authority to delete another user' device !"
				_ = logger.Log("Revise_Device_Err:", "no authority to delete another user' device !")
				return errors.New("no authority to delete another user' device")
			}

			_, errUpd := colDevice.UpdateOne(sessionContext, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{
				{"title", r.Title},
				{"external", r.External},
			}}})
			if errUpd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = "fail to update"
				_ = logger.Log("Revise_Device_Err:", errUpd.Error())
				return errors.New("fail to update")
			}
			errF := colDevice.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&revieddev)
			if errF != nil {
				_ = sessionContext.AbortTransaction(sessionContext)
				response.Msg = errF.Error()
				response.Msg = "can't find recently revised device"
				return errors.New("can't find recently revised device")
			}
		} else {
			//传入的是devid
			errFind := colDevice.FindOne(sessionContext, bson.D{{"deviceid", r.Deviceid}}).Decode(&device)
			if device == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = "find no device to revise "
				_ = logger.Log("Revise_Device_Err:", errFind.Error())
				return errors.New("find no device to revise ")
			}
			if tokenuser != device.Userid {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "no authority to delete another user' device !"
				_ = logger.Log("Revise_Device_Err:", "no authority to delete another user' device !")
				return errors.New("no authority to delete another user' device")
			}
			_, errUpd := colDevice.UpdateOne(sessionContext, bson.D{{"deviceid", r.Deviceid}}, bson.D{{"$set", bson.D{
				{"title", r.Title},
				{"external", r.External},
			}}})
			if errUpd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = "fail to update device"
				_ = logger.Log("ReviseDeviceErr:", errUpd.Error())
				return errors.New("fail to update device")
			}
			_ = colDevice.FindOne(sessionContext, bson.D{{"deviceid", r.Deviceid}}).Decode(&revieddev)
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})
	if sessionErr != nil {
		_ = logger.Log("ReviseDeviceErr:", sessionErr)
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
	colDevice := Baseinfo.Client.Database("test").Collection("device")
	colSpace := Baseinfo.Client.Database("test").Collection("space")

	errChecktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("BindDeviceErr:", errChecktoken)
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
		_ = logger.Log("BindDeviceErr:", "logining user and binding user is unmathced!")
		return response
	}

	var dev *Baseinfo.Device
	id, errObj := primitive.ObjectIDFromHex(deviceid)
	var isid bool
	if errObj == nil {
		isid = true
		errFind := colDevice.FindOne(context.Background(), bson.M{"_id": id}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = errFind.Error()
			_ = logger.Log("Bind_Device_Err:", errFind.Error())
			return response
		}
	} else {
		isid = false
		errFind := colDevice.FindOne(context.Background(), bson.M{"deviceid": deviceid}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = errFind.Error()
			_ = logger.Log("BindDeviceErr:", errFind)
			return response
		}
	}

	var newdev *Baseinfo.Device
	ctx := context.Background()
	sessionErr := Baseinfo.Client.Database("test").Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		//绑定账号
		if userid != "" {
			if dev.Userid != "" {
				response.Code = Baseinfo.CONST_PARAM_ERROR
				response.Msg = "device has been bound!"
				_ = logger.Log("BindDeviceErr:", "device has been bound!")
				return errors.New("device has been bound")
			}
			if isid {
				_, err_upd := colDevice.UpdateOne(context.Background(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"userid", userid}}}})
				if err_upd != nil {
					response.Code = Baseinfo.CONST_UPDATE_FAIL
					response.Msg = "fail to update userid"
					_ = logger.Log("Bind_Device_Err:", err_upd.Error())
					return errors.New("fail to update userid")
				}
			} else {
				_, errUpd := colDevice.UpdateOne(context.Background(), bson.D{{"deviceid", deviceid}}, bson.D{{"$set", bson.D{{"userid", userid}}}})
				if errUpd != nil {
					response.Code = Baseinfo.CONST_UPDATE_FAIL
					response.Msg = "fail to update userid"
					_ = logger.Log("BindDeviceErr:", errUpd.Error())
					return errors.New("fail to update userid")
				}
			}
		}

		//sid为真实有效的id时，绑定房源
		sid_obj, errSid := primitive.ObjectIDFromHex(sid)
		if errSid == nil {
			var devids []primitive.ObjectID
			var space *Baseinfo.Space
			_ = colSpace.FindOne(sessionContext, bson.D{{"_id", sid_obj}}).Decode(&space)
			if space == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = "can't find space by sid!" //包含房源不存在的情况
				_ = sessionContext.AbortTransaction(sessionContext)
				_ = logger.Log("BindDeviceErr:", "find no space by sid!")
				return errors.New("find no space by sid")
			}
			devids = space.Devids
			if isid {
				devids = append(devids, id)
			} else {
				devids = append(devids, dev.Id)
			}
			//更新space表
			_, errUpd := colSpace.UpdateOne(sessionContext, bson.D{{"_id", sid_obj}}, bson.D{{"$set", bson.D{{"devids", devids}}}})
			if errUpd != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = "fail to update space info"
				_ = sessionContext.AbortTransaction(sessionContext)
				_ = logger.Log("BindDeviceErr:", errUpd)
				return errors.New("fail to update space info")
			}
			//跟新device表
			if isid {
				_, errUpdate := colDevice.UpdateOne(sessionContext, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{
					{"sid", sid_obj},
					{"addr", space.Addr},
					{"spacecode", space.Spacecode},
				}}})
				if errUpdate != nil {
					response.Code = Baseinfo.CONST_UPDATE_FAIL
					response.Msg = "fail to update device'sid"
					_ = sessionContext.AbortTransaction(sessionContext)
					_ = logger.Log("BindDeviceErr:", errUpdate)
					return errors.New("fail to update device'sid")
				}
			} else {
				_, err1 := colDevice.UpdateOne(sessionContext, bson.D{{"deviceid", deviceid}}, bson.D{{"$set", bson.D{
					{"sid", sid_obj},
					{"addr", space.Addr},
					{"spacecode", space.Spacecode},
				}}})
				if err1 != nil {
					response.Code = Baseinfo.CONST_UPDATE_FAIL
					response.Msg = "fail to update device'sid"
					_ = sessionContext.AbortTransaction(sessionContext)
					_ = logger.Log("BindDeviceErr:", err1)
					return errors.New("fail to update device'sid")
				}
			}
		} else {
			response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
			response.Msg = "invalid spaceid"
			_ = sessionContext.AbortTransaction(sessionContext)
			_ = logger.Log("BindDeviceErr:", errSid.Error())
			return errors.New("invalid spaceid")
		}

		// 绑定网关
		if gatewayid != "" { //gatewayid只允许输入编号，不是_id
			var gateway *Baseinfo.Device
			_ = colDevice.FindOne(sessionContext, bson.D{{"deviceid", gatewayid}}).Decode(&gateway)
			if gateway == nil {
				response.Code = Baseinfo.CONST_FIND_FAIL
				response.Msg = "can't find gateway"
				_ = sessionContext.AbortTransaction(sessionContext)
				_ = logger.Log("BindDeviceErr:", "can't find gateway")
				return errors.New("can't find gateway")
			}
			if gateway.Userid != tokenuser {
				response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
				response.Msg = "gateway has been bound in another user or not bound"
				_ = sessionContext.AbortTransaction(sessionContext)
				_ = logger.Log("Bind_Device_Err:", "gateway has been bound in another user or not bound")
				return errors.New("gateway has been bound in another user or not bound")
			}
			if isid {
				_, err := colDevice.UpdateOne(sessionContext, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"gatewayid", gatewayid}}}})
				if err != nil {
					response.Code = Baseinfo.CONST_UPDATE_FAIL
					response.Msg = "fail to bind gateway for device"
					_ = sessionContext.AbortTransaction(sessionContext)
					_ = logger.Log("Bind_Device_Err:", err.Error())
					return errors.New("fail to bind gateway for device")
				}
			} else {
				_, err := colDevice.UpdateOne(sessionContext, bson.D{{"deviceid", deviceid}}, bson.D{{"$set", bson.D{{"gatewayid", gatewayid}}}})
				if err != nil {
					response.Code = Baseinfo.CONST_UPDATE_FAIL
					response.Msg = "fail to bind gateway for device"
					_ = sessionContext.AbortTransaction(sessionContext)
					_ = logger.Log("Bind_Device_Err:", err.Error())
					return errors.New("fail to bind gateway for device")
				}
			}
		}

		var e error
		if isid {
			e = colDevice.FindOne(sessionContext, bson.D{{"_id", id}}).Decode(&newdev)
		} else {
			e = colDevice.FindOne(sessionContext, bson.D{{"deviceid", deviceid}}).Decode(&newdev)
		}
		if e != nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "can't find recently bound device"
			_ = sessionContext.AbortTransaction(sessionContext)
			_ = logger.Log("BindDeviceErr:", e.Error())
			return errors.New("can't finf recently bound device")
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil

	})
	if sessionErr != nil {
		_ = logger.Log("BindDeviceErr:", sessionErr)
		return response
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
	colDevice := Baseinfo.Client.Database("test").Collection("device")
	colSpace := Baseinfo.Client.Database("test").Collection("space")

	errChecktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("UnboundDeviceErr:", errChecktoken)
		return response
	}
	deviceid := r.Deviceid
	kind := r.Type

	var errcode int64
	var errmsg string
	var data interface{}
	id, errObj := primitive.ObjectIDFromHex(deviceid)

	if errObj == nil {
		var dev *Baseinfo.Device
		_ = colDevice.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "find no device!"
			_ = logger.Log("UnboundDeviceErr:", "find no device!")
			return response
		}
		if tokenuser != dev.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't operate another user' device!"
			_ = logger.Log("UnboundDeviceErr:", "cant't operate another user' device!")
			return response
		}
		var node *Baseinfo.Device
		_ = colDevice.FindOne(context.Background(), bson.D{{"gatewayid", dev.Deviceid}, {"isnode", true}}).Decode(&node)
		if node != nil {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't unbound gateway with nodes under it!"
			_ = logger.Log("UnboundDeviceErr:", "cant't unbound gateway with nodes under it!")
			return response
		}
		_ = Baseinfo.Client.Database("test").Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
			err := sessionContext.StartTransaction()
			if err != nil {
				return err
			}
			errcode, errmsg, data = Baseinfo.UnboundDeviceByid(dev, kind, sessionContext, colDevice, colSpace)
			if errmsg != "" {
				_ = sessionContext.AbortTransaction(sessionContext)
			}
			return nil
		})
	} else {
		var dev *Baseinfo.Device
		_ = colDevice.FindOne(context.Background(), bson.D{{"deviceid", deviceid}}).Decode(&dev)
		if dev == nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "find no device!"
			_ = logger.Log("UnboundDeviceErr:", "find no device!")
			return response
		}
		if tokenuser != dev.Userid {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't operate another user' device!"
			_ = logger.Log("unbound device err:", "cant't operate another user' device!")
			return response
		}
		var node *Baseinfo.Device
		_ = colDevice.FindOne(context.Background(), bson.D{{"gatewayid", deviceid}, {"isnode", true}}).Decode(&node)
		if node != nil {
			response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
			response.Msg = "cant't unbound gateway with nodes under it!"
			_ = logger.Log("UnboundDeviceErr:", "cant't unbound gateway with nodes under it!")
			return response
		}
		_ = Baseinfo.Client.Database("test").Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
			err := sessionContext.StartTransaction()
			if err != nil {
				return err
			}
			errcode, errmsg, data = Baseinfo.UnboundDeviceBydeviceid(dev, kind, sessionContext, colDevice, colSpace)
			if errmsg != "" {
				_ = sessionContext.AbortTransaction(sessionContext)
			}
			return nil
		})
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
	colDevice := Baseinfo.Client.Database("test").Collection("device")
	colHistory := Baseinfo.Client.Database("test").Collection("history")

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

	sessionErr := Baseinfo.Client.Database("test").Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		//记录至history表
		_, errIns := colHistory.InsertOne(sessionContext, history)
		if errIns != nil {
			response.Code = Baseinfo.CONST_INSERT_FAIL
			response.Msg = "fail to insert data to history "
			_ = logger.Log("UploadDataErr:", errIns)
			return errors.New("fail to insert data to history ")
		}

		//更新device表的expand
		id, errObj := primitive.ObjectIDFromHex(deviceid)
		if errObj == nil {
			filter := bson.D{{"_id", id}}
			update := bson.D{{"$set", bson.D{{"expand", &struct {
				T    string
				Data interface{}
			}{
				T:    r.T,
				Data: r.Data,
			},
			}}}}
			updateResult := colDevice.FindOneAndUpdate(sessionContext, filter, update)
			if updateResult.Err() != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = "fail to update device's history"
				_ = sessionContext.AbortTransaction(sessionContext)
				_ = logger.Log("UploadDataErr:", updateResult)
				return errors.New("fail to update device's history")
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
			updateResult := colDevice.FindOneAndUpdate(sessionContext, filter, update)
			if updateResult.Err() != nil {
				response.Code = Baseinfo.CONST_UPDATE_FAIL
				response.Msg = "fail to update device's history"
				_ = sessionContext.AbortTransaction(sessionContext)
				_ = logger.Log("Upload_Data_Err:", updateResult.Err().Error())
				return errors.New("fail to update device's history")
			}
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})
	if sessionErr != nil {
		_ = logger.Log("UploadDataErr:", sessionErr)
		return response
	}
	response.Code = Baseinfo.Success
	return response
}
