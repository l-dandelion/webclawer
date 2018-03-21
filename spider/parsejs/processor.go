package parsejs

import (
	"github.com/l-dandelion/webcrawler/module"
	"github.com/robertkrimen/otto"
)

func GetProcessorFromJS(jsCode string) module.ProcessItem {
	return func(item module.Item) (module.Item, error) {
		errorList := NewErrorList()
		vm := otto.New()
		vm.Set("aid", aid)
		vm.Set("item", item)
		vm.Set("errorList", errorList)
		_, err := vm.Run(jsCode)
		if err != nil {
			return item, err
		}
		if errorList != nil && len(errorList.Errs) > 0 {
			err = errorList.Errs[0]
		}
		return item, err
	}
}
