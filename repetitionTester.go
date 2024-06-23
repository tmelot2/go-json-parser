package main

import (
	"fmt"
)

type TestMode int
const (
	TestMode_Uninitialized TestMode = iota
	TestMode_Testing
	TestMode_Completed
	TestMode_Error
)

type RepetitionTestResults struct {
	testCount uint64
	totalTime uint64
	minTime  uint64
	maxTime  uint64
}

type RepetitionTester struct {
	targetProcessedByteCount uint64
	cpuTimerFreq             uint64
	tryForTime               uint64
	testStartedAt            uint64

	test_mode                  TestMode
	printNewMinimums           bool
	openBlockCount             uint32
	closeBlockCount            uint32
	timeAccumulatedOnThisTest  uint32
	bytesAccumulatedOnThisTest uint32

	results RepetitionTestResults
}

func (rt *RepetitionTester) NewTestWave(targetProcessedByteCount, cpuTimerFreq uint64) {
	if rt.test_mode == TestMode_Uninitialized {
		rt.test_mode = TestMode_Testing
		rt.targetProcessedByteCount = targetProcessedByteCount
		rt.cpuTimerFreq = cpuTimerFreq
	}
}

func newRepetitionTester() *RepetitionTester {
	return &RepetitionTester{ test_mode: TestMode_Uninitialized }
}

func main() {
	rt := newRepetitionTester()
	fmt.Println(rt)
}