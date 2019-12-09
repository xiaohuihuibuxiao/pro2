package Baseinfo

const (
	Success             = 200
	Fail                = 500
	CONST_TOEKN_INVALID = 1001 //token过期或者无效

	CONST_INSERT_FAIL      = 1002 //数据库插入失败
	CONST_UNMARSHALL_FAIL  = 1003 //数据解析错误
	CONST_UPDATE_FAIL      = 1004 //数据库更新失败
	CONST_DATABASE_ERROR   = 1005 //数据库错误
	CONST_PARAM_LACK       = 1006 //缺少参数
	CONST_UNAUTHORUTY_USER = 1007 //缺少权限
	CONST_FIND_FAIL        = 1008 //数据查询失败
	CONST_PARAM_ERROR      = 1009 //参数错误
	CONST_DELETE_FAIL      = 1010 //数据删除失败

	CONST_USER_NOTEXIST   = 1011 //用户不存在
	CONST_USERPWD_UNMATCH = 1012 //账号密码不匹配
	CONST_USER_OCCUPIED   = 1014 //该账号已存在

	CONST_TOEKN_ERROR      = 1013 //token过期或者无效
	CONST_DATA_HASEXISTED  = 1015 //数据已存在
	CONST_ACTION_UNALLOWED = 1016 //操作不允许
	CONST_DATA_UNEXISTED   = 1017 //数据不存在
)

type Response struct {
	Code   int64
	Msg    interface{}
	Data   interface{}
	Expand interface{}
}
