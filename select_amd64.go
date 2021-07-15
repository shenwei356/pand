// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

package pand

import "golang.org/x/sys/cpu"

var andFuncs = []andImpl{
	// {andAvx512, cpu.X86.HasBMI2 && cpu.X86.HasAVX512BW},
	{andAvx, "avx", cpu.X86.HasAVX},
	{andGeneric, "generic", true},
}
