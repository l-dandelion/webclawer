package module

import (
	"net/http"
)

//计数
type Counts struct {
	CalledCount    uint64 //调用次数
	AcceptedCount  uint64 //接受次数
	CompletedCount uint64 //完成次数
	HandlingNumber uint64 //实时处理数量
}

//摘要
type SummaryStruct struct {
	ID        MID         `json:"id"`
	Called    uint64      `json:"called"`
	Accepted  uint64      `json:"accepted"`
	Completed uint64      `json:"completed"`
	Handling  uint64      `json:"handling"`
	Extra     interface{} `json:"extra,omitempty"`
}

//组件的基础接口类型
type Module interface {
	ID() MID
	Addr() string
	Score() uint64
	SetScore(score uint64)
	ScoreCalculator() CalculateScore
	CalledCount() uint64
	AcceptedCount() uint64
	CompletedCount() uint64
	HandlingNumber() uint64
	Counts() Counts
	Summary() SummaryStruct
}

type Downloader interface {
	Module
	Download(req *Request) (*Response, error)
}

type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]Data, []error)

type Analyzer interface {
	Module
	RespParsers() []ParseResponse
	Analyze(resp *Response) ([]Data, []error)
}

type ProcessItem func(item Item) (result Item, err error)

type Pipeline interface {
	Module
	ItemProcessors() []ProcessItem
	Send(item Item) []error
	FailFast() bool
	SetFailFast(failFast bool)
}
