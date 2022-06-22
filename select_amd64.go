// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

//go:build amd64

package pand

import "golang.org/x/sys/cpu"

var andInplaceFuncs = []andInplaceImpl{
	{andInplaceAVX512, "AVX512", cpu.X86.HasAVX512F},
	{andInplaceAVX2, "AVX2", cpu.X86.HasAVX2},
	{andInplaceSSE2, "SSE2", cpu.X86.HasSSE2},
	{andInplaceGeneric, "generic", true},
}

var andFuncs = []andImpl{
	{andAVX512, "AVX512", cpu.X86.HasAVX512F},
	{andAVX2, "AVX2", cpu.X86.HasAVX2},
	{andSSE2, "SSE2", cpu.X86.HasSSE2},
	{andGeneric, "generic", true},
}
