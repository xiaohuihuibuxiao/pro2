package Baseinfo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type Space struct {
	Id        primitive.ObjectID   `json:"id" bson:"_id"`
	Mastered  string               //上级区域id
	Master    []string             //下级区域
	Devids    []primitive.ObjectID //objectid
	Level     int64                //用于区分该房源所属的位置的级别，具体到哪个位置
	Spacecode string
	Title     string
	Addr      string
	Userid    string
	External  interface{}
}

func AddDev(hid primitive.ObjectID, devids []primitive.ObjectID) error {
	return nil
}

func RemoveDev() {

}

func RemoveSpace(s *Space, col *mongo.Collection) (int64, interface{}) {
	masterspace := s.Master
	fmt.Println("自己的spaceid", s.Id, "读取到的master为", masterspace)
	_, err := col.DeleteOne(context.Background(), bson.M{"_id": s.Id})
	if err != nil {
		return CONST_DELETE_FAIL, err.Error()
	}
	if masterspace != nil || len(masterspace) > 0 {
		fmt.Println("masterspace不为空")
		for _, nextspaceid := range masterspace {
			fmt.Println("当前循环的spaceid", nextspaceid)
			var nextsapce *Space
			nextsid, _ := primitive.ObjectIDFromHex(nextspaceid)
			col.FindOne(context.Background(), bson.D{{"_id", nextsid}}).Decode(&nextsapce)
			if nextsapce != nil {
				fmt.Println("即将下一级循环", nextsapce.Id, nextsapce.Master)
				errcode, errmsg := RemoveSpace(nextsapce, col)
				if errmsg != nil {
					fmt.Println("出错了，将返回数据", errcode, errmsg)
					return errcode, errmsg
				}
			}
			//如果master里的id没找到相应的space,则忽略，就当已经被删
		}
	}
	return Success, nil
}

func GetFirstPartCode(mergename string, col *mongo.Collection) (int64, interface{}, string, *Dictionary) {
	fmt.Println("获取第一段时，输入的mergename ", mergename, col.Name())
	var dic *Dictionary
	err := col.FindOne(context.Background(), bson.D{{"mergername", mergename}}).Decode(&dic)
	if err != nil {
		return CONST_FIND_FAIL, err, "-1", nil
	}
	fmt.Println("获取第一段1111", dic)
	return Success, nil, dic.Code, dic
}

func GetSecondPartCode(upperdic *Dictionary, district, building, storey, room, place string, level int, col *mongo.Collection) (int64, interface{}, string) {
	errcode, errmsg, r, updistrictcode, name, upname := GetInfomation(upperdic, level, district, building, storey, room, place, col)
	fmt.Println("GetInfomation 函数返回的结果为 errcode", errcode, "errmsg", errmsg, "updistrictcode", updistrictcode, "name", name)
	fmt.Println("r", r)
	if errcode != Success {
		return errcode, errmsg, ""
	} else {
		if r != nil { //找到了对应的区域，直接返回code
			return Success, nil, r.Code
		}
		//	return 0,"test",""
		//在district表新建该区域，然后返回
		return NewDistrict(updistrictcode, upperdic, name, upname, level, col)
	}
}

func GetInfomation(dic *Dictionary, level int, district, building, storey, room, place string, col *mongo.Collection) (errcode int64, errmsg interface{}, r *District, updistrictcode string, name, upname string) {
	switch level {
	case 4:
		errcode, errmsg, r = FindDistrictByNameandLevel(dic.Code, district, 4, col)
		updistrictcode = "0000000000"
		name = district
		upname = ""
		return
	case 5:
		_, _, upr := FindDistrictBymergename(dic.Mergername+","+district, col)
		errcode, errmsg, r = FindDistrict(dic.Code, upr.Code, building, 5, col)
		if r != nil {
			updistrictcode = r.Parentcode
		} else { //TODO 默认上一级level区域肯定是存在的
			updistrictcode = upr.Code
		}
		name = building
		upname = district
		return
	case 6:
		_, _, upr := FindDistrictBymergename(dic.Mergername+","+district+","+building, col)
		fmt.Println("upr----", upr == nil, upr.Code, upr.Mergeaddr, upr.Dicaddr, upr.Parentcode)
		errcode, errmsg, r = FindDistrict(dic.Code, upr.Code, storey, 6, col)
		if r != nil {
			fmt.Println("111")
			updistrictcode = r.Parentcode
		} else { //TODO 默认上一级level区域肯定是存在的
			fmt.Println("2222")
			updistrictcode = upr.Code
		}
		name = storey
		upname = building
		return
	case 7:
		_, _, upr := FindDistrictBymergename(dic.Mergername+","+district+","+building+","+storey, col)
		errcode, errmsg, r = FindDistrict(dic.Code, upr.Code, room, 7, col)
		if r != nil {
			updistrictcode = r.Parentcode
		} else { //TODO 默认上一级level区域肯定是存在的
			updistrictcode = upr.Code
		}
		name = room
		upname = storey
		return
	case 8:
		_, _, upr := FindDistrictBymergename(dic.Mergername+","+district+","+building+","+storey+","+room, col) //TODO 默认上一级level区域肯定是存在的
		errcode, errmsg, r = FindDistrict(dic.Code, upr.Code, place, 8, col)
		if r != nil {
			updistrictcode = r.Parentcode
		} else {
			updistrictcode = upr.Code
		}
		name = place
		upname = room
		return
	default:
		return CONST_PARAM_ERROR, "pls check level !", nil, "", "", ""
	}
}

func Getaddr(s string, flag string) (str string) {

	strs := strings.Split(s, flag)
	for _, v := range strs {
		str = str + v
	}
	return str
}
