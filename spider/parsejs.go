package spider

import (
	"encoding/xml"

	"github.com/l-dandelion/webcrawler/module"
	"github.com/l-dandelion/webcrawler/spider/parsejs"
)

type (
	SpiderModle struct {
		Name                  string   `xml:"Name"`
		InitRequests          string   `xml:"Init>Script"`
		MaxDepth              uint32   `xml:"MaxDepth"`
		AcceptedPrimaryDomain []string `xml:"AcceptedPrimaryDomain>PrimaryDomain"`
		RespParsers           []string `xml:"ResponseParser>Script"`
		ItemProcessors        []string `xml:"ItemProcessor>Script"`

		ReqBufferCap         uint32 `xml:"DataArgs>Request>BufferCap"`
		ReqMaxBufferNumber   uint32 `xml:"DataArgs>Request>MaxBufferNumber"`
		RespBufferCap        uint32 `xml:"DataArgs>Response>BufferCap"`
		RespMaxBufferNumber  uint32 `xml:"DataArgs>Response>MaxBufferNumber"`
		ItemBufferCap        uint32 `xml:"DataArgs>Item>BufferCap"`
		ItemMaxBufferNumber  uint32 `xml:"DataArgs>Item>MaxBufferNumber"`
		ErrorBufferCap       uint32 `xml:"DataArgs>Error>BufferCap"`
		ErrorMaxBufferNumber uint32 `xml:"DataArgs>Error>MaxBufferNumber"`

		DownloaderNumber uint8 `xml:"ModuleArgs>DownloaderNumber"`
		AnalyzerNumber   uint8 `xml:"ModuleArgs>AnalyzerNumber"`
		PipelineNumber   uint8 `xml:"ModuleArgs>PipelineNumber"`
	}
)

func GetSpiderModuleByXML(xmlStr string) (*SpiderModle, error) {
	xmlByte := []byte(xmlStr)
	m := &SpiderModle{}
	err := xml.Unmarshal(xmlByte, m)
	return m, err
}

func GetSpiderByXML(xmlStr string) (*Spider, error) {
	spiderModule, err := GetSpiderModuleByXML(xmlStr)
	if err != nil {
		return nil, err
	}
	spider := &Spider{
		Name:                  spiderModule.Name,
		MaxDepth:              spiderModule.MaxDepth,
		AcceptedPrimaryDomain: spiderModule.AcceptedPrimaryDomain,

		ReqBufferCap:         spiderModule.ReqBufferCap,
		ReqMaxBufferNumber:   spiderModule.ReqMaxBufferNumber,
		RespBufferCap:        spiderModule.RespBufferCap,
		RespMaxBufferNumber:  spiderModule.RespMaxBufferNumber,
		ItemBufferCap:        spiderModule.ItemBufferCap,
		ItemMaxBufferNumber:  spiderModule.ItemMaxBufferNumber,
		ErrorBufferCap:       spiderModule.ErrorBufferCap,
		ErrorMaxBufferNumber: spiderModule.ErrorMaxBufferNumber,

		DownloaderNumber: spiderModule.DownloaderNumber,
		AnalyzerNumber:   spiderModule.DownloaderNumber,
		PipelineNumber:   spiderModule.PipelineNumber,
	}
	initialReqs, err := parsejs.GetReqsFromJS(spiderModule.InitRequests)
	if err != nil {
		return nil, err
	}
	respParsers := []module.ParseResponse{}
	for _, jsCode := range spiderModule.RespParsers {
		respParser := parsejs.GetParserFromJS(jsCode)
		respParsers = append(respParsers, respParser)
	}
	itemProcessors := []module.ProcessItem{}
	for _, jsCode := range spiderModule.ItemProcessors {
		itemProcessor := parsejs.GetProcessorFromJS(jsCode)
		itemProcessors = append(itemProcessors, itemProcessor)
	}

	spider.InitialReqs = initialReqs
	spider.RespParsers = respParsers
	spider.ItemProcessors = itemProcessors
	return spider, nil

}
