// Used to profile code repeatedly to find the "stars align" best-possible speed over
// many runs.

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

type RepetitionValueType int
const (
	RepValue_TestCount RepetitionValueType = iota

	RepValue_CPUTimer
	RepValue_MemPageFaults
	RepValue_ByteCount

	RepValue_Count
)

type RepetitionValue struct {
	e [RepValue_Count]uint64
}

type RepetitionTestResults struct {
	// testCount uint64
	// totalTime uint64
	// minTime   uint64
	// maxTime   uint64
	total RepetitionValue
	min   RepetitionValue
	max   RepetitionValue
}

type RepetitionTester struct {
	targetProcessedByteCount uint64
	cpuTimerFreq             uint64
	tryForTime               uint64
	testStartedAt            uint64

	testMode                        TestMode
	printNewMinimums                bool
	openBlockCount                  uint64
	closeBlockCount                 uint64

	accumulatedOnThisTest RepetitionValue
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
		rt.results.min.e[RepValue_CPUTimer] = math.MaxUint64 - 1
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

	accum := &rt.accumulatedOnThisTest
	accum.e[RepValue_MemPageFaults] -= GetPageFaultCount()
	accum.e[RepValue_CPUTimer] -= profiler.ReadCPUTimer()
}

func (rt *RepetitionTester) EndTime() {
	rt.closeBlockCount += 1

	accum := &rt.accumulatedOnThisTest
	accum.e[RepValue_CPUTimer] += profiler.ReadCPUTimer()
	accum.e[RepValue_MemPageFaults] += GetPageFaultCount()
}

func (rt *RepetitionTester) CountBytes(byteCount uint64) {
	accum := &rt.accumulatedOnThisTest
	// rt.bytesAccumulatedOnThisTest += byteCount
	accum.e[RepValue_ByteCount] += byteCount
}

func (rt *RepetitionTester) IsTesting() bool {
	if rt.testMode == TestMode_Testing {
		accum := rt.accumulatedOnThisTest
		currentTime := profiler.ReadCPUTimer()

		if rt.openBlockCount > 0 {
			// Error if blocks are unbalanced or there's a byte count mismatch.
			if rt.openBlockCount != rt.closeBlockCount {
				panic("Unbalanced begin/end time blocks")
			}
			// if rt.bytesAccumulatedOnThisTest != rt.targetProcessedByteCount {
			if accum.e[RepValue_ByteCount] != rt.targetProcessedByteCount {
				panic("Processed byte count mismatch")
			}

			if rt.testMode == TestMode_Testing {
				results := &rt.results

				// Increment results.
				accum.e[RepValue_TestCount] = 1
				for eIndex := 0; eIndex < len(accum.e); eIndex++ {
					results.total.e[eIndex] += accum.e[eIndex]
				}

				// Set new max or min if found.
				if results.max.e[RepValue_CPUTimer] < accum.e[RepValue_CPUTimer] {
					results.max = accum
				}
				if results.min.e[RepValue_CPUTimer] > accum.e[RepValue_CPUTimer] {
					results.min = accum

					// New min time found, restart to full trial time.
					rt.testStartedAt = currentTime

					if rt.printNewMinimums {
						rt.PrintValue("Min", results.min, rt.cpuTimerFreq)
						fmt.Printf("                                   \r")
					}
				}

				rt.openBlockCount = 0
				rt.closeBlockCount = 0
				rt.accumulatedOnThisTest = RepetitionValue{}
			}
		}

		if (currentTime - rt.testStartedAt) > rt.tryForTime {
			rt.testMode = TestMode_Completed
			fmt.Printf("                                                          \r")
			rt.PrintResults(rt.results, rt.cpuTimerFreq)
		}
	}

	result := rt.testMode == TestMode_Testing
	return result
}

func (rt *RepetitionTester) PrintValue(label string, value RepetitionValue, cpuTimerFreq uint64) {
	var divisor float64
	testCount := value.e[RepValue_TestCount]
	if testCount > 0 {
		divisor = float64(testCount)
	} else {
		divisor = float64(1)
	}

	var e [RepValue_Count]float64
	for eIndex := 0; eIndex < len(e); eIndex++ {
		e[eIndex] = float64(value.e[eIndex]) / divisor
	}

	fmt.Printf("%s: %.0f", label, e[RepValue_CPUTimer])
	if cpuTimerFreq > 0 {
		seconds := rt.secondsFromCPUTime(e[RepValue_CPUTimer], cpuTimerFreq)
		fmt.Printf(" (%fms)", 1000.0*seconds)

		if e[RepValue_ByteCount] > 0 {
			gigabyte := float64(1024.0 * 1024.0 * 1024.0)
			bandwidth := float64(e[RepValue_ByteCount]) / (gigabyte * seconds)
			fmt.Printf(" %fgb/s", bandwidth)
		}
	}

	pfs := uint64(e[RepValue_MemPageFaults])
    if pfs > 0 {
        fmt.Printf(" PF: %d (%0.4fk/fault)", pfs, float64(e[RepValue_ByteCount]) / (float64(pfs) * 1024.0))
    }
}

func (rt *RepetitionTester) PrintResults(results RepetitionTestResults, cpuTimerFreq uint64) {
	rt.PrintValue("Min", results.min, cpuTimerFreq)
	fmt.Println("")
	rt.PrintValue("Max", results.max, cpuTimerFreq)
	fmt.Println("")
	rt.PrintValue("Avg", results.total, cpuTimerFreq)
	fmt.Println("")
	fmt.Printf("Test Count: %d\n", results.total.e[RepValue_TestCount])
	fmt.Println("")
}
