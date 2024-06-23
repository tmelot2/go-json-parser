package repetitionTester

import (
	// "fmt"
	"math"

	"tmelot.jsonparser/internal/profiler"
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
	minTime   uint64
	maxTime   uint64
}

type RepetitionTester struct {
	targetProcessedByteCount uint64
	cpuTimerFreq             uint64
	tryForTime               uint64
	testStartedAt            uint64

	testMode                  TestMode
	printNewMinimums           bool
	openBlockCount             uint32
	closeBlockCount            uint32
	timeAccumulatedOnThisTest  uint32
	bytesAccumulatedOnThisTest uint32

	results RepetitionTestResults
}

func NewRepetitionTester() *RepetitionTester {
	return &RepetitionTester{
		testMode: TestMode_Uninitialized,
	}
}

func (rt *RepetitionTester) secondsFromCPUTime(cpuTime float64, cpuTimerFreq uint64) float64 {
	var result float64
	if cpuTimerFreq > 0 {
		result = cpuTime / float64(cpuTimerFreq)
	}
	return result
}

func (rt *RepetitionTester) NewTestWave(targetProcessedByteCount, cpuTimerFreq uint64, secondsToTry uint32) {
	if rt.testMode == TestMode_Uninitialized {
		rt.testMode = TestMode_Testing
		rt.targetProcessedByteCount = targetProcessedByteCount
		rt.cpuTimerFreq = cpuTimerFreq
		rt.printNewMinimums = true
		rt.results.minTime = math.MaxUint64 - 1
	} else if rt.testMode == TestMode_Completed {
		rt.testMode = TestMode_Testing

		if rt.targetProcessedByteCount != targetProcessedByteCount {
			panic("targetProcessedByteCount changed")
		}

		if rt.cpuTimerFreq != cpuTimerFreq {
			panic("cpuTimerFreq changed")
		}
	}

	rt.tryForTime = uint64(secondsToTry) * cpuTimerFreq
	rt.testStartedAt = profiler.ReadCPUTimer()
}
