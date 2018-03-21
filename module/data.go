package module

import (
	"net/http"
)

type Data interface {
	Valid() bool
}

//Extra 记录使用者提供的额外信息
//获取完请求后会将额外信息放入响应中
type Request struct {
	httpReq *http.Request
	depth   uint32
	Extra   map[string]interface{}
}

func (req *Request) HTTPReq() *http.Request {
	return req.httpReq
}

func (req *Request) Depth() uint32 {
	return req.depth
}

func (req *Request) SetExtra(key string, val interface{}) {
	req.Extra[key] = val
}

func (req *Request) SetDepth(depth uint32) {
	req.depth = depth
}

func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

func NewRequest(httpReq *http.Request, extras ...map[string]interface{}) *Request {
	var extra map[string]interface{}
	if len(extras) != 0 {
		extra = extras[0]
	} else {
		extra = map[string]interface{}{}
	}
	return &Request{
		httpReq: httpReq,
		Extra:   extra,
	}
}

//Extra 值为对应的Request的值
type Response struct {
	httpResp *http.Response
	depth    uint32
	Extra    map[string]interface{}
}

func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

func (resp *Response) Depth() uint32 {
	return resp.depth
}

func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

func NewResponse(httpResp *http.Response, depth uint32, extras ...map[string]interface{}) *Response {
	var extra map[string]interface{}
	if len(extras) != 0 {
		extra = extras[0]
	} else {
		extra = map[string]interface{}{}
	}
	return &Response{
		httpResp: httpResp,
		depth:    depth,
		Extra:    extra,
	}
}

type Item map[string]interface{}

func (item Item) Valid() bool {
	return item != nil
}
