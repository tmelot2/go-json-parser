package repetitionTester

import (
	"fmt"
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

	testMode                   TestMode
	printNewMinimums           bool
	openBlockCount             uint64
	closeBlockCount            uint64
	timeAccumulatedOnThisTest  uint64
	bytesAccumulatedOnThisTest uint64

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
			s := fmt.Sprintf("got %d, expected %d", cpuTimerFreq, rt.cpuTimerFreq)
			panic(s)
		}
	}

	rt.tryForTime = uint64(secondsToTry) * cpuTimerFreq
	rt.testStartedAt = profiler.ReadCPUTimer()
}

func (rt *RepetitionTester) BeginTime() {
	rt.openBlockCount += 1
	rt.timeAccumulatedOnThisTest -= profiler.ReadCPUTimer()
}

func (rt *RepetitionTester) EndTime() {
	rt.closeBlockCount += 1
	rt.timeAccumulatedOnThisTest += profiler.ReadCPUTimer()
}

func (rt *RepetitionTester) IsTesting() bool {
	if rt.testMode == TestMode_Testing {
		currentTime := profiler.ReadCPUTimer()

		if rt.openBlockCount > 0 {
			// Error if blocks are unbalanced or there's a byte count mismatch.
			if rt.openBlockCount != rt.closeBlockCount {
				panic("Unbalanced begin/end time blocks")
			}
			if rt.bytesAccumulatedOnThisTest != rt.targetProcessedByteCount {
				panic("Processed byte count mismatch")
			}

			if rt.testMode == TestMode_Testing {
				// Increment test stuff.
				elapsedTime := rt.timeAccumulatedOnThisTest
				results := &rt.results
				results.testCount += 1
				results.totalTime += elapsedTime

				// Set new max or min if found.
				if results.maxTime < elapsedTime {
					results.maxTime = elapsedTime
				}
				if results.minTime > elapsedTime {
					results.minTime = elapsedTime
					// New min time found, restart to full trial time.
					rt.testStartedAt = currentTime

					if rt.printNewMinimums {
						rt.PrintTime("Min", float64(results.minTime), rt.cpuTimerFreq, rt.bytesAccumulatedOnThisTest)
						fmt.Printf("                       \r")
					}
				}

				rt.openBlockCount = 0
				rt.closeBlockCount = 0
				rt.timeAccumulatedOnThisTest = 0
				rt.bytesAccumulatedOnThisTest = 0
			}
		}

		if (currentTime - rt.testStartedAt) > rt.tryForTime {
			rt.testMode = TestMode_Completed
			fmt.Printf("                                    \r")
			rt.PrintResults(rt.results, rt.cpuTimerFreq, rt.targetProcessedByteCount)
		}
	}

	result := rt.testMode == TestMode_Testing
	return result
}

func (rt *RepetitionTester) PrintTime(label string, cpuTime float64, cpuTimerFreq, byteCount uint64) {
	fmt.Printf("%s: %.0f", label, cpuTime)
	if cpuTimerFreq > 0 {
		seconds := rt.secondsFromCPUTime(cpuTime, cpuTimerFreq)
		fmt.Printf(" (%fms)", 1000.0*seconds)

		if byteCount > 0 {
			gigabyte := float64(1024.0 * 1024.0 * 1024.0)
			bestBandwidth := float64(byteCount) / (gigabyte * seconds)
			fmt.Printf(" %fgb/s", bestBandwidth)
		}
	}
}

func (rt *RepetitionTester) PrintResults(results RepetitionTestResults, cpuTimerFreq, byteCount uint64) {
	rt.PrintTime("Min", float64(results.minTime), cpuTimerFreq, byteCount)
	fmt.Println("")

	rt.PrintTime("Max", float64(results.maxTime), cpuTimerFreq, byteCount)
	fmt.Println("")

	if results.testCount > 0 {
		rt.PrintTime("Avg", float64(results.totalTime)/float64(results.testCount), cpuTimerFreq, byteCount)
		fmt.Println("")
	}
}
