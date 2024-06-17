package main

import (
	"fmt"
	"strings"
)

var MAX_BLOCKS = 4096

// Holds running data for a profiled Block.
type Block struct {
	startTSC uint64 // We only ever store the start TSC.
	total    uint64 // Accumulates on every EndBlock() call.
	hitCount uint64 // Accumulates +1 for every StartBlock() call.
}

// Special bucket used to profile the profiler.
var PROFILER_BLOCK_NAME = "__Profiler"

type Profiler struct {
	// Maps block names to blocks.
	blocks map[string]Block
	// Block for profiling the profiler (yo dawg).
	profilerBlock *Block
	// Remembers ordering of blocks by name.
	order  []string
	// Remembers when the profiler started so we can compute total time.
	startTSC uint64
}

func newProfiler() *Profiler {
	return &Profiler{
		blocks:        make(map[string]Block, MAX_BLOCKS),
		profilerBlock: &Block{},
		order:         make([]string, 0, MAX_BLOCKS),
		startTSC:      0,
	}
}

func (p *Profiler) BeginProfile() {
	p.startTSC = ReadCPUTimer()
}

// Starts a new time stamp counter for the given block name.
// NOTE: Will replace startTSC for the given block.
func (p *Profiler) StartBlock(name string) {
	profStart := ReadCPUTimer()

	if name == PROFILER_BLOCK_NAME {
		return
	}

	// Get block & add to ordering. It uses same names as blocks, so check with a hash
	// lookup instead of an array scan.
	block, ok := p.blocks[name]
	if !ok {
		if len(p.order) > MAX_BLOCKS {
			msg := fmt.Sprintf("Number of blocks has exceeded maximum of %d", MAX_BLOCKS)
			panic(msg)
		}
		p.order = append(p.order, name)
	}

	// Timestamp start of measurement.
	block.hitCount += 1
	block.startTSC = ReadCPUTimer()
	p.blocks[name] = block

	// Profile the profiler.
	p.profilerBlock.total += ReadCPUTimer() - profStart
	p.profilerBlock.hitCount += 1
}

// Ends a time stamp counter & accumulates its total duration.
func (p *Profiler) EndBlock(name string) {
	profStart := ReadCPUTimer()

	// Ignore if block name does not exist.
	block, ok := p.blocks[name]
	if !ok {
		return
	}

	// Calc duration, return if negative.
	duration := ReadCPUTimer() - block.startTSC
	if duration < 0 {
		return
	}

	// Update total & write back.
	block.total += duration
	p.blocks[name] = block

	// Profile the profiler.
	p.profilerBlock.total += ReadCPUTimer() - profStart
	p.profilerBlock.hitCount += 1
}

// Prints block names & their durations in the order that Start() was called for each block.
func (p *Profiler) EndAndPrintProfile() {
	totalCycles := ReadCPUTimer() - p.startTSC

	// Print CPU info
	printer := GetPrinter()
	cpuFreq := EstimateCPUTimerFreq(false)
	ms := 1000.0 * float64(totalCycles) / float64(cpuFreq)
	fmt.Println("\n[CPU profiling stats]")
	printer.Printf("Total time: %0.4fms (CPU freq %*d)\n", ms, 14, cpuFreq)

	// Print block profiles.
	p.printBlockHeader()
	for _, blockName := range p.order {
		p.printBlockTimeElapsed(blockName, p.blocks[blockName].hitCount, p.blocks[blockName].total, totalCycles)
	}

	// Print profiler block profile
	p.printBlockTimeElapsed(PROFILER_BLOCK_NAME[2:], p.profilerBlock.hitCount, p.profilerBlock.total, totalCycles)

	// Print total
	fmt.Println(strings.Repeat("=", 60))
	p.printBlockTimeElapsed("Total", 1, totalCycles, totalCycles)
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

func (p *Profiler) printBlockHeader() {
	printer := GetPrinter()
	longestBlockNameLen := p.getLongestBlockNameLen() + 1
	printer.Printf("  %*s %21s | %*s | %s\n", longestBlockNameLen, "Block", "Cycles", 8, "Hit Cnt", "Percent")
}

// Prints profiling information for the given block data.
func (p *Profiler) printBlockTimeElapsed(label string, hitCount, durationCycles, totalCycles uint64) {
	printer := GetPrinter()
	durationCyclesStr := printer.Sprintf("%*d", 20, durationCycles)
	longestBlockNameLen := p.getLongestBlockNameLen() + 1
	percent := float64(100.0 * (float64(durationCycles) / float64(totalCycles)))

	printer.Printf("  %*s: %14s | %*d | %.2f%%\n", longestBlockNameLen, label, durationCyclesStr, 8, hitCount, percent)
}
