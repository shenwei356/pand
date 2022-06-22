// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

//go:build !amd64

package pand

var andInplaceFuncs = []andInplaceImpl{
	{andInplaceGeneric, "generic", true},
}

var andFuncs = []andImpl{
	{andGeneric, "generic", true},
}
