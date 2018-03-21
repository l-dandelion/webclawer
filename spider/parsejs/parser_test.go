package parsejs

import (
	"fmt"
	"testing"

	"github.com/l-dandelion/webcrawler/module"
)

func TestGetPaeserFromJS(t *testing.T) {
	jsCode := `
        req = aid.NewRequest("www.baidu.com", "get");
        item = aid.NewItem({URL:"www.baidu.com"});
        err = aid.NewError("error_test")
        dataList.Push(req)
        dataList.Push(item)
        dataList.Push(err)
    `
	parser := getParserFromJs(jsCode)

	dataList, errs := parser(nil, 0)
	if len(dataList) != 2 {
		t.Fatalf("Incorrect length of dataList, expected: %d, actual: %d", 2, len(dataList))
	}

	req, ok := dataList[0].(*module.Request)
	if !ok {
		t.Fatalf("Incorrect type of dataList[0], type: %T", dataList[0])
	}
	if req.HTTPReq().Method != "get" {
		t.Fatalf("Incorrect method of request, expected: %s, actual: %s", "get", req.HTTPReq().Method)
	}
	if req.HTTPReq().URL.String() != "www.baidu.com" {
		t.Fatalf("Incorrect url of request, expected: %s, actual: %s", "www.baidu.com", req.HTTPReq().URL)
	}
	item, ok := dataList[1].(module.Item)
	if _, ok := item["URL"]; !ok {
		t.Fatalf(`Incorrect item, could't find item["URL"]`)
	}
	if item["URL"] != "www.baidu.com" {
		t.Fatalf("Incorrect URL of item, expected: %s, actual: %s", "www.baidu.com", item["URL"])
	}

	if len(errs) != 1 {
		t.Fatalf("Incorrect length of errorList, expected: %d, actual: %d", 1, len(errs))
	}
}
