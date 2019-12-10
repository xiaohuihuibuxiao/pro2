package Baseinfo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"strconv"
	"time"
)

type District struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	//Areacode string //所属省市区的在dictionary表中的code(限定园区大楼的空间范围)
	Code           string //当前区域编码--即dictionary表中的code
	Name           string //如：漕河泾园区，5号楼，4楼 402室 xx工位
	Parentcode     string //上层区域编码 //园区则填写所在省市区的code
	Mergeaddr      string //园区到当前这一级的地址连在一起的地址
	Level          int    //1-xx园区/小区  2-xx楼 3-xx层 4-xx室/公共区域（如会议室） 5-xx工位
	Dictionarycode string //省市区的code，即dictionary
	Dicaddr        string //所属省市区的地址
	Title          string
}

//创建新district（园区 楼 层室 工位）
func NewDistrict(updistrictcode string, dic *Dictionary, name, upname string, level int, col *mongo.Collection) (int64, interface{}, string) {
	//name 新建区域的名字 upditrictcode 上级区域的code(不同level的c上级code不一样)
	opt := options.Find().SetSort(bson.D{{"code", -1}}) //按照code倒序，查询目前该level，该diccode下的区域情况
	curror, err_d := col.Find(context.Background(), bson.D{
		{"level", level},
		{"parentcode", updistrictcode},
		{"dictionarycode", dic.Code},
	}, opt)
	if err_d != nil {
		return CONST_FIND_FAIL, err_d.Error(), ""
	}
	if curror.Err() != nil {
		return CONST_FIND_FAIL, curror.Err().Error(), ""
	}
	var allDis []*District
	err_all := curror.All(context.Background(), &allDis)
	if err_all != nil {
		return CONST_FIND_FAIL, err_all.Error(), ""
	}
	curror.Close(context.Background())
	//此时已拿到全部数据
	largestcode := Getnewcode(updistrictcode, level, allDis)

	newmergeaddr := Getnewmergeaddr(level, dic, updistrictcode, upname, col)
	//填补新district的信息
	newdistrict := &District{
		Id:             primitive.NewObjectIDFromTimestamp(time.Now()),
		Code:           largestcode,
		Name:           name,
		Parentcode:     updistrictcode,
		Mergeaddr:      newmergeaddr + "," + name,
		Level:          level,
		Dictionarycode: dic.Code,
		Dicaddr:        dic.Mergername,
	}
	insert_r, err := col.InsertOne(context.Background(), newdistrict)
	if err != nil {
		return CONST_INSERT_FAIL, err.Error(), ""
	}
	var distri *District
	filter0 := bson.D{
		{"_id", insert_r.InsertedID.(primitive.ObjectID)},
	}
	col.FindOne(context.Background(), filter0).Decode(&distri)
	return Success, nil, distri.Code
}

func CreateDistrict(s *Space, col_dis, col_dic *mongo.Collection) (error, string) {
	//col--是district表
	var district0 *District
	col_dis.FindOne(context.Background(), bson.D{{"code", s.Spacecode[6:16]}, {"dictionarycode", s.Spacecode[:6]}}).Decode(&district0)
	if district0 == nil {
		return errors.New("no coresponding district for the space"), ""
	}
	var Up_district *District
	col_dis.FindOne(context.Background(), bson.D{{"code", district0.Parentcode}, {"dictionarycode", s.Spacecode[:6]}}).Decode(&Up_district)
	if Up_district == nil {
		return errors.New("no coresponding upper district for the space"), ""
	}
	var dictionary0 *Dictionary
	col_dic.FindOne(context.Background(), bson.D{{"code", district0.Dictionarycode}}).Decode(&dictionary0)
	if dictionary0 == nil {
		return errors.New("no coresponding dictionary for the space"), ""
	}
	_, error_insert, newdistrictcode := NewDistrict(district0.Parentcode, dictionary0, district0.Name, Up_district.Name, int(s.Level), col_dis)
	if error_insert != nil {
		return error_insert.(error), ""
	}
	return nil, district0.Dictionarycode + newdistrictcode
}

//---------------------------------------------------------------
func Getnewcode(updistrictcode string, level int, allDis []*District) (largestcode string) {
	//updic--所属省市区的数据，在level位4且需要新建区域时会用到
	if len(allDis) > 0 { //在该level下已经存在空间，code+1
		num, _ := strconv.Atoi(allDis[0].Code)                                   //TODO 万一是99了怎么处理？？？？------------------------
		largestcode = strconv.Itoa(num + 1*int(math.Pow(100, float64(8-level)))) //+1 +100 +10000...
		if len([]rune(largestcode)) == 9 {                                       //只可能出现9位或10位，10位不需要动
			largestcode = "0" + largestcode
		}
		return largestcode
	} else {
		//没找到，code就是它上级的code（xxx0000）开始加1 ，xxx0100
		num, _ := strconv.Atoi(updistrictcode)
		num = num + 1*int(math.Pow(100, float64(8-level)))
		largestcode := strconv.Itoa(num)
		if len([]rune(largestcode)) == 9 { //只可能出现9位或10位，10位不需要动。 出现其他位数的话肯定是有问题！！！！理论上不会出现，这里不做考虑
			largestcode = "0" + largestcode
		}
		return largestcode
	}
}

func Getnewmergeaddr(level int, dic *Dictionary, updistrictcode, upname string, col *mongo.Collection) string {
	if level == 4 {
		return dic.Mergername
	} else {
		var r *District
		filter := bson.D{
			{"name", upname},
			{"code", updistrictcode},
			{"level", (level - 1)},
			{"dictionarycode", dic.Code},
		}
		err := col.FindOne(context.Background(), filter).Decode(&r)
		if err != nil {
			return "find_error" + err.Error()
		}
		return r.Mergeaddr
	}
}

//根据name和level查找
func FindDistrictByNameandLevel(diccode string, name string, level int, col *mongo.Collection) (int64, interface{}, *District) {

	filter := bson.D{{"name", name}, {"level", level}, {"dictionarycode", diccode}}
	var result *District
	err := col.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return CONST_FIND_FAIL, "find no document in district", nil
		} else {
			return Success, "", nil
		}
	} else {
		return Success, "", result
	}

}

//根据name和level查找
func FindDistrict(dictionarycode, updiscode, name string, level int, col *mongo.Collection) (int64, interface{}, *District) {

	filter := bson.D{{"name", name}, {"level", level}, {"parentcode", updiscode}, {"dictionarycode", dictionarycode}}
	var result *District
	err := col.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return CONST_FIND_FAIL, "find no document in district", nil
		} else {
			return Success, "", nil
		}
	} else {
		return Success, "", result
	}
}

func FindDistrictBymergename(mergename string, col *mongo.Collection) (int64, interface{}, *District) {

	filter := bson.D{{"mergeaddr", mergename}}
	var result *District
	err := col.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return CONST_FIND_FAIL, "find no document in district", nil
		} else {
			return Success, "", nil
		}
	} else {
		return Success, "", result
	}

}
