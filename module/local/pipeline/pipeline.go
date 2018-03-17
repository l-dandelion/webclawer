package pipeline

import (
	"fmt"

	"github.com/l-dandelion/webcrawler/module"
	"github.com/l-dandelion/webcrawler/module/stub"
	"gopcp.v2/helper/log"
)

var logger = log.DLogger()

type myPipeline struct {
	stub.ModuleInternal
	itemProcessors []module.ProcessItem
	failFast       bool
}

func New(mid module.MID,
	itemProcessors []module.ProcessItem,
	scoreCalculator module.CalculateScore) (module.Pipeline, error) {
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, genError(err.Error())
	}
	if itemProcessors == nil {
		return nil, genParameterError("nil item processor list")
	}
	if len(itemProcessors) == 0 {
		return nil, genParameterError("empty item processor list")
	}
	var processors []module.ProcessItem
	for i, processor := range itemProcessors {
		if processor == nil {
			return nil, genParameterError(fmt.Sprintf("nil item processor[%d]", i))
		}
		processors = append(processors, processor)
	}
	return &myPipeline{
		ModuleInternal: moduleBase,
		itemProcessors: processors,
	}, nil
}

func (pipeline *myPipeline) ItemProcessors() []module.ProcessItem {
	processors := make([]module.ProcessItem, len(pipeline.itemProcessors))
	copy(processors, pipeline.itemProcessors)
	return processors
}

func (pipeline *myPipeline) Send(item module.Item) (errs []error) {
	pipeline.IncrHandlingNumber()
	defer pipeline.DecrHandlingNumber()
	pipeline.IncrCalledCount()
	if item == nil {
		errs = append(errs, genParameterError("nil item"))
		return
	}
	pipeline.IncrAcceptedCount()
	logger.Infof("Process item %+v... \n", item)
	currentItem := item
	for _, processor := range pipeline.itemProcessors {
		processedItem, err := processor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if pipeline.failFast {
				break
			}
		}
		if processedItem != nil {
			currentItem = processedItem
		}
	}
	if len(errs) == 0 {
		pipeline.IncrCompletedCount()
	}
	return
}

func (pipeline *myPipeline) FailFast() bool {
	return pipeline.failFast
}

func (pipeline *myPipeline) SetFailFast(failFast bool) {
	pipeline.failFast = failFast
}

type extraSummaryStruct struct {
	FailFast        bool `json:fail_fast`
	ProcessorNumber int  `json:"processor_number"`
}

func (pipeline *myPipeline) Summary() module.SummaryStruct {
	summary := pipeline.ModuleInternal.Summary()
	summary.Extra = extraSummaryStruct{
		FailFast:        pipeline.failFast,
		ProcessorNumber: len(pipeline.itemProcessors),
	}
	return summary
}
