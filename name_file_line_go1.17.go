//go:build go1.17
// +build go1.17

package loc

import "unsafe"

func (l PC) nameFileLine() (name, file string, line int) {
	if l == 0 {
		return
	}

	funcInfo := findfunc(l)
	if funcInfo.entry == nil {
		return
	}

	if uintptr(l) > *funcInfo.entry {
		// We store the pc of the start of the instruction following
		// the instruction in question (the call or the inline mark).
		// This is done for historical reasons, and to make FuncForPC
		// work correctly for entries in the result of runtime.Callers.
		l--
	}

	name = funcname(funcInfo)
	file, line32 := funcline1(funcInfo, l, false)
	line = int(line32)
	if inldata := funcdata(funcInfo, _FUNCDATA_InlTree); inldata != nil {
		ix := pcdatavalue1(funcInfo, _PCDATA_InlTreeIndex, l, nil, false)
		if ix >= 0 {
			inltree := (*[1 << 20]inlinedCall)(inldata)
			// Note: entry is not modified. It always refers to a real frame, not an inlined one.
			name = funcnameFromNameoff(funcInfo, inltree[ix].func_)
			// File/line is already correct.
			// TODO: remove file/line from InlinedCall?
		}
	}

	return
}

//go:linkname findfunc runtime.findfunc
func findfunc(pc PC) funcInfo

//go:linkname funcline1 runtime.funcline1
func funcline1(f funcInfo, targetpc PC, strict bool) (file string, line int32)

//go:linkname funcname runtime.funcname
func funcname(f funcInfo) string

//go:linkname funcdata runtime.funcdata
func funcdata(f funcInfo, i uint8) unsafe.Pointer

//go:linkname pcdatavalue runtime.pcdatavalue
func pcdatavalue(f funcInfo, table int32, targetpc PC, cache unsafe.Pointer) int32

//go:linkname pcdatavalue1 runtime.pcdatavalue1
func pcdatavalue1(f funcInfo, table int32, targetpc PC, cache unsafe.Pointer, strict bool) int32

//go:linkname funcnameFromNameoff runtime.funcnameFromNameoff
func funcnameFromNameoff(f funcInfo, nameoff int32) string
