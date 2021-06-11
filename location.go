package loc

import (
	"path/filepath"
	"strings"
	"sync/atomic"
	"unsafe"
)

type (
	// PC is a program counter alias.
	// Function name, file name and line can be obtained from it but only in the same binary where Caller or Funcentry was called.
	PC uintptr

	// PCs is a stack trace.
	// It's quiet the same as runtime.CallerFrames but more efficient.
	PCs []PC
)

// Caller returns information about the calling goroutine's stack. The argument s is the number of frames to ascend, with 0 identifying the caller of Caller.
//
// It's hacked version of runtime.Caller with no allocs.
func Caller(s int) (r PC) {
	caller1(1+s, &r, 1, 1)

	return
}

// Funcentry returns information about the calling goroutine's stack. The argument s is the number of frames to ascend, with 0 identifying the caller of Caller.
//
// It's hacked version of runtime.Callers -> runtime.CallersFrames -> Frames.Next -> Frame.Entry with no allocs.
func Funcentry(s int) (r PC) {
	caller1(1+s, &r, 1, 1)

	return r.Entry()
}

func CallerOnce(s int, pc *PC) (r PC) {
	r = PC(atomic.LoadUintptr((*uintptr)(unsafe.Pointer(pc))))

	// TODO: may be mutex to prevent concurrention?

	if r == 0 {
		caller1(1+s, &r, 1, 1)

		atomic.StoreUintptr((*uintptr)(unsafe.Pointer(pc)), uintptr(r))
	}

	return
}

func FuncentryOnce(s int, pc *PC) (r PC) {
	r = PC(atomic.LoadUintptr((*uintptr)(unsafe.Pointer(pc))))

	if r == 0 {
		caller1(1+s, &r, 1, 1)

		r = r.Entry()

		atomic.StoreUintptr((*uintptr)(unsafe.Pointer(pc)), uintptr(r))
	}

	return
}

// Callers returns callers stack trace.
//
// It's hacked version of runtime.Callers -> runtime.CallersFrames -> Frames.Next -> Frame.Entry with only one alloc (resulting slice).
func Callers(skip, n int) PCs {
	tr := make([]PC, n)
	n = callers(1+skip, tr)
	return tr[:n]
}

// CallersFill puts callers stack trace into provided slice.
//
// It's hacked version of runtime.Callers -> runtime.CallersFrames -> Frames.Next -> Frame.Entry with no allocs.
func CallersFill(skip int, tr PCs) PCs {
	n := callers(1+skip, tr)
	return tr[:n]
}

func cropFilename(fn, tp string) string {
	p := strings.LastIndexByte(tp, '/')
	pp := strings.IndexByte(tp[p+1:], '.')
	tp = tp[:p+pp]

again:
	if p = strings.Index(fn, tp); p != -1 {
		return fn[p:]
	}

	p = strings.IndexByte(tp, '/')
	if p == -1 {
		return filepath.Base(fn)
	}
	tp = tp[p+1:]
	goto again
}
