// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

// +build !amd64

package pand

var andGenerics = []andImpl{{andGeneric, "generic", true}}
