package trace_error

import (
// "encoding/json"
// "github.com/google/uuid"
)

type Response struct {
	RequestId string `json:"requestId"`
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data"`
}

type ResponseOption func(*Response)

/*
func (res *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(res)
}
*/

func (self *Response) XRequestId(requestId string) *Response {
	self.RequestId = requestId
	return self
}

func RequestId(requestId string) ResponseOption {
	return func(self *Response) {
		self.RequestId = requestId
	}
}

func NewResponse(data any, opts ...ResponseOption) Response {
	self := Response{Code: "ok", Msg: "ok", Data: data}

	for _, opt := range opts {
		opt(&self)
	}

	if self.Data == nil {
		self.Data = map[string]any{}
	}

	/*
		if res.RequestId == "" {
			if id, e := uuid.NewUUID(); e == nil {
				res.RequestId = id.String()
			}
		}
	*/

	return self
}

func (self *Error) IntoResponse(opts ...ResponseOption) Response {
	res := Response{
		Code: self.CodeStr,
		Msg:  self.Msg,
	}

	for _, opt := range opts {
		opt(&res)
	}

	if res.Data == nil {
		res.Data = map[string]any{}
	}

	return res
}

func FromError(err *Error, opts ...ResponseOption) Response {
	res := Response{
		Code: err.CodeStr,
		Msg:  err.Msg,
	}

	for _, opt := range opts {
		opt(&res)
	}

	if res.Data == nil {
		res.Data = map[string]any{}
	}

	return res
}
