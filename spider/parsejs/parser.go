package parsejs

import (
	"github.com/l-dandelion/webcrawler/module"
	"github.com/robertkrimen/otto"
)

func GetParserFromJS(jsCode string) module.ParseResponse {
	return func(resp *module.Response) ([]module.Data, []error) {
		dataList := NewDataList()
		errorList := NewErrorList()
		vm := otto.New()
		vm.Set("aid", aid)
		vm.Set("dataList", dataList)
		vm.Set("errorList", errorList)
		vm.Set("resp", resp)
		_, err := vm.Run(jsCode)
		if err != nil {
			return dataList.Data, []error{err}
		}
		return dataList.Data, errorList.Errs
	}
}
