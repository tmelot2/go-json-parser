// Simple code profiler. Enable with build tags: go run -tags="profile".
//

// +build profile

package main

import (
	"fmt"
	"strings"
)

var MAX_BLOCKS = 4096

// Holds running data for a profiled Block.
type Block struct {
	startTSC      uint64 // We only ever store the start TSC.
	total         uint64 // Accumulates cycle count on EndBlock().
	totalChildren uint64 // Accumulates children cycle count on EndBlock(), so we can subtract out nested blocks.
	hitCount      uint64 // Accumulates +1 for every StartBlock() call.
	byteCount     uint64 // Accumulates bytes written by this block.
	parentName    string // Tracks parent block name to handle nested blocks.
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
func (p *Profiler) startBlock(name string, byteCount uint64) {
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

	// Setup measurement
	block.parentName = globalProfilerParent
	block.hitCount += 1
	block.byteCount += byteCount
	globalProfilerParent = name
	block.startTSC = ReadCPUTimer()
	p.blocks[name] = block

	// Profile the profiler.
	p.profilerBlock.total += ReadCPUTimer() - profStart
	p.profilerBlock.hitCount += 1
}

// Ends a time stamp counter & accumulates its total duration.
func (p *Profiler) endBlock(name string) {
	tsc := ReadCPUTimer()

	// Ignore if block name does not exist.
	block, ok := p.blocks[name]
	if !ok {
		return
	}

	// Unwind parent
	globalProfilerParent = block.parentName

	// Calc duration, return if negative.
	duration := tsc - block.startTSC
	if duration < 0 {
		return
	}

	// Update total & write back.
	block.total += duration
	p.blocks[name] = block

	// Update parent block
	parentBlock := p.blocks[block.parentName]
	parentBlock.totalChildren += duration
	p.blocks[block.parentName] = parentBlock

	// Profile the profiler.
	p.profilerBlock.total += ReadCPUTimer() - tsc
	p.profilerBlock.hitCount += 1
}

// Start a block measurement without tracking byte count.
func (p *Profiler) StartBlock(name string) {
	p.startBlock(name, 0)
}

// Start a block measurement and track byte count.
func (p *Profiler) StartBandwidth(name string, byteCount uint64) {
	p.startBlock(name, byteCount)
}

// End block.
func (p *Profiler) EndBlock(name string) {
	p.endBlock(name)
}

// End block.
func (p *Profiler) EndBandwidth(name string) {
	p.endBlock(name)
}

// Prints block names & their durations in the order that Start() was called for each block.
func (p *Profiler) EndAndPrintProfile() {
	totalCycles := ReadCPUTimer() - p.startTSC

	// Print CPU info
	printer := GetPrinter()
	cpuFreq := EstimateCPUTimerFreq(false)
	ms := 1000.0 * float64(totalCycles) / float64(cpuFreq)
	fmt.Println("\n[CPU profiling stats]")
	printer.Printf("Total time: %0.4fms (CPU freq %*d Hz)\n", ms, 14, cpuFreq)

	// Print block profiles.
	p.printBlockHeader()
	for _, blockName := range p.order {
		block := p.blocks[blockName]
		p.printBlockTimeElapsed(blockName, totalCycles, cpuFreq, &block)
	}

	// Print profiler & total
	fmt.Println(strings.Repeat("=", 60))
	p.printBlockTimeElapsed(PROFILER_BLOCK_NAME[2:], totalCycles, cpuFreq, p.profilerBlock)
	p.printBlockTimeElapsed("Total", totalCycles, cpuFreq, &Block{total: totalCycles})
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
func (p *Profiler) printBlockTimeElapsed(label string, totalCycles, timerFreq uint64, block *Block) {
	printer := GetPrinter()
	// func (p *Profiler) printBlockTimeElapsed(label string, hitCount, durationCycles, childrenCycles, totalCycles, byteCount uint64) {
	//   p.printBlockTimeElapsed(blockName, block.hitCount, block.total, block.totalChildren, totalCycles, block.byteCount)
	elapsedCycles := block.total - block.totalChildren
	durationCyclesStr := printer.Sprintf("%*d", 20, elapsedCycles)
	longestBlockNameLen := p.getLongestBlockNameLen() + 1
	percent := float64(100.0 * (float64(elapsedCycles) / float64(totalCycles)))

	printer.Printf("  %*s: %14s | %*d | %.2f%%", longestBlockNameLen, label, durationCyclesStr, 8, block.hitCount, percent)

	if block.totalChildren > 0 {
		percentWithChildren := float64(100.0 * (float64(block.total) / float64(totalCycles)))
		printer.Printf(", %.2f%% w/children", percentWithChildren)
	}

	if block.byteCount > 0 {
		megabyte := float64(1024.0) * float64(1024.0)
		gigabyte := megabyte * float64(1024.0)

		seconds := float64(elapsedCycles) / float64(timerFreq)
		bytesPerSecond := float64(block.byteCount) / seconds
		megabytes := float64(block.byteCount) / megabyte
		gigabytesPerSecond := bytesPerSecond / gigabyte

		printer.Printf(" [%d, %2.2f, %2.2f] ", timerFreq, seconds, bytesPerSecond)
		printer.Printf(", %0.3fmb at %.2fgb/s", megabytes, gigabytesPerSecond)
	}

	printer.Println("")
}
