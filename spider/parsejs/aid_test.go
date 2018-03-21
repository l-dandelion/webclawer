package parsejs

import (
	"testing"

	"github.com/l-dandelion/webcrawler/module"
	"github.com/robertkrimen/otto"
)

var tAid = &Aid{}

func TestFuncNewRequest(t *testing.T) {
	vm := otto.New()
	vm.Set("aid", tAid)
	_, err := vm.Run(`
        req = aid.NewRequest("http://www.baidu.com", "get", {kind: "test"})
        `)
	if err != nil {
		t.Fatalf("vm run error: %s", err)
	}
	gReq, err := vm.Get("req")
	if err != nil {
		t.Fatalf("get req from vm fail: %s", err)
	}
	iReq, err := gReq.Export()
	if err != nil {
		t.Fatalf("export req from vm fail: gReq: %v, err: %s", gReq, err)
	}
	req, ok := iReq.(*module.Request)
	if ok {
		if req.HTTPReq().Method != "get" ||
			req.HTTPReq().URL.String() != "http://www.baidu.com" ||
			req.Extra["kind"] != "test" {
			t.Fatalf("Incorrect valure.(req: %v)", req)
		}
	} else {
		t.Fatalf("convert error, %v", iReq)
	}
}

func TestFuncNewItem(t *testing.T) {
	vm := otto.New()
	vm.Set("aid", tAid)
	_, err := vm.Run(`
        item = aid.NewItem({file:"result/file1.xlsx"})
        `)
	if err != nil {
		t.Fatalf("vm run error: %s", err)
	}
	gItem, err := vm.Get("item")
	if err != nil {
		t.Fatalf("get item from vm fail: %s", err)
	}
	iItem, err := gItem.Export()
	if err != nil {
		t.Fatalf("export item from vm fail: gReq: %v, err: %s", gItem, err)
	}
	item, ok := iItem.(module.Item)
	if ok {
		file, ok := item["file"]
		if !ok {
			t.Fatalf("Incorrect valure.(item: %v)", item)
		}
		fileStr, ok := file.(string)
		if !ok {
			t.Fatalf("Incorrect valure.(item: %v)", item)
		}
		if fileStr != "result/file1.xlsx" {
			t.Fatalf("Incorrect valure.(item: %v)", item)
		}

	} else {
		t.Fatalf("convert error, %v", iItem)
	}
}

func TestFuncExtra(t *testing.T) {
	vm := otto.New()
	vm.Set("aid", tAid)
	_, err := vm.Run(`
        req = aid.NewRequest("http://www.baidu.com", "get", {kind: "test"});
        req.Extra["file"] = "haha.go";
        `)
	if err != nil {
		t.Fatalf("vm run error: %s", err)
	}
	gReq, err := vm.Get("req")
	if err != nil {
		t.Fatalf("get req from vm fail: %s", err)
	}
	iReq, err := gReq.Export()
	if err != nil {
		t.Fatalf("export req from vm fail: gReq: %v, err: %s", gReq, err)
	}
	req, ok := iReq.(*module.Request)
	if ok {
		if req.HTTPReq().Method != "get" ||
			req.HTTPReq().URL.String() != "http://www.baidu.com" ||
			req.Extra["kind"] != "test" {
			t.Fatalf("Incorrect valure.(req: %v)", req)
		}
	} else {
		t.Fatalf("convert error, %v", iReq)
	}
	if req.Extra["file"] != "haha.go" {
		t.Fatalf("use extra fail")
	}
}
