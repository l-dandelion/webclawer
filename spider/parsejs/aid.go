package parsejs

import (
	"errors"
	"gopcp.v2/helper/log"
	"net/http"
	"net/url"

	"github.com/l-dandelion/webcrawler/module"
)

var (
	logger = log.DLogger()
	aid    = &Aid{}
)

//辅助在js代码中生成相应类型的对象
type Aid struct{}

func (aid *Aid) NewRequest(url, method string, extras ...map[string]interface{}) *module.Request {
	httpReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		logger.Warnf("An Error occur when new a HTTP request: %s", err)
		return nil
	}
	return module.NewRequest(httpReq, extras...)
}

func (aid *Aid) NewItem(extras ...map[string]interface{}) module.Item {
	item := module.Item{}
	if len(extras) > 0 && extras[0] != nil {
		for key, val := range extras[0] {
			item[key] = val
		}
	}
	return item
}

func (aid *Aid) NewError(errMsg string) error {
	return errors.New(errMsg)
}

func (aid *Aid) ParseURL(href string) (*url.URL, error) {
	return url.Parse(href)
}
