package main

import (
	"fmt"
	"strings"
)

type TSC struct {
	start uint64
	end   uint64
}

type Profiler struct {
	// Maps block names to TSC start & end counters.
	tscs   map[string]TSC
	// Maps block names to total counted time from TSCs
	totals map[string]uint64
	// Remembers ordering of TSCs by block name
	order  []string
}

func newProfiler() *Profiler {
	return &Profiler{
		tscs:   make(map[string]TSC),
		totals: make(map[string]uint64),
		order:  make([]string, 0),
	}
}

// Starts a new time stamp counter for the given block name.
// NOTE: WILL ERASE existing data for the given block.
func (p *Profiler) StartBlock(name string) {
	// Add to ordering. It uses same names as TSCs, so check with a hash lookup instead of
	// an array scan.
	_, ok := p.tscs[name]
	if !ok {
		p.order = append(p.order, name)
	}

	// Timestamp start of measurement
	p.tscs[name] = TSC{
		start: ReadCPUTimer(),
		end:   0,
	}
}

// Ends a time stamp counter & accumulates its total duration.
func (p *Profiler) EndBlock(name string) {
	endTSC := ReadCPUTimer()

	// Ignore if TSC block name does not exist
	tscVal, ok := p.tscs[name]
	if !ok {
		return
	}

	// Set end TSC
	tscVal.end = endTSC
	p.tscs[name] = tscVal

	// Add sum for block if it doesn't exist
	totalVal, ok := p.totals[name]
	if !ok {
		p.totals[name] = 0
		totalVal = 0
	}

	// Calculate total
	if tscVal.start != 0 && tscVal.end != 0 {
		// TODO: Check that end > start
		duration := tscVal.end - tscVal.start
		totalVal += duration
		p.totals[name] = totalVal
	}
}

// Prints block names & their durations in the order that Start() was called for each block.
func (p *Profiler) Print() {
	// Compute total
	totalCycles := uint64(0)
	for _, t := range p.totals {
		totalCycles += t
	}

	// Print CPU info
	printer := GetPrinter()
	cpuFreq := EstimateCPUTimerFreq(false)
	ms := 1000.0 * float64(totalCycles) / float64(cpuFreq)
	fmt.Println("\n[CPU profiling stats]")
	printer.Printf("Total time: %0.4fms (CPU freq %*d)\n", ms, 14, cpuFreq)

	// Print block profiles
	for _, blockName := range p.order {
		p.printBlockTimeElapsed(blockName, p.totals[blockName], totalCycles)
	}

	fmt.Println(strings.Repeat("=", 50))

	// Print total
	p.printBlockTimeElapsed("Total", totalCycles, totalCycles)
}

// Returns len of longest block name string. Used for print formatting.
func (p *Profiler) getLongestBlockNameLen() int {
	longest := 0
	for _,b := range p.order {
		if len(b) > longest {
			longest = len(b)
		}
	}
	return longest
}

// Prints profiling information for the given block data.
func (p *Profiler) printBlockTimeElapsed(label string, durationCycles, totalCycles uint64) {
	printer := GetPrinter()
	durationCyclesStr := printer.Sprintf("%*d", 20, durationCycles)
	longestBlockNameLen := p.getLongestBlockNameLen() + 1
	percent := float64(100.0 * (float64(durationCycles) / float64(totalCycles)))

	printer.Printf("  %*s: %14s (%.2f%%)\n", longestBlockNameLen, label, durationCyclesStr, percent)
}
