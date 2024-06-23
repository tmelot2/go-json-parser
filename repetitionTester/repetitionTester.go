package main

import (
	"fmt"
	"math"
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

	test_mode                  TestMode
	printNewMinimums           bool
	openBlockCount             uint32
	closeBlockCount            uint32
	timeAccumulatedOnThisTest  uint32
	bytesAccumulatedOnThisTest uint32

	results RepetitionTestResults
}

func newRepetitionTester() *RepetitionTester {
	return &RepetitionTester{
		test_mode: TestMode_Uninitialized,
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
	if rt.test_mode == TestMode_Uninitialized {
		rt.test_mode = TestMode_Testing
		rt.targetProcessedByteCount = targetProcessedByteCount
		rt.cpuTimerFreq = cpuTimerFreq
		rt.printNewMinimums = true
		rt.results.minTime = math.MaxUint64 - 1
	} else if rt.test_Mode == TestMode_Completed {
		rt.test_mode = TestMode_Testing

		if rt.targetProcessedByteCount != targetProcessedByteCount {
			panic("targetProcessedByteCount changed")
		}

		if rt.cpuTimerFreq != cpuTimerFreq {
			panic("cpuTimerFreq changed")
		}
	}

	rt.tryForTime = secondsToTry * cpuTimerFreq
	rt.testStartedAt = ReadCPUTimer()
}


func main() {
	rt := newRepetitionTester()
	fmt.Println(rt)
}
