// Stub profiler that doesn't do anything. May still have slight overhead, but I think Go compiles
// it out as part of dead code elimination. We'll see for sure when we start analyzing compiled asm.
//

// +build !profile

package profiler

var GlobalProfiler = newProfiler()
var GlobalProfilerParent string

type Profiler struct {}

func newProfiler() *Profiler { return &Profiler{} }
func (p *Profiler) BeginProfile() {}
func (p *Profiler) StartBlock(name string) {}
func (p *Profiler) EndBlock(name string) {}
func (p *Profiler) StartBandwidth(name string, byteCount uint64) {}
func (p *Profiler) EndBandwidth(name string) {}
func (p *Profiler) EndAndPrintProfile() {}
