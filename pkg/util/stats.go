package util

import "github.com/wklken/httptest/pkg/log"

const tableTPL = `
┌─────────────────────────┬─────────────────┬─────────────────┬─────────────────┐
│                         │           total │          passed │          failed │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│                   cases │          %6d │          %6d │          %6d │
├─────────────────────────┼─────────────────┼─────────────────┼─────────────────┤
│              assertions │          %6d │          %6d │          %6d │
├─────────────────────────┴─────────────────┴─────────────────┴─────────────────┤
│                    Time : %6d ms                                           │
└───────────────────────────────────────────────────────────────────────────────┘`

type Stats struct {
	okCaseCount     int64
	failCaseCount   int64
	okAssertCount   int64
	failAssertCount int64
}

func (s *Stats) MergeAssertCount(s1 Stats) {
	// NOTE: here only
	s.okAssertCount += s1.okAssertCount
	s.failAssertCount += s1.failAssertCount

	// if got fail assert, the case is fail
	if s1.AllPassed() {
		s.okCaseCount++
	} else {
		s.failCaseCount++
	}
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

func (s *Stats) AllPassed() bool {
	return s.failCaseCount == 0 && s.failAssertCount == 0
}

func (s *Stats) Report(totalCaseCount int, latency int64) {
	log.Info(tableTPL,
		totalCaseCount, s.okCaseCount, s.failCaseCount,
		s.okAssertCount+s.failAssertCount, s.okAssertCount, s.failAssertCount,
		latency)
}
