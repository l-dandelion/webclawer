package scheduler

import (
	"encoding/json"
	"sort"

	"github.com/l-dandelion/webcrawler/module"
	"github.com/l-dandelion/webcrawler/toolkit/buffer"
)

type SchedSummary interface {
	Struct() SummaryStruct
	String() string
}

type mySchedSummary struct {
	requestArgs RequestArgs
	dataArgs    DataArgs
	moduleArgs  ModuleArgs
	maxDepth    uint32
	sched       *myScheduler
}

func (ss *mySchedSummary) Struct() SummaryStruct {
	registrar := ss.sched.registrar
	return SummaryStruct{
		RequestArgs:     ss.requestArgs,
		DataArgs:        ss.dataArgs,
		ModuleArgs:      ss.moduleArgs.Summary(),
		Status:          GetStatusDescription(ss.sched.Status()),
		Downloaders:     getModuleSummaries(registrar, module.TYPE_DOWNLOADER),
		Analyzers:       getModuleSummaries(registrar, module.TYPE_ANALYZER),
		Pipelines:       getModuleSummaries(registrar, module.TYPE_PIPELINE),
		ReqBufferPool:   getBufferPoolSummary(ss.sched.reqBufferPool),
		RespBufferPool:  getBufferPoolSummary(ss.sched.respBufferPool),
		ItemBufferPool:  getBufferPoolSummary(ss.sched.itemBufferPool),
		ErrorBufferPool: getBufferPoolSummary(ss.sched.errorBufferPool),
		NumURL:          ss.sched.urlMap.Len(),
	}
}

func (ss *mySchedSummary) String() string {
	b, err := json.MarshalIndent(ss.Struct(), "", "    ")
	if err != nil {
		logger.Errorf("An error occurs when generating scheduler summary: %s\n", err)
		return ""
	}
	return string(b)
}

func newSchedSummary(requestArgs RequestArgs, dataArgs DataArgs, moduleArgs ModuleArgs, sched *myScheduler) SchedSummary {
	if sched == nil {
		return nil
	}
	return &mySchedSummary{
		requestArgs: requestArgs,
		dataArgs:    dataArgs,
		moduleArgs:  moduleArgs,
		maxDepth:    requestArgs.MaxDepth,
		sched:       sched,
	}
}

type SummaryStruct struct {
	RequestArgs     RequestArgs             `json:"request_args"`
	DataArgs        DataArgs                `json:"data_args"`
	ModuleArgs      ModuleArgsSummary       `json:"module_args"`
	Status          string                  `json:"status"`
	Downloaders     []module.SummaryStruct  `json:"downloaders"`
	Analyzers       []module.SummaryStruct  `json:"analyzers"`
	Pipelines       []module.SummaryStruct  `json:"pipelines"`
	ReqBufferPool   BufferPoolSummaryStruct `json:"request_buffer_pool"`
	RespBufferPool  BufferPoolSummaryStruct `json:"response_buffer_pool"`
	ItemBufferPool  BufferPoolSummaryStruct `json:"item_buffer_pool"`
	ErrorBufferPool BufferPoolSummaryStruct `json:"error_buffer_pool"`
	NumURL          uint64                  `json:"url_number"`
}

func (one *SummaryStruct) Same(anthor SummaryStruct) bool {
	if !one.RequestArgs.Same(&anthor.RequestArgs) {
		return false
	}
	if one.DataArgs != anthor.DataArgs ||
		one.ModuleArgs != anthor.ModuleArgs {
		return false
	}
	if one.Status != anthor.Status {
		return false
	}

	if anthor.Downloaders == nil || len(anthor.Downloaders) != len(one.Downloaders) {
		return false
	}
	for i, d := range anthor.Downloaders {
		if d != one.Downloaders[i] {
			return false
		}
	}

	if anthor.Analyzers == nil || len(anthor.Analyzers) != len(one.Analyzers) {
		return false
	}
	for i, a := range anthor.Analyzers {
		if a != one.Analyzers[i] {
			return false
		}
	}

	if anthor.Pipelines == nil || len(anthor.Pipelines) != len(one.Pipelines) {
		return false
	}
	for i, p := range anthor.Pipelines {
		if p != one.Pipelines[i] {
			return false
		}
	}

	if one.ReqBufferPool != anthor.ReqBufferPool ||
		one.RespBufferPool != anthor.RespBufferPool ||
		one.ItemBufferPool != anthor.ItemBufferPool ||
		one.ErrorBufferPool != anthor.ErrorBufferPool {
		return false
	}

	if one.NumURL != anthor.NumURL {
		return false
	}

	return true
}

type BufferPoolSummaryStruct struct {
	BufferCap       uint32 `json:"buffer_cap"`
	MaxBufferNumber uint32 `json:"max_buffer_number"`
	BufferNumber    uint32 `json:"buffer_number"`
	Total           uint64 `json:"total"`
}

func getBufferPoolSummary(bufferPool buffer.Pool) BufferPoolSummaryStruct {
	return BufferPoolSummaryStruct{
		BufferCap:       bufferPool.BufferCap(),
		MaxBufferNumber: bufferPool.MaxBufferNumber(),
		BufferNumber:    bufferPool.BufferNumber(),
		Total:           bufferPool.Total(),
	}
}

func getModuleSummaries(registrar module.Registrar, mType module.Type) []module.SummaryStruct {
	summaries := []module.SummaryStruct{}
	moduleMap, _ := registrar.GetAllByType(mType)
	if len(moduleMap) > 0 {
		for _, module := range moduleMap {
			summaries = append(summaries, module.Summary())
		}
	}
	if len(summaries) > 1 {
		sort.Slice(summaries, func(i, j int) bool {
			return summaries[i].ID < summaries[j].ID
		})
	}
	return summaries
}
