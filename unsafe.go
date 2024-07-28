package loc

import _ "unsafe"

//go:noescape
//go:linkname callers runtime.callers
func callers(skip int, pc []PC) int

//go:noescape
//go:linkname caller1 runtime.callers
func caller1(skip int, pc *PC, len, cap int) int //nolint:predeclared
