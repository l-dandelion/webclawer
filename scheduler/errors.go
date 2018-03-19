package scheduler

import (
	"github.com/l-dandelion/webcrawler/errors"
	"github.com/l-dandelion/webcrawler/module"
	"github.com/l-dandelion/webcrawler/toolkit/buffer"
)

func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER, errMsg)
}

func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_SCHEDULER, errors.NewIllegalParameterError(errMsg))
}

func genErrorByError(err error) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER, err.Error())
}

func sendError(err error, mid module.MID, errorBufferPool buffer.Pool) bool {
	if err == nil || errorBufferPool == nil || errorBufferPool.Closed() {
		return false
	}
	var (
		crawlerError errors.CrawlerError
		ok           bool
	)
	crawlerError, ok = err.(errors.CrawlerError)
	if !ok {
		var (
			moduleType module.Type
			errorType  errors.ErrorType
		)
		ok, moduleType = module.GetType(mid)
		if !ok {
			errorType = errors.ERROR_TYPE_SCHEDULER
		} else {
			switch moduleType {
			case module.TYPE_DOWNLOADER:
				errorType = errors.ERROR_TYPE_DOWNLOADER
			case module.TYPE_ANALYZER:
				errorType = errors.ERROR_TYPE_ANALYZER
			case module.TYPE_PIPELINE:
				errorType = errors.ERROR_TYPE_PIPELINE
			}
		}
		crawlerError = errors.NewCrawlerErrorBy(errorType, err)
	}
	if errorBufferPool.Closed() {
		return false
	}
	go func(crawlerError errors.CrawlerError) {
		if err := errorBufferPool.Put(crawlerError); err != nil {
			logger.Warnln("The error buffer pool was closed. Ignore error sending.")
		}
	}(crawlerError)
	return true
}
