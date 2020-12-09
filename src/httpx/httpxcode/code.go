package httpxcode

type (
	CodeType int64
)

const (
	CodeOK   CodeType = 0
	CodeBusy CodeType = 10

	CodeOkMessage   = "ok"
	CodeBusyMessage = "很抱歉，系统忙，请稍后重试"
)
