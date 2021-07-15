// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

package pand

import "golang.org/x/sys/cpu"

var andInplaceFuncs = []andInplaceImpl{
	{andInplaceAvx512, "avx512", cpu.X86.HasAVX512F},
	{andInplaceAvx2, "avx2", cpu.X86.HasAVX2},
	{andInplaceGeneric, "generic", true},
}

var andFuncs = []andImpl{
	{andAvx512, "avx512", cpu.X86.HasAVX512F},
	{andAvx2, "avx2", cpu.X86.HasAVX2},
	{andGeneric, "generic", true},
}
