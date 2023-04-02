package models

type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	Id      string      `json:"id"`
}

type RpcMessage struct {
	Method string `json:"method"`
	Params any    `json:"params,omitempty"`
}

type Notification struct {
	RpcMessage
	Jsonrpc string `json:"jsonrpc"`
}

func NewJsonRpcNotification(method string) Notification {
	return Notification{
		RpcMessage: RpcMessage{
			Method: method,
		},
		Jsonrpc: "2.0",
	}
}

type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *ErrorResp  `json:"error,omitempty"`
}

func NewJsonRpcResponseError(id string, code int, message string) Response {
	return Response{
		Jsonrpc: "2.0",
		Id:      id,
		Result:  nil,
		Error: &ErrorResp{
			Code:    code,
			Message: message,
		},
	}
}

func NewJsonRpcResponseOk(id string) Response {
	return Response{
		Jsonrpc: "2.0",
		Id:      id,
		Result: &ErrorResp{
			Code:    0,
			Message: "OK",
		},
		Error: nil,
	}
}
