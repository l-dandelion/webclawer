package module

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/PuerkitoBio/goquery"
	"github.com/l-dandelion/webcrawler/toolkit/reader"
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
	text     []byte            // 下载内容Body的字节流格式
	dom      *goquery.Document //下载内容Body为html时，可转换为Dom的对象
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

func (resp *Response) GetText() ([]byte, error) {
	if resp.text != nil {
		return resp.text, nil
	}
	multiReader, err := reader.NewMultipleReader(resp.httpResp.Body)
	resp.httpResp.Body = multiReader.Reader()
	defer func() {
		resp.httpResp.Body.Close()
		resp.httpResp.Body = multiReader.Reader()
	}()
	var contentType, pageEncode string
	// 优先从响应头读取编码类型
	contentType = resp.httpResp.Header.Get("Content-Type")
	if _, params, err := mime.ParseMediaType(contentType); err == nil {
		if cs, ok := params["charset"]; ok {
			pageEncode = strings.ToLower(strings.TrimSpace(cs))
		}
	}
	// 响应头未指定编码类型时，从请求头读取
	if len(pageEncode) == 0 {
		contentType = resp.httpResp.Request.Header.Get("Content-Type")
		if _, params, err := mime.ParseMediaType(contentType); err == nil {
			if cs, ok := params["charset"]; ok {
				pageEncode = strings.ToLower(strings.TrimSpace(cs))
			}
		}
	}

	switch pageEncode {
	// 不做转码处理
	case "utf8", "utf-8", "unicode-1-1-utf-8":
	default:
		// 指定了编码类型，但不是utf8时，自动转码为utf8
		// get converter to utf-8
		// Charset auto determine. Use golang.org/x/net/html/charset. Get response body and change it to utf-8
		var destReader io.Reader

		if len(pageEncode) == 0 {
			destReader, err = charset.NewReader(resp.httpResp.Body, "")
		} else {
			destReader, err = charset.NewReaderLabel(pageEncode, resp.httpResp.Body)
		}

		if err == nil {
			resp.text, err = ioutil.ReadAll(destReader)
			if err == nil {
				return resp.text, err
			}
		}

	}
	// 不做转码处理
	resp.text, err = ioutil.ReadAll(resp.httpResp.Body)

	return resp.text, err
}

func (resp *Response) GetDom() (*goquery.Document, error) {
	if resp.dom != nil {
		return resp.dom, nil
	}
	text, err := resp.GetText()
	if err != nil {
		return nil, err
	}
	resp.dom, err = goquery.NewDocumentFromReader(bytes.NewReader(text))
	return resp.dom, err
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
