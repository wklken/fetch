package util

import (
	"fmt"
	"sync"

	"github.com/wklken/httptest/pkg/log"
)

const tableTPL = `
┌─────────────────────────┬─────────────────┬─────────────────┬─────────────────┐
│                         │           total │          passed │          failed │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│                   files │          %6d │          %6d │          %6d │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│                   cases │          %6d │          %6d │          %6d │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│              assertions │          %6d │          %6d │          %6d │
├─────────────────────────┴─────────────────┴─────────────────┴─────────────────┤
│                    Time : %6d ms                                           │
└───────────────────────────────────────────────────────────────────────────────┘`

type Message struct {
	Type string
	Text string
}

type Stats struct {
	okFileCount     int64
	failFileCount   int64
	okCaseCount     int64
	failCaseCount   int64
	okAssertCount   int64
	failAssertCount int64
	messages        []Message
}

func (s *Stats) MergeAssertCount(s1 Stats) {
	// NOTE: here only
	s.okAssertCount += s1.okAssertCount
	s.failAssertCount += s1.failAssertCount

	messages := s1.GetMessages()
	if len(messages) > 0 {
		s.messages = append(s.messages, messages...)
	}
}

func (s *Stats) MergeAssertAndCaseCount(s1 Stats) {
	// NOTE: here only
	s.okAssertCount += s1.okAssertCount
	s.failAssertCount += s1.failAssertCount
	s.okCaseCount += s1.okCaseCount
	s.failCaseCount += s1.failCaseCount

	if s1.AllPassed() {
		s.okFileCount++
	} else {
		s.failFileCount++
	}

	messages := s1.GetMessages()
	if len(messages) > 0 {
		s.messages = append(s.messages, messages...)
	}
}

func (s *Stats) IncrOkFileCount() {
	s.okFileCount++
}

func (s *Stats) IncrFailFileCount() {
	s.failFileCount++
}

func (s *Stats) IncrOkCaseCount() {
	s.okCaseCount++
}

func (s *Stats) IncrFailCaseCount() {
	s.failCaseCount++
}

func (s *Stats) IncrOkAssertCount() {
	s.okAssertCount++
}

func (s *Stats) IncrFailAssertCount() {
	s.failAssertCount++
}

func (s *Stats) IncrFailAssertCountByN(n int64) {
	s.failAssertCount += n
}

func (s *Stats) GetOkAssertCount() int64 {
	return s.okAssertCount
}

func (s *Stats) GetFailAssertCount() int64 {
	return s.failAssertCount
}

func (s *Stats) GetOkCaseCount() int64 {
	return s.okCaseCount
}

func (s *Stats) GetFailCaseCount() int64 {
	return s.failCaseCount
}

func (s *Stats) AllPassed() bool {
	return s.failCaseCount == 0 && s.failAssertCount == 0 && s.failFileCount == 0
}

func (s *Stats) Report(latency int64) {
	log.Info(
		tableTPL,
		s.okFileCount+s.failFileCount,
		s.okFileCount,
		s.failFileCount,
		s.okCaseCount+s.failCaseCount,
		s.okCaseCount,
		s.failCaseCount,
		s.okAssertCount+s.failAssertCount,
		s.okAssertCount,
		s.failAssertCount,
		latency,
	)
}

func (s *Stats) AddMessage(msg Message) {
	s.messages = append(s.messages, msg)
}

func (s *Stats) GetMessages() []Message {
	return s.messages
}

func (s *Stats) AddTipMessage(format string, args ...interface{}) {
	s.messages = append(s.messages, Message{
		Type: "tip",
		Text: fmt.Sprintf(format, args...),
	})
}

func (s *Stats) AddErrorMessage(format string, args ...interface{}) {
	s.messages = append(s.messages, Message{
		Type: "error",
		Text: fmt.Sprintf(format, args...),
	})
}

func (s *Stats) AddInfofMessage(format string, args ...interface{}) {
	s.messages = append(s.messages, Message{
		Type: "infof",
		Text: fmt.Sprintf(format, args...),
	})
}

func (s *Stats) AddInfoMessage(format string, args ...interface{}) {
	s.messages = append(s.messages, Message{
		Type: "info",
		Text: fmt.Sprintf(format, args...),
	})
}

func (s *Stats) AddPassMessage() {
	s.messages = append(s.messages, Message{
		Type: "pass",
	})
}

func (s *Stats) AddFailMessage(format string, args ...interface{}) {
	s.messages = append(s.messages, Message{
		Type: "fail",
		Text: fmt.Sprintf(format, args...),
	})
}

func (s *Stats) PrintMessages() {
	for _, msg := range s.messages {
		switch msg.Type {
		case "tip":
			log.Tip(msg.Text)
		case "error":
			log.Error(msg.Text)
		case "infof":
			log.Infof(msg.Text)
		case "info":
			log.Info(msg.Text)
		case "pass":
			log.Pass()
		case "fail":
			log.Fail(msg.Text)
		default:
			log.Info(msg.Text)
		}
	}
}

type StatsCollection struct {
	stats Stats
	mu    sync.Mutex
}

func (sc *StatsCollection) Add(s Stats) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.stats.MergeAssertAndCaseCount(s)
}

func (sc *StatsCollection) GetStats() Stats {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	return sc.stats
}
