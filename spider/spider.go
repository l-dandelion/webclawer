package spider

import (
	"github.com/l-dandelion/webcrawler/module"
	"github.com/l-dandelion/webcrawler/scheduler"
	lib "github.com/l-dandelion/webcrawler/spider/internal"
)

type Spider struct {
	Name                  string
	MaxDepth              uint32
	AcceptedPrimaryDomain []string
	InitialHTTPReq        []*module.Request
	RespParsers           []module.ParseResponse
	ItemProcessors        []module.ProcessItem
	ReqBufferCap          uint32
	ReqMaxBufferNumber    uint32
	RespBufferCap         uint32
	RespMaxBufferNumber   uint32
	ItemBufferCap         uint32
	ItemMaxBufferNumber   uint32
	ErrorBufferCap        uint32
	ErrorMaxBufferNumber  uint32
	DownloaderNumber      uint8
	AnalyzerNumber        uint8
	PipelineNumber        uint8
}

func NewSpider(
	name string,
	maxDepth uint32,
	acceptedPrimaryDomain []string,
	initialHTTPReq []*module.Request,
	RespParsers []module.ParseResponse,
	ItemProcessors []module.ProcessItem,
	DownloaderNumber, AnalyzerNumber, PipelineNumber uint8) *Spider {
	return &Spider{
		Name:                  name,
		MaxDepth:              maxDepth,
		AcceptedPrimaryDomain: acceptedPrimaryDomain,
		InitialHTTPReq:        initialHTTPReq,
		RespParsers:           RespParsers,
		ItemProcessors:        ItemProcessors,
		DownloaderNumber:      DownloaderNumber,
		AnalyzerNumber:        AnalyzerNumber,
		PipelineNumber:        PipelineNumber,
		ReqBufferCap:          50,
		ReqMaxBufferNumber:    1000,
		RespBufferCap:         50,
		RespMaxBufferNumber:   10,
		ItemBufferCap:         50,
		ItemMaxBufferNumber:   100,
		ErrorBufferCap:        50,
		ErrorMaxBufferNumber:  1,
	}
}

func (sp *Spider) GenAndStartScheduler() (scheduler.Scheduler, error) {
	sched := scheduler.NewScheduler()
	requestArgs := scheduler.RequestArgs{
		AcceptedDomains: sp.AcceptedPrimaryDomain,
		MaxDepth:        sp.MaxDepth,
	}

	dataArgs := scheduler.DataArgs{
		ReqBufferCap:         sp.ReqBufferCap,
		ReqMaxBufferNumber:   sp.ReqMaxBufferNumber,
		RespBufferCap:        sp.RespBufferCap,
		RespMaxBufferNumber:  sp.RespMaxBufferNumber,
		ItemBufferCap:        sp.ItemBufferCap,
		ItemMaxBufferNumber:  sp.ItemMaxBufferNumber,
		ErrorBufferCap:       sp.ErrorBufferCap,
		ErrorMaxBufferNumber: sp.ErrorMaxBufferNumber,
	}

	downloaders, err := lib.GetDownloaders(sp.DownloaderNumber)
	if err != nil {
		return nil, err
	}
	analyzers, err := lib.GetAnalyzers(sp.AnalyzerNumber, sp.RespParsers)
	if err != nil {
		return nil, err
	}
	pipelines, err := lib.GetPipelines(sp.PipelineNumber, sp.ItemProcessors)
	if err != nil {
		return nil, err
	}

	moduleArgs := scheduler.ModuleArgs{
		Downloaders: downloaders,
		Analyzers:   analyzers,
		Pipelines:   pipelines,
	}
	err = sched.Init(
		requestArgs,
		dataArgs,
		moduleArgs)
	if err != nil {
		return nil, err
	}

	err = sched.Start(sp.InitialHTTPReq)

	return sched, nil
}
