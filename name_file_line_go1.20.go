//go:build go1.20

package loc

import "unsafe"

func (l PC) nameFileLine() (name, file string, line int) {
	funcInfo := findfunc(l)
	if funcInfo._func == nil {
		return
	}

	entry := funcInfoEntry(funcInfo)

	if l > entry {
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
		inltree := (*[1 << 20]inlinedCall)(inldata)
		// Non-strict as cgoTraceback may have added bogus PCs
		// with a valid funcInfo but invalid PCDATA.
		ix := pcdatavalue1(funcInfo, _PCDATA_InlTreeIndex, l, nil, false)
		if ix >= 0 {
			// Note: entry is not modified. It always refers to a real frame, not an inlined one.
			ic := inltree[ix]
			name = funcnameFromNameOff(funcInfo, ic.nameOff)
			// File/line from funcline1 below are already correct.
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

//go:linkname funcnameFromNameOff runtime.funcnameFromNameOff
func funcnameFromNameOff(f funcInfo, nameoff int32) string

//go:linkname funcInfoEntry runtime.funcInfo.entry
func funcInfoEntry(f funcInfo) PC
