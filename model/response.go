package model

// 响应体相关
type LarkResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}
type Data struct {
}
