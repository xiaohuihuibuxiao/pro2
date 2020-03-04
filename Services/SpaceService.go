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

//--创建空间--
type WSpaceCreateService interface {
	NewSpace(r *SpaceCreateRequest) *CommonResponse
}
type SpaceCreateService struct{}

func (this SpaceCreateService) NewSpace(r *SpaceCreateRequest) *CommonResponse {
	response := &CommonResponse{}
	colSpace := Baseinfo.Client.Database("test").Collection("space")
	colDic := Baseinfo.Client.Database("test").Collection("dictionary")
	colDis := Baseinfo.Client.Database("test").Collection("district")

	errChecktoken, tokenuser := Baseinfo.LoginTokenAuth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("CreateSpaceErr", errChecktoken.Error())
		return response
	}
	//....生成space所需的全部信息
	level := r.Level
	province := r.Province
	city := r.City
	area := r.Area
	if province == "" || city == "" || area == "" || level == 0 {
		response.Code = Baseinfo.CONST_PARAM_LACK
		response.Msg = "provice ,city , area or level can't be ignored !"
		_ = logger.Log("CreateSpaceErr", "provice ,city , area or level can't be ignored !")
		return response
	}
	if level < 4 || level > 8 {
		response.Code = Baseinfo.CONST_PARAM_ERROR
		response.Msg = "province ,city , area or level can't be ignored !"
		_ = logger.Log("CreateSpaceErr", "provice ,city , area or level can't be ignored !")
		return response
	}
	district := r.District
	if district == "" {
		response.Code = Baseinfo.CONST_PARAM_LACK
		response.Msg = "disctrict can't be nil"
		_ = logger.Log("CreateSpaceErr", "disctrict can't be nil")
		return response
	}

	//获取省市区编码
	errcode, errmsg, firstpartcode, areadic := Baseinfo.GetFirstPartCode("中国,"+province+","+city+","+area, colDic)
	if errmsg != nil {
		response.Code = errcode
		response.Msg = errmsg.(error).Error()
		_ = logger.Log("CreateSpaceErr", errmsg.(error).Error())
		return response
	}

	building := r.Building
	storey := r.Storey
	room := r.Room
	place := r.Place

	errcode1, errmsg1, secondpartcode := Baseinfo.GetSecondPartCode(areadic, district, building, storey, room, place, level, colDis)
	if errmsg1 != nil {
		response.Code = errcode1
		response.Msg = errmsg1.(error).Error()
		_ = logger.Log("CreateSpaceErr", errmsg1.(error).Error())
		return response
	}

	spacecode := firstpartcode + secondpartcode
	var nowdis *Baseinfo.District
	_ = colDis.FindOne(context.Background(), bson.D{{"code", secondpartcode}, {"dictionarycode", areadic.Code}}).Decode(&nowdis)

	var checkspace *Baseinfo.Space
	_ = colSpace.FindOne(context.Background(), bson.D{{"spacecode", spacecode}}).Decode(&checkspace)
	if checkspace != nil {
		response.Code = Baseinfo.CONST_DATA_HASEXISTED
		response.Msg = "this space has existed!"
		_ = logger.Log("Create_Space_Err", "this space has existed!")
		return response
	}

	var upspaceid string
	//masteredspace--上级区域
	masteredSpace, _ := Baseinfo.FindMasteredSpace(spacecode, int64(r.Level), colSpace)
	if masteredSpace != nil {
		upspaceid = masteredSpace.Id.Hex()
	} else {
		upspaceid = "000000000000000000000000"
	}

	newspace := &Baseinfo.Space{
		Id:        primitive.NewObjectIDFromTimestamp(time.Now()),
		Mastered:  upspaceid,
		Master:    nil,
		Devids:    nil,
		Level:     int64(r.Level),
		Spacecode: spacecode,
		Title:     r.Title,
		Addr:      Baseinfo.Getaddr(nowdis.Mergeaddr, ","),
		Userid:    tokenuser,
		External:  nil,
	}
	insertResult, errIns := colSpace.InsertOne(context.Background(), newspace)
	if errIns != nil {
		response.Code = Baseinfo.CONST_INSERT_FAIL
		response.Msg = errIns.Error()
		_ = logger.Log("CreateSpaceErr", errIns.Error())
		return response
	}
	//在上级空间的master中添加新空间的id
	if masteredSpace != nil {
		if r.Level > 4 && r.Level < 9 {
			var m []string
			for _, v := range masteredSpace.Master {
				m = append(m, v)
			} //把原本的master先赋值给m
			m = append(m, insertResult.InsertedID.(primitive.ObjectID).Hex())
			colSpace.FindOneAndUpdate(context.Background(), bson.D{{"_id", masteredSpace.Id}}, bson.D{{"$set", bson.D{{"master", m}}}})
		}
	}
	var newsapceinfo *Baseinfo.Space
	_ = colSpace.FindOne(context.Background(), bson.D{{"_id", insertResult.InsertedID}}).Decode(&newsapceinfo)
	response.Code = Baseinfo.Success
	response.Data = newsapceinfo
	return response
}

//--查询空间--
type WSpaceQueryService interface {
	QuerySpace(r *SpaceQueryRequest) *CommonResponse
}
type SpaceQueryService struct{}

func (this SpaceQueryService) QuerySpace(r *SpaceQueryRequest) *CommonResponse {
	response := &CommonResponse{}
	colSpace := Baseinfo.Client.Database("test").Collection("space")

	errChecktoken, tokenuser := Baseinfo.LoginTokenAuth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("QuerySpaceErr:", errChecktoken.Error())
		return response
	}

	sid := r.Sid
	sid_obj, errObj := primitive.ObjectIDFromHex(sid)
	if errObj != nil {
		response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
		response.Msg = errObj.Error()
		_ = logger.Log("QuerySpaceErr:", errObj.Error())
		return response
	}
	var s *Baseinfo.Space
	err := colSpace.FindOne(context.Background(), bson.M{"_id": sid_obj}).Decode(&s)
	if s == nil {
		response.Code = Baseinfo.CONST_FIND_FAIL
		response.Msg = err.Error()
		_ = logger.Log("Query_Space_Err:", err.Error())
		return response
	}
	if tokenuser != s.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't query another user's space"
		_ = logger.Log("QuerySpaceErr:", "can't query another user's space")
		return response
	}
	var masterspace []*Baseinfo.Space
	if s.Master != nil || len(s.Master) > 0 {
		for _, v := range s.Master {
			var masterS *Baseinfo.Space
			mastersid, _ := primitive.ObjectIDFromHex(v)
			_ = colSpace.FindOne(context.Background(), bson.D{{"_id", mastersid}}).Decode(&masterS)
			if masterS != nil {
				masterspace = append(masterspace, masterS)
			}
		}
	}
	response.Code = Baseinfo.Success
	response.Data = s
	response.Expand = masterspace
	return response
}

//--修改空间--
type WSpaceReviseService interface {
	ReviseSapce(r *SpaceReviseRequest) *CommonResponse
}
type SpaceReviseService struct{}

func (this SpaceReviseService) ReviseSapce(r *SpaceReviseRequest) *CommonResponse {
	response := &CommonResponse{}
	colSpace := Baseinfo.Client.Database("test").Collection("space")

	errChecktoken, tokenuser := Baseinfo.LoginTokenAuth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("ReviseSpaceErr:", errChecktoken.Error())
		return response
	}

	sid := r.Sid
	sidObj, errObj := primitive.ObjectIDFromHex(sid)
	if errObj != nil {
		response.Code = Baseinfo.CONST_PARAM_ERROR
		response.Msg = errObj.Error()
		_ = logger.Log("ReviseSpaceErr:", errObj.Error())
		return response
	}
	filter := bson.D{{"_id", sidObj}}

	var s *Baseinfo.Space
	_ = colSpace.FindOne(context.Background(), filter).Decode(&s)
	if s == nil {
		response.Code = Baseinfo.CONST_DATA_UNEXISTED
		response.Msg = "find no space by sid "
		_ = logger.Log("ReviseSpaceErr:", "find no space by sid ")
		return response
	}
	if tokenuser != s.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't revise another user's space"
		_ = logger.Log("ReviseSpaceErr:", "can't revise another user's space")
		return response
	}

	var reviseddpace *Baseinfo.Space
	sessionErr := Baseinfo.Client.Database("test").Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		update := bson.D{{"$set", bson.D{{"title", r.Title}}}}
		_, errUpd := colSpace.UpdateOne(sessionContext, filter, update)
		if errUpd != nil {
			response.Code = Baseinfo.CONST_UPDATE_FAIL
			response.Msg = "fail to update space "
			_ = logger.Log("ReviseSpaceErr:", errUpd.Error())
			return errors.New("fail to update space ")
		}

		e := colSpace.FindOne(sessionContext, bson.D{{"_id", sidObj}}).Decode(&reviseddpace)
		if e != nil {
			response.Code = Baseinfo.CONST_FIND_FAIL
			response.Msg = "fail to find recently revised space "
			_ = sessionContext.AbortTransaction(sessionContext)
			_ = logger.Log("ReviseSpaceErr:", errUpd.Error())
			return errors.New("fail to find recently revised space ")
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})
	if sessionErr != nil {
		_ = logger.Log("ReviseSpaceErr:", sessionErr)
		return response
	}
	response.Code = Baseinfo.Success
	response.Data = reviseddpace
	return response

}

//--删除空间--
type WSpaceDelService interface {
	DelSapce(r *SpaceDelRequest) *CommonResponse
}
type SpaceDelService struct{}

func (this SpaceDelService) DelSapce(r *SpaceDelRequest) *CommonResponse {
	response := &CommonResponse{}
	colSpace := Baseinfo.Client.Database("test").Collection("space")

	errChecktoken, tokenuser := Baseinfo.LoginTokenAuth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("DeleteSpaceErr:", errChecktoken.Error())
		return response
	}
	sid, errObj := primitive.ObjectIDFromHex(r.Sid)
	if errObj != nil {
		response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
		response.Msg = errObj.Error()
		_ = logger.Log("DeleteSpaceErr:", errObj.Error())
		return response
	}

	var s *Baseinfo.Space
	errFind := colSpace.FindOne(context.Background(), bson.D{{"_id", sid}}).Decode(&s)
	if s == nil {
		response.Code = Baseinfo.CONST_PARAM_ERROR
		response.Msg = "no suc space ,pls check space id"
		_ = logger.Log("DeleteSpaceErr:", "no suc space ,pls check space id")
		return response
	}
	if errFind != nil {
		response.Code = Baseinfo.CONST_FIND_FAIL
		response.Msg = errFind.Error()
		_ = logger.Log("DeleteSpaceErr:", errFind.Error())
		return response
	}

	if tokenuser != s.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't revise another user's space"
		_ = logger.Log("DeleteSpaceErr:", "can't revise another user's space")
		return response
	}

	//区域绑定有设备时不允许删除区域
	if s.Devids != nil || len(s.Devids) > 0 {
		response.Code = Baseinfo.CONST_ACTION_UNALLOWED
		response.Msg = "devices are now hound in space,pls remove devices first !"
		_ = logger.Log("DeleteSpaceErr:", "devices are now hound in space,pls remove devices first !")
		return response
	}

	//该区域，以及下属的区域，全部删除！！！

	sessionErr := Baseinfo.Client.Database("test").Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		e, m := Baseinfo.RemoveSpace(s, sessionContext, colSpace)
		response.Code = e
		if m != nil {
			response.Msg = m.(error).Error()
			_ = sessionContext.AbortTransaction(sessionContext)
			_ = logger.Log("Delete_Space_Err:", m.(error).Error())
		}
		_ = sessionContext.CommitTransaction(sessionContext)
		return nil
	})
	if sessionErr != nil {
		_ = logger.Log("DeleteSpaceErr:", sessionErr)
		return response
	}
	return response
}

//--复制空间--
type WSpaceCloneService interface {
	CloneSpace(r *SpaceCloneRequest) *CommonResponse
}
type SpaceCloneService struct{}

func (this SpaceCloneService) CloneSpace(r *SpaceCloneRequest) *CommonResponse {
	response := &CommonResponse{}
	colSpace := Baseinfo.Client.Database("test").Collection("space")
	colDis := Baseinfo.Client.Database("test").Collection("district")
	colDic := Baseinfo.Client.Database("test").Collection("dictionary")

	errChecktoken, tokenuser := Baseinfo.LoginTokenAuth(r.Token)
	if errChecktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = errChecktoken.Error()
		_ = logger.Log("CloneSpaceErr:", errChecktoken.Error())
		return response
	}

	sid, errObj := primitive.ObjectIDFromHex(r.Sid)
	if errObj != nil {
		response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
		response.Msg = errObj.Error()
		_ = logger.Log("CloneSpaceErr:", errObj.Error())
		return response
	}

	var originlspace *Baseinfo.Space
	_ = colSpace.FindOne(context.Background(), bson.D{{"_id", sid}}).Decode(&originlspace)
	if originlspace == nil {
		response.Code = Baseinfo.CONST_DATA_UNEXISTED
		response.Msg = "can't find original space !"
		_ = logger.Log("CloneSpaceErr:", "can't find original space !")
		return response
	}

	if originlspace.Userid != tokenuser && originlspace.Userid != "" {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't clone other user's space !"
		_ = logger.Log("Clone_Space_Err:", "can't clone other user's space !")
		return response
	}

	//除了创建新的space，还需要创建新的district存储起来
	errNewdis, newDistrictcode := Baseinfo.CreateDistrict(originlspace, colDis, colDic)
	if errNewdis != nil {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't clone other user's space !"
		_ = logger.Log("CloneSpaceErr:", "can't clone other user's space !")
		return response
	}

	var newdistrict *Baseinfo.District
	code := newDistrictcode[6:]
	discode := originlspace.Spacecode[:6]
	_ = colDis.FindOne(context.Background(), bson.D{{"code", code}, {"dictionarycode", discode}}).Decode(&newdistrict)
	var newaddr string
	if newdistrict != nil {
		newaddr = Baseinfo.Getaddr(newdistrict.Mergeaddr, ",")
	}
	newspace := &Baseinfo.Space{
		Id:        primitive.NewObjectIDFromTimestamp(time.Now()),
		Mastered:  originlspace.Mastered,
		Master:    nil,
		Devids:    nil,
		Level:     originlspace.Level,
		Spacecode: newDistrictcode,
		Title:     originlspace.Title,
		Addr:      newaddr,
		Userid:    originlspace.Userid,
		External:  nil,
	}
	insertresult, errIns := colSpace.InsertOne(context.Background(), newspace)
	if errIns != nil {
		response.Code = Baseinfo.CONST_INSERT_FAIL
		response.Msg = "clone space failed:" + errIns.Error()
		_ = logger.Log("CloneSpaceErr:", "clone space failed:"+errIns.Error())
		return response
	}

	var newspaceinfo *Baseinfo.Space
	_ = colSpace.FindOne(context.Background(), bson.D{{"_id", insertresult.InsertedID.(primitive.ObjectID)}}).Decode(&newspaceinfo)
	response.Code = Baseinfo.Success
	response.Data = newspaceinfo
	return response
}
