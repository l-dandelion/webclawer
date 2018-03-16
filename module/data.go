package module

import (
	"net/http"
)

type Data interface {
	Valid() bool
}

type Request struct {
	httpReq *http.Request
	depth   uint32
}

func (req *Request) HTTPReq() *http.Request {
	return req.httpReq
}

func (req *Request) Depth() uint32 {
	return req.depth
}

func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

func NewRequest(httpReq *http.Request, depth uint32) *Request {
	return &Request{
		httpReq: httpReq,
		depth:   depth,
	}
}

type Response struct {
	httpResp *http.Response
	depth    uint32
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

func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{
		httpResp: httpResp,
		depth:    depth,
	}
}

type Item map[string]interface{}

func (item Item) Valid() bool {
	return item != nil
}
