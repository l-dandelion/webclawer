package parsejs

import (
	"github.com/l-dandelion/webcrawler/module"
)

type DataList struct {
	Data []module.Data
}

func NewDataList() *DataList {
	return &DataList{
		Data: []module.Data{},
	}
}

func (dataList *DataList) Push(datum module.Data) {
	dataList.Data = append(dataList.Data, datum)
}

type ErrorList struct {
	Errs []error
}

func NewErrorList() *ErrorList {
	return &ErrorList{
		Errs: []error{},
	}
}

func (errorList *ErrorList) Push(err error) {
	errorList.Errs = append(errorList.Errs, err)
}
