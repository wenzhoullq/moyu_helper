package lib

import (
	"github.com/gin-gonic/gin"
	"weixin_LLM/lib/constant"
)

type Response struct {
	ErrNo   int         `json:"err_no"`            // 错误码
	ErrMsg  string      `json:"err_msg"`           // 错误信息
	Results interface{} `json:"results,omitempty"` // 返回结果
}

func SetErrMsg(msg string) func(response *Response) {
	return func(r *Response) {
		r.ErrMsg = msg
	}
}

func SetErrNo(errNo int) func(*Response) {
	return func(r *Response) {
		r.ErrNo = errNo
	}
}

func SetResults(results interface{}) func(*Response) {
	return func(r *Response) {
		r.Results = results
	}
}

func NewResponse(ops ...func(response *Response)) *Response {
	//默认的resp
	resp := &Response{
		//ErrNo: Success,
	}
	for _, op := range ops {
		op(resp)
	}
	return resp
}

func SetResponse(resp *Response, ops ...func(response *Response)) {
	for _, op := range ops {
		op(resp)
	}
}

func SetContextResponse(c *gin.Context, resp *Response) {
	c.JSON(200, resp)
	return
}

func SetContextErrorResponse(c *gin.Context, err error) {
	resp := NewResponse(SetErrNo(constant.ParamErr), SetErrMsg(err.Error()))
	c.JSON(200, resp)
	return
}
