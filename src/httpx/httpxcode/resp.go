package httpxcode

type RespBase struct {
	Code    CodeType `json:"code"`
	Message string   `json:"message"`
}

type RespStruct struct {
	RespBase
	Data interface{} `json:"data"`
}

func NewRespBase(code CodeType, message string) RespBase {
	return RespBase{Code: code, Message: message}
}
