package errors

import (
	"bytes"
	"fmt"
	"strings"
)

type ErrorType string

const (
	// ERROR_TYPE_DOWNLOADER 代表下载器错误。
	ERROR_TYPE_DOWNLOADER ErrorType = "downloader error"
	// ERROR_TYPE_ANALYZER 代表分析器错误。
	ERROR_TYPE_ANALYZER ErrorType = "analyzer error"
	// ERROR_TYPE_PIPELINE 代表条目处理管道错误。
	ERROR_TYPE_PIPELINE ErrorType = "pipeline error"
	// ERROR_TYPE_SCHEDULER 代表调度器错误。
	ERROR_TYPE_SCHEDULER ErrorType = "scheduler error"
)

type CrawlerError interface {
	Type() ErrorType
	Error() string
}

type myCrawlerError struct {
	errType    ErrorType
	errMsg     string
	fullErrMsg string
}

func (ce *myCrawlerError) Type() ErrorType {
	return ce.errType
}

func (ce *myCrawlerError) Error() string {
	if ce.fullErrMsg == "" {
		ce.genFullErrMsg()
	}
	return ce.fullErrMsg
}

func (ce *myCrawlerError) genFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("crawler error: ")
	if ce.errType != "" {
		buffer.WriteString(string(ce.errType))
		buffer.WriteString(string(": "))
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = buffer.String()
}

func NewCrawlerError(errType ErrorType, errMsg string) CrawlerError {
	return &myCrawlerError{
		errType: errType,
		errMsg:  strings.TrimSpace(errMsg),
	}
}
func NewCrawlerErrorBy(errType ErrorType, err error) CrawlerError {
	return NewCrawlerError(errType, err.Error())
}

type IllegalParameterError struct {
	msg string
}

func (ipe IllegalParameterError) Error() string {
	return ipe.msg
}

func NewIllegalParameterError(msg string) IllegalParameterError {
	return IllegalParameterError{
		msg: fmt.Sprintf("illegal parameter: %s", strings.TrimSpace(msg)),
	}
}
