//go:build !go1.20
// +build !go1.20

package loc

import "unsafe"

//nolint
type (
	funcID uint8

	funcInfo struct {
		entry *uintptr
		datap unsafe.Pointer
	}

	inlinedCall struct {
		parent   int16  // index of parent in the inltree, or < 0
		funcID   funcID // type of the called function
		_        byte
		file     int32 // fileno index into filetab
		line     int32 // line number of the call site
		func_    int32 // offset into pclntab for name of called function
		parentPc int32 // position of an instruction whose source position is the call site (offset from entry)
	}
)

// FuncEntry is functions entry point.
func (l PC) FuncEntry() PC {
	funcInfo := findfunc(l)
	if funcInfo.entry == nil {
		return 0
	}
	return funcInfoEntry(funcInfo)
}
