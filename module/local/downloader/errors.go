package downloader

import (
	"github.com/l-dandelion/webcrawler/errors"
)

func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_DOWNLOADER, errMsg)
}

func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_DOWNLOADER, errors.NewIllegalParameterError(errMsg))
}
