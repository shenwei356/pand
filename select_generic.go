// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

//go:build !amd64

package pand

var andInplaceGenerics = []andInplaceImpl{{andInplaceGeneric, "generic", true}}

var andGenerics = []andImpl{{andGeneric, "generic", true}}
