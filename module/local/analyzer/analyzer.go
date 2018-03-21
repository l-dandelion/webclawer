package analyzer

import (
	"fmt"

	"github.com/l-dandelion/webcrawler/module"
	"github.com/l-dandelion/webcrawler/module/stub"
	"github.com/l-dandelion/webcrawler/toolkit/reader"
	"gopcp.v2/helper/log"
)

var logger = log.DLogger()

type myAnalyzer struct {
	stub.ModuleInternal
	respParsers []module.ParseResponse
}

func New(
	mid module.MID,
	respParsers []module.ParseResponse,
	scoreCalculator module.CalculateScore) (module.Analyzer, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, genError(err.Error())
	}
	if respParsers == nil {
		return nil, genParameterError("nil response parsers")
	}
	if len(respParsers) == 0 {
		return nil, genParameterError("empty response parser list")
	}
	parsers := []module.ParseResponse{}
	for i, parser := range respParsers {
		if parser == nil {
			return nil, genParameterError(fmt.Sprintf("empty response parser[%d]", i))
		}
		parsers = append(parsers, parser)
	}
	return &myAnalyzer{
		ModuleInternal: moduleBase,
		respParsers:    parsers,
	}, nil
}

func (analyzer *myAnalyzer) RespParsers() []module.ParseResponse {
	parsers := make([]module.ParseResponse, len(analyzer.respParsers))
	copy(parsers, analyzer.respParsers)
	return parsers
}

func (analyzer *myAnalyzer) Analyze(resp *module.Response) (dataList []module.Data, errList []error) {
	analyzer.IncrHandlingNumber()
	defer analyzer.DecrHandlingNumber()
	analyzer.IncrCalledCount()
	if resp == nil {
		errList = append(errList, genParameterError("nil response"))
		return
	}
	httpResp := resp.HTTPResp()
	if httpResp == nil {
		errList = append(errList, genParameterError("nil HTTP response"))
		return
	}
	httpReq := httpResp.Request
	if httpReq == nil {
		errList = append(errList, genParameterError("nil HTTP request"))
		return
	}
	reqURL := httpReq.URL
	if reqURL == nil {
		errList = append(errList, genParameterError("nil HTTP request URL"))
		return
	}
	analyzer.IncrAcceptedCount()
	respDepth := resp.Depth()
	logger.Infof("Parse the response (URL: %s, depth: %d)... \n",
		reqURL, respDepth)
	if httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	multiReader, err := reader.NewMultipleReader(httpResp.Body)
	if err != nil {
		errList = append(errList, genError(err.Error()))
		return
	}
	parsers := analyzer.RespParsers()
	for _, parser := range parsers {
		httpResp.Body = multiReader.Reader()
		pDataList, pErrList := parser(httpResp, respDepth)
		if pDataList != nil {
			for _, pData := range pDataList {
				dataList = appendDataList(dataList, pData, respDepth)
			}
		}
		if pErrList != nil {
			for _, pErr := range pErrList {
				if pErr != nil {
					errList = append(errList, pErr)
				}
			}
		}
	}
	if len(errList) == 0 {
		analyzer.IncrCompletedCount()
	}
	return
}

func appendDataList(dataList []module.Data, data module.Data, respDepth uint32) []module.Data {
	if data == nil {
		return dataList
	}
	req, ok := data.(*module.Request)
	if !ok {
		return append(dataList, data)
	}
	newDepth := respDepth + 1
	if req.Depth() != newDepth {
		req = module.NewRequest(req.HTTPReq())
	}
	return append(dataList, req)
}
