package downloader

import (
	"net/http"

	"github.com/l-dandelion/webcrawler/module"
	"github.com/l-dandelion/webcrawler/module/stub"
	"gopcp.v2/helper/log"
)

var logger = log.DLogger()

type myDownloader struct {
	stub.ModuleInternal
	httpClient http.Client
}

func New(mid module.MID, client *http.Client, scoreCalculator module.CalculateScore) (module.Downloader, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, genError(err.Error())
	}
	if client == nil {
		return nil, genParameterError("nil http client")
	}
	return &myDownloader{
		ModuleInternal: moduleBase,
		httpClient:     *client,
	}, nil
}

func (downloader *myDownloader) Download(req *module.Request) (*module.Response, error) {
	downloader.IncrHandlingNumber()
	defer downloader.DecrHandlingNumber()
	downloader.IncrCalledCount()
	if req == nil {
		return nil, genParameterError("nil request")
	}
	httpReq := req.HTTPReq()
	if httpReq == nil {
		return nil, genParameterError("nil HTTP request")
	}
	downloader.IncrAcceptedCount()
	logger.Infof("Do the request (URL: %s, depth: %d)... \n", httpReq.URL, req.Depth())
	httpResp, err := downloader.httpClient.Do(httpReq)
	if err != nil {
		return nil, genError(err.Error())
	}
	downloader.IncrCompletedCount()
	return module.NewResponse(httpResp, req.Depth()), nil
}
