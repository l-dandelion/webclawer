package parsejs

import (
	"testing"
)

func TestGetReqsFromJS(t *testing.T) {
	jsCode := `
        req1 = aid.NewRequest("www.baidu.com", "get")
        req2 = aid.NewRequest("www.codeforce.com", "get")
        dataList.Push(req1)
        dataList.Push(req2)
    `
	reqs, err := getReqsFromJS(jsCode)
	if err != nil {
		t.Fatalf("An error occurs when get requests from JS: %s", err)
	}
	if len(reqs) != 2 {
		t.Fatalf("Incorrect length of reqs. expected: %d, actual: %d.", 2, len(reqs))
	}
	if reqs[0].HTTPReq().Method != "get" {
		t.Fatalf("Incorrect method of req[0]. expected: %s, actual: %s", "get", reqs[0].HTTPReq().Method)
	}
	if reqs[0].HTTPReq().URL.String() != "www.baidu.com" {
		t.Fatalf("Incorrect URL of req[0]. expected: %s, actual: %s",
			"www.baidu.com", reqs[0].HTTPReq().URL.String())
	}

	if reqs[1].HTTPReq().Method != "get" {
		t.Fatalf("Incorrect method of req[1]. expected: %s, actual: %s", "get", reqs[1].HTTPReq().Method)
	}
	if reqs[1].HTTPReq().URL.String() != "www.codeforce.com" {
		t.Fatalf("Incorrect URL of req[1]. expected: %s, actual: %s",
			"www.codeforce.com", reqs[1].HTTPReq().URL.String())
	}
}
