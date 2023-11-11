//go:build !nikandfor_loc_unsafe
// +build !nikandfor_loc_unsafe

package loc

import (
	"runtime"
)

func (l PC) nameFileLine() (name, file string, line int) {
	fs := callersFrames(PCs{l})
	f, _ := fs.Next()
	return f.Function, f.File, f.Line
}

func (l PC) FuncEntry() PC {
	if l == 0 {
		return 0
	}

	f := runtime.FuncForPC(uintptr(l))
	if f == nil {
		return 0
	}

	return PC(f.Entry())
}
