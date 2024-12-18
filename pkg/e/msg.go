package e

var MsgFlags = map[int]string{
	SUCCESS:           "ok",
	ERROR:             "fail",
	InvalidParams:     "请求参数错误",
	ErrorTemplateLoad: "模版加载错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
