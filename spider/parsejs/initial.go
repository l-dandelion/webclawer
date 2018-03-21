package parsejs

import (
	"github.com/l-dandelion/webcrawler/module"
	"github.com/robertkrimen/otto"
)

func getReqsFromJS(jsCode string) ([]*module.Request, error) {
	initialReqs := []*module.Request{}
	dataList := NewDataList()
	vm := otto.New()
	vm.Set("aid", aid)
	vm.Set("dataList", dataList)
	_, err := vm.Run(jsCode)
	if err != nil {
		return initialReqs, err
	}
	for _, data := range dataList.Data {
		req, ok := data.(*module.Request)
		if ok {
			initialReqs = append(initialReqs, req)
		}
	}
	return initialReqs, nil
}
