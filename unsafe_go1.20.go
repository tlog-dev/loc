package loc

import "unsafe"

type (
	funcID uint8

	funcInfo struct {
		*_func
		datap unsafe.Pointer
	}

	_func struct {
		entryOff uint32 // start pc, as offset from moduledata.text/pcHeader.textStart
		nameOff  int32  // function name, as index into moduledata.funcnametab.
	}

	// inlinedCall is the encoding of entries in the FUNCDATA_InlTree table.
	inlinedCall struct {
		funcID    funcID // type of the called function
		_         [3]byte
		nameOff   int32 // offset into pclntab for name of called function
		parentPc  int32 // position of an instruction whose source position is the call site (offset from entry)
		startLine int32 // line number of start of function (func keyword/TEXT directive)
	}
)

// FuncEntry is functions entry point.
func (l PC) FuncEntry() PC {
	funcInfo := findfunc(l)
	if funcInfo._func == nil {
		return 0
	}

	return funcInfoEntry(funcInfo)
}
