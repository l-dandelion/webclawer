package parsejs

import (
	"testing"

	"github.com/l-dandelion/webcrawler/module"
	"github.com/robertkrimen/otto"
)

func TestDataList(t *testing.T) {
	vm := otto.New()
	vm.Set("aid", tAid)
	dataList := NewDataList()
	vm.Set("dataList", dataList)
	_, err := vm.Run(`
        req = aid.NewRequest("http://www.baidu.com", "get", {kind: "test"})
        item = aid.NewItem({file:"result/file1.xlsx"})
        dataList.Push(req)
        dataList.Push(item)
        `)
	if err != nil {
		t.Fatalf("vm run error: %s", err)
	}
	gDataList, err := vm.Get("dataList")
	if err != nil {
		t.Fatalf("get gDataList from vm fail: %s", err)
	}
	iDataList, err := gDataList.Export()
	if err != nil {
		t.Fatalf("export iDataList from vm fail: gDataList: %v, err: %s", gDataList, err)
	}
	dataList2, ok := iDataList.(*DataList)
	if ok {
		for _, data := range dataList2.Data {
			switch d := data.(type) {
			case *module.Request:
				if d.HTTPReq().URL.String() != "http://www.baidu.com" ||
					d.HTTPReq().Method != "get" ||
					d.Extra["kind"] != "test" {
					t.Fatalf("Incorrect valure.(req: %v)", d)
				}
			case module.Item:
				if d["file"].(string) != "result/file1.xlsx" {
					t.Fatalf("Incorrect valure.(item: %v)", d)
				}
			default:
				t.Fatalf("Incorrect type.(data: %T)", data)
			}
		}
	} else {
		t.Fatalf("convert error, %v", iDataList)
	}
}

func TestErrorList(t *testing.T) {
	vm := otto.New()
	vm.Set("aid", tAid)
	errorList := NewErrorList()
	vm.Set("errorList", errorList)
	_, err := vm.Run(`
        err1 = aid.NewError("test_error1")
        err2 = aid.NewError("test_error2")
        errorList.Push(err1)
        errorList.Push(err2)
        `)
	if err != nil {
		t.Fatalf("vm run error: %s", err)
	}
	errs := errorList.Errs
	if len(errs) != 2 {
		t.Fatalf("Incorrect error list: %v", errorList)
	}
	if errs[0].Error() != "test_error1" || errs[1].Error() != "test_error2" {
		t.Fatalf("Incorrect valure of errro list: %v", errorList)
	}
}
