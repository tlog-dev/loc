//go:build nikandfor_loc_unsafe && go1.21
// +build nikandfor_loc_unsafe,go1.21

package loc

import "unsafe"

type (
	inlineUnwinder struct {
		f       funcInfo
		cache   *uintptr
		inlTree *uintptr
	}

	inlineFrame struct {
		pc    uintptr
		index int32
	}

	srcFunc struct {
		datap     *uintptr
		nameOff   int32
		startLine int32
		funcID    funcID
	}
)

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

	// It's important that interpret pc non-strictly as cgoTraceback may
	// have added bogus PCs with a valid funcInfo but invalid PCDATA.
	u, uf := newInlineUnwinder(funcInfo, l, nil)
	sf := inlineUnwinder_srcFunc(&u, uf)

	name = funcNameForPrint(srcFunc_name(sf))
	file, line32 := funcline1(funcInfo, l, false)
	line = int(line32)

	return
}

//go:linkname findfunc runtime.findfunc
func findfunc(pc PC) funcInfo

//go:linkname funcInfoEntry runtime.funcInfo.entry
func funcInfoEntry(f funcInfo) PC

//go:linkname newInlineUnwinder runtime.newInlineUnwinder
func newInlineUnwinder(f funcInfo, pc PC, cache unsafe.Pointer) (inlineUnwinder, inlineFrame)

//go:linkname inlineUnwinder_srcFunc runtime.(*inlineUnwinder).srcFunc
func inlineUnwinder_srcFunc(*inlineUnwinder, inlineFrame) srcFunc

//go:linkname inlineUnwinder_isInlined runtime.(*inlineUnwinder).isInlined
func inlineUnwinder_isInlined(*inlineUnwinder, inlineFrame) bool

//go:linkname srcFunc_name runtime.srcFunc.name
func srcFunc_name(srcFunc) string

//go:linkname funcNameForPrint runtime.funcNameForPrint
func funcNameForPrint(name string) string

//go:linkname funcline1 runtime.funcline1
func funcline1(f funcInfo, targetpc PC, strict bool) (file string, line int32)
