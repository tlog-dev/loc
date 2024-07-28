//go:build nikandfor_loc_unsafe
// +build nikandfor_loc_unsafe

package loc

import (
	"runtime"
	"unsafe"
)

func (l PC) nameFileLine() (name, file string, line int) {
	f := l.frame()

	return f.Function, f.File, f.Line
}

func (l PC) FuncEntry() PC {
	f := l.frame()

	return PC(f.Entry)
}

//go:nocheckptr
func (l PC) frame() runtimeFrame {
	fs0 := &runtimeFrames{}

	x := (uintptr)(unsafe.Pointer(fs0))
	fs := (*runtimeFrames)(unsafe.Pointer(x ^ 0))

	fs.buf = l
	fs.ptr = &fs.buf
	fs.len = 1
	fs.frames = fs.frameStore[:0]

	r := (*runtime.Frames)(unsafe.Pointer(fs))

	f, _ := r.Next()

	return f
}
