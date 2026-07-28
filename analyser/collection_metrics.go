package analyser

import (
	"github.com/HugoDrl/LogParser/parser"
	"sort"
	"time"
)

func (m *CollectionMetric) handleService(log *parser.Log) {
	s, ok := m.ServicePerformance[log.Service]
	if !ok {
		s = ServiceMetric{
			Name: log.Service,
		}
	}
	s.Lines++

	s.AverageDuration = s.AverageDuration * time.Duration(s.Lines-1)
	s.AverageDuration += log.Duration
	s.AverageDuration /= time.Duration(s.Lines)
	m.ServicePerformance[log.Service] = s
}

func (m *CollectionMetric) handleSlowestLogs(slowestLogsToRetrieve int, log *parser.Log) {
	if slowestLogsToRetrieve == 0 {
		return
	}
	if len(m.SlowestInput) < slowestLogsToRetrieve {
		m.SlowestInput = append(m.SlowestInput, log)
		return
	}
	sort.Slice(m.SlowestInput, func(i, j int) bool {
		return m.SlowestInput[i].Duration > m.SlowestInput[j].Duration
	})

	if log.Duration > m.SlowestInput[len(m.SlowestInput)-1].Duration {
		m.SlowestInput = append(m.SlowestInput[:len(m.SlowestInput)-1], log)
	}
}
