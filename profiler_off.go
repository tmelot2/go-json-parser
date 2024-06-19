// Stub profiler that doesn't do anything. May still have slight overhead (we'll see when we start analyzing compiled asm).
//

// +build !profile

package main

type Profiler struct {}

func newProfiler() *Profiler { return &Profiler{} }
func (p *Profiler) BeginProfile() {}
func (p *Profiler) StartBlock(name string) {}
func (p *Profiler) EndBlock(name string) {}
func (p *Profiler) EndAndPrintProfile() {}
