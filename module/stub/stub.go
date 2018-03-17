package stub

import (
	"fmt"
	"sync/atomic"

	"github.com/l-dandelion/webcrawler/errors"
	"github.com/l-dandelion/webcrawler/module"
	"gopcp.v2/helper/log"
)

// logger 代表日志记录器。
var logger = log.DLogger()

type myModule struct {
	mid             module.MID
	addr            string
	score           uint64
	scoreCalculator module.CalculateScore
	calledCount     uint64
	acceptedCount   uint64
	completedCount  uint64
	handlingNumber  uint64
}

func NewModuleInternal(mid module.MID, scoreCalculator module.CalculateScore) (ModuleInternal, error) {
	parts, err := module.SplitMID(mid)
	if err != nil {
		return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal ID %q: %s", mid, err))
	}
	return &myModule{
		mid:             mid,
		addr:            parts[2],
		scoreCalculator: scoreCalculator,
	}, nil
}

func (m *myModule) ID() module.MID {
	return m.mid
}

func (m *myModule) Addr() string {
	return m.addr
}

func (m *myModule) Score() uint64 {
	return atomic.LoadUint64(&m.score)
}

func (m *myModule) SetScore(score uint64) {
	atomic.StoreUint64(&m.score, score)
}

func (m *myModule) ScoreCalculator() module.CalculateScore {
	return m.scoreCalculator
}

func (m *myModule) CalledCount() uint64 {
	return atomic.LoadUint64(&m.calledCount)
}

func (m *myModule) AcceptedCount() uint64 {
	return atomic.LoadUint64(&m.acceptedCount)
}

func (m *myModule) CompletedCount() uint64 {
	return atomic.LoadUint64(&m.completedCount)
}

func (m *myModule) HandlingNumber() uint64 {
	return atomic.LoadUint64(&m.handlingNumber)
}

func (m *myModule) Counts() module.Counts {
	return module.Counts{
		CalledCount:    m.CalledCount(),
		AcceptedCount:  m.AcceptedCount(),
		CompletedCount: m.CompletedCount(),
		HandlingNumber: m.HandlingNumber(),
	}
}

func (m *myModule) Summary() module.SummaryStruct {
	return module.SummaryStruct{
		ID:        m.mid,
		Called:    m.CalledCount(),
		Accepted:  m.AcceptedCount(),
		Completed: m.CompletedCount(),
		Handling:  m.HandlingNumber(),
	}
}

func (m *myModule) IncrCalledCount() {
	atomic.AddUint64(&m.calledCount, 1)
}

func (m *myModule) IncrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, 1)
}

func (m *myModule) IncrCompletedCount() {
	atomic.AddUint64(&m.completedCount, 1)
}

func (m *myModule) IncrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, 1)
}

func (m *myModule) DecrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, ^uint64(0))
}

func (m *myModule) Clear() {
	atomic.StoreUint64(&m.calledCount, 0)
	atomic.StoreUint64(&m.acceptedCount, 0)
	atomic.StoreUint64(&m.completedCount, 0)
	atomic.StoreUint64(&m.handlingNumber, 0)
}
