package Services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	col_space := Baseinfo.Client.Database("test").Collection("space")
	col_dic := Baseinfo.Client.Database("test").Collection("dictionary")
	col_dis := Baseinfo.Client.Database("test").Collection("district")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
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
		return response
	}
	if level < 4 || level > 8 {
		response.Code = Baseinfo.CONST_PARAM_ERROR
		response.Msg = "provice ,city , area or level can't be ignored !"
		return response
	}

	errcode, errmsg, firstpartcode, areadic := Baseinfo.GetFirstPartCode("中国,"+province+","+city+","+area, col_dic)
	//fmt.Println("获取第一段code", errcode, errmsg, firstpartcode)
	if errmsg != nil {
		response.Code = errcode
		response.Msg = errmsg.(error).Error()
		return response
	}

	district := r.District
	if district == "" {
		response.Code = Baseinfo.CONST_PARAM_LACK
		response.Msg = "disctrict can't be nil"
		return response
	}
	building := r.Building
	storey := r.Storey
	room := r.Room
	place := r.Place

	errcode1, errmsg1, secondpartcode := Baseinfo.GetSecondPartCode(areadic, district, building, storey, room, place, level, col_dis)
	//	fmt.Println("获取第二段地址", errcode, errmsg, secondpartcode)
	if errmsg1 != nil {
		response.Code = errcode1
		response.Msg = errmsg1.(error).Error()
		return response
	}

	spacecode := firstpartcode + secondpartcode
	var nowdis *Baseinfo.District
	col_dis.FindOne(context.Background(), bson.D{{"code", secondpartcode}, {"dictionarycode", areadic.Code}}).Decode(&nowdis)

	var checkspace *Baseinfo.Space
	col_space.FindOne(context.Background(), bson.D{{"spacecode", spacecode}}).Decode(&checkspace)
	if checkspace != nil {
		response.Code = Baseinfo.CONST_DATA_HASEXISTED
		response.Msg = "this space has existed!"
		return response
	}

	var upspaceid string
	masteredSpace, err_upspace := Baseinfo.FindMasteredSpace(spacecode, int64(r.Level), col_space)
	if err_upspace != nil {
		if masteredSpace == nil {
			upspaceid = "000000000000000000000000"
		} else {
			upspaceid = masteredSpace.Id.Hex()

		}
	}

	newspace := &Baseinfo.Space{
		Id:        primitive.NewObjectIDFromTimestamp(time.Now()),
		Mastered:  upspaceid, //这个直接填上
		Master:    nil,
		Devids:    nil,
		Level:     int64(r.Level),
		Spacecode: spacecode,
		Title:     r.Title,
		Addr:      Baseinfo.Getaddr(nowdis.Mergeaddr, ","),
		Userid:    tokenuser,
		External:  nil,
	}
	insert_result, err_ins := col_space.InsertOne(context.Background(), newspace)
	if err_ins != nil {
		response.Code = Baseinfo.CONST_INSERT_FAIL
		response.Msg = err_ins.Error()
		return response
	}
	//在上级空间的master中添加新空间的id
	if r.Level > 4 && r.Level < 9 {
		var m []string
		m = masteredSpace.Master
		m = append(m, insert_result.InsertedID.(string))
		col_space.FindOneAndUpdate(context.Background(), bson.D{{"_id", masteredSpace.Id}}, bson.D{{"$set", bson.D{{"master", m}}}})
	}
	response.Code = Baseinfo.Success
	response.Data = insert_result.InsertedID
	return response
}

//--查询空间--
type WSpaceQueryService interface {
	QuerySpace(r *SpaceQueryRequest) *CommonResponse
}
type SpaceQueryService struct{}

func (this SpaceQueryService) QuerySpace(r *SpaceQueryRequest) *CommonResponse {
	response := &CommonResponse{}
	col_space := Baseinfo.Client.Database("test").Collection("space")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}

	sid := r.Sid
	sid_obj, err_obj := primitive.ObjectIDFromHex(sid)
	if err_obj != nil {
		response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
		response.Msg = err_obj.Error()
		return response
	}
	var s *Baseinfo.Space
	err := col_space.FindOne(context.Background(), bson.M{"_id": sid_obj}).Decode(&s)
	if s == nil {
		response.Code = Baseinfo.CONST_FIND_FAIL
		response.Msg = err.Error()
		return response
	}
	if tokenuser != s.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't query another user's space"
		return response
	}
	response.Code = Baseinfo.Success
	response.Data = s
	return response
}

//--修改空间--
type WSpaceReviseService interface {
	ReviseSapce(r *SpaceReviseRequest) *CommonResponse
}
type SpaceReviseService struct{}

func (this SpaceReviseService) ReviseSapce(r *SpaceReviseRequest) *CommonResponse {
	response := &CommonResponse{}
	col_space := Baseinfo.Client.Database("test").Collection("space")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}

	sid := r.Sid
	sid_obj, err_obj := primitive.ObjectIDFromHex(sid)
	if err_obj != nil {
		response.Code = Baseinfo.CONST_PARAM_ERROR
		response.Msg = err_obj.Error()
		return response
	}
	filter := bson.D{{"_id", sid_obj}}

	var s *Baseinfo.Space
	col_space.FindOne(context.Background(), filter).Decode(&s)
	if s == nil {
		response.Code = Baseinfo.CONST_DATA_UNEXISTED
		response.Msg = "find no space by sid "
		return response
	}
	if tokenuser != s.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't revise another user's space"
		return response
	}

	update := bson.D{{"$set", bson.D{{"title", r.Title}}}}
	_, err_upd := col_space.UpdateOne(context.Background(), filter, update)
	if err_upd != nil {
		response.Code = Baseinfo.CONST_UPDATE_FAIL
		response.Msg = err_upd.Error()
		return response
	}
	response.Code = Baseinfo.Success
	return response

}

//--删除空间--
type WSpaceDelService interface {
	DelSapce(r *SpaceDelRequest) *CommonResponse
}
type SpaceDelService struct{}

func (this SpaceDelService) DelSapce(r *SpaceDelRequest) *CommonResponse {
	response := &CommonResponse{}
	col_space := Baseinfo.Client.Database("test").Collection("space")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}
	sid, err_obj := primitive.ObjectIDFromHex(r.Sid)
	if err_obj != nil {
		response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
		response.Msg = err_obj.Error()
		return response
	}

	var s *Baseinfo.Space
	err_find := col_space.FindOne(context.Background(), bson.D{{"_id", sid}}).Decode(&s)
	if s == nil {
		response.Code = Baseinfo.CONST_PARAM_ERROR
		response.Msg = "no suc space ,pls check space id"
		return response
	}
	if err_find != nil {
		response.Code = Baseinfo.CONST_FIND_FAIL
		response.Msg = err_find.Error()
		return response
	}

	if tokenuser != s.Userid {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't revise another user's space"
		return response
	}

	//区域绑定有设备时不允许删除区域
	if s.Devids != nil || len(s.Devids) > 0 {
		response.Code = Baseinfo.CONST_ACTION_UNALLOWED
		response.Msg = "devices are now hound in space,pls remove devices first !"
		return response
	}

	//该区域，以及下属的区域，全部删除！！！
	e, m := Baseinfo.RemoveSpace(s, col_space)
	response.Code = e
	response.Msg = m.(error).Error()
	return response
}

//--复制空间--
type WSpaceCloneService interface {
	CloneSpace(r *SpaceCloneRequest) *CommonResponse
}
type SpaceCloneService struct{}

func (this SpaceCloneService) CloneSpace(r *SpaceCloneRequest) *CommonResponse {
	response := &CommonResponse{}
	col_space := Baseinfo.Client.Database("test").Collection("space")
	col_dis := Baseinfo.Client.Database("test").Collection("district")
	col_dic := Baseinfo.Client.Database("test").Collection("dictionary")

	err_checktoken, tokenuser := Baseinfo.Logintokenauth(r.Token)
	if err_checktoken != nil {
		response.Code = Baseinfo.CONST_TOEKN_INVALID
		response.Msg = err_checktoken.Error()
		return response
	}

	sid, err_obj := primitive.ObjectIDFromHex(r.Sid)
	if err_obj != nil {
		response.Code = Baseinfo.CONST_UNMARSHALL_FAIL
		response.Msg = err_obj.Error()
		return response
	}

	var originlspace *Baseinfo.Space
	col_space.FindOne(context.Background(), bson.D{{"_id", sid}}).Decode(&originlspace)
	if originlspace == nil {
		response.Code = Baseinfo.CONST_DATA_UNEXISTED
		response.Msg = "can't find original space !"
		return response
	}

	if originlspace.Userid != tokenuser && originlspace.Userid != "" {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't clone other user's space !"
		return response
	}

	//除了创建新的space，还需要创建新的district存储起来 TODO
	err_newdis := Baseinfo.CreateDistrict(originlspace, col_dis, col_dic)
	if err_newdis != nil {
		response.Code = Baseinfo.CONST_UNAUTHORUTY_USER
		response.Msg = "can't clone other user's space !"
		return response
	}
	newspace := &Baseinfo.Space{
		Id:        primitive.NewObjectIDFromTimestamp(time.Now()),
		Mastered:  originlspace.Mastered,
		Master:    nil,
		Devids:    nil,
		Level:     originlspace.Level,
		Spacecode: "",
		Title:     "",
		Addr:      "",
		Userid:    originlspace.Userid,
		External:  nil,
	}
	insertresult, err_ins := col_space.InsertOne(context.Background(), newspace)
	if err_ins != nil {
		response.Code = Baseinfo.CONST_INSERT_FAIL
		response.Msg = "clone space failed:" + err_ins.Error()
		return response
	}

	response.Code = Baseinfo.Success
	response.Data = insertresult.InsertedID
	return response
}
