// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

package pand

import "golang.org/x/sys/cpu"

var andInplaceFuncs = []andInplaceImpl{
	// {andAvx512, cpu.X86.HasBMI2 && cpu.X86.HasAVX512BW},
	{andInplaceAvx, "avx", cpu.X86.HasAVX2},
	{andInplaceGeneric, "generic", true},
}

var andFuncs = []andImpl{
	// {andAvx512, cpu.X86.HasBMI2 && cpu.X86.HasAVX512BW},
	{andAvx, "avx", cpu.X86.HasAVX2},
	{andGeneric, "generic", true},
}
