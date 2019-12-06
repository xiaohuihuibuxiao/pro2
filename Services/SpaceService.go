package Services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pro2/Baseinfo"
	"time"
)

//--创建康健--
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
	//fmt.Println("----------/////", checkspace)
	if checkspace != nil {
		response.Code = Baseinfo.CONST_DATA_HASEXISTED
		response.Msg = "this space has existed!"
		return response
	}

	newspace := &Baseinfo.Space{
		Id:        primitive.NewObjectIDFromTimestamp(time.Now()),
		Mastered:  "",
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

	response.Code = Baseinfo.Success
	response.Data = insert_result.InsertedID
	return response
}
