package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
	// "time"
)

// Enums, basically
type BenchmarkType int
const (
	BenchmarkTypeA BenchmarkType = iota
	BenchmarkTypeB
	BenchmarkTypeC
	BenchmarkTypeCount // Make sure this is always last!
)
var benchmarkInfo = map[BenchmarkType]string{
	BenchmarkTypeA: "BenchmarkTypeA name",
	BenchmarkTypeB: "BenchmarkTypeB name",
	BenchmarkTypeC: "BenchmarkTypeC name",
}

/*
   CpuBenchmark keeps track of CPU cycles* we want to measure in the query tool.

   It uses the gotsc measuring instrument to count cycles. Cycle measurement is
   done with assembly RDTSC/P instructions (and other assembly around the RDTRC
   to make sure the reading is accurate).

   * TODO: This doesn't measure *actual* cycles (yet), because modern RDTSC calls
   use invariant TSC, which does not scale up or down as CPUs adjust frequency. To
   get a more accurate reading, we need to get a calibration of a known timer,
   which is not yet implemented.
*/
type CpuBenchmark struct {
	tscOverhead	uint64

	// Tracks cycle counts for loading db configs from the .env file.
	readEnvFile         []uint64
}

func NewCpuBenchmark(capacity int) *CpuBenchmark {
	// Capacity must be positive
	if capacity < 0 {
		capacity = 0
	}

	readEnvFile    := make([]uint64, 0, capacity)

	return &CpuBenchmark{
		// tscOverhead:  gotsc.TSCOverhead(),
		tscOverhead:  1,

		readEnvFile:         readEnvFile,
	}
}

// Adds cycle count to the list for the given type.
// NOTE: This will subtract TSC overhead for you!
func (b *CpuBenchmark) Add(bt BenchmarkType, cycles uint64) {
	count := cycles - b.tscOverhead

	switch bt {
	// Load db env file
	case BenchmarkTypeA:
		b.readEnvFile = append(b.readEnvFile, count)
	}
}

func (b *CpuBenchmark) Print() {
	// Compute benchmark totals & grand total
	var totalCycles uint64

	// Env file
	var totalEnvFileCycles uint64
	for _,q := range b.readEnvFile {
		totalEnvFileCycles += q
	}
	totalCycles += totalEnvFileCycles

	fmt.Println("[CPU cycle totals]")

	// We should have > 0 measurements, if not, this arch isn't supported & the library just
	// returned 0 for cycle counts.
	if totalCycles == 0 {
		fmt.Println("CPU architecture not supported, only x86 works currently")
		return
	}

	// Calculate percents
	loadEnvPercent        := 100 * float64(totalEnvFileCycles)        / float64(totalCycles)

	// Print with separators because cycle counts are large numbers
	width := 10
	p := message.NewPrinter(language.English)
	p.Printf( "   Load env file:  %*d  | %s%%\n", width, totalEnvFileCycles,        alignFloatAsStr(loadEnvPercent))
	p.Println(" =================================================")
	p.Printf( "    Total cycles:  %*d\n\n", width, totalCycles)

	p.Printf("* This is not quite fully accurate cycle measurement (see TODO at top of cpuBenchmark.go).")
	p.Printf("  BUT, it does give a window into how long parts of the app are taking.")
}

func alignFloatAsStr(n float64) string {
	str := strconv.FormatFloat(n, 'f', 2, 64)
	if n < 10 {
		str = " " + str
	}
	return str
}
