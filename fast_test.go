//+build amd64

package loc

import (
	"fmt"
	"runtime"
	"testing"
	"unsafe"
)

func stack(p uintptr, fa int) {
	st := -14 * 0x10
	f := 14 * 0x10

	if p&0xf != 0 {
		st -= 8
		f -= 8
	}

	println("-----------")

	for x := f; x != st; x -= 16 {
		var l string
		if x == 0 || x == -8 {
			if x == 0 {
				l = " <    "
			} else {
				l = "    < "
			}

			l += fmt.Sprintf("ptr %x  caller %x", p, Caller2(1))
		}

		p0 := unsafe.Pointer(p + uintptr(x))
		p1 := unsafe.Pointer(p + uintptr(x+8))

		var v0, v1 uintptr
		v0 = *(*uintptr)(unsafe.Pointer(p0))
		v1 = *(*uintptr)(unsafe.Pointer(p1))

		println(fmt.Sprintf("%x %16x %16x%s", p0, v0, v1, l))
	}
}

func TestFastCaller(t *testing.T) {
	for i := 0; i < 6; i++ {
		pc := f1(i, 0x999)
		f := runtime.FuncForPC(pc)

		var file string
		var line int
		if pc != 0 {
			file, line = f.FileLine(pc)
		}

		t.Logf("i %2d  pc %x  func: %v   file %v:%d", i, pc, f.Name(), file, line)
	}
}

func TestFastFuncentry(t *testing.T) {
	for i := 0; i < 6; i++ {
		pc := e1(i, 0x999)
		f := runtime.FuncForPC(pc)

		var file string
		var line int
		if pc != 0 {
			file, line = f.FileLine(pc)
		}

		t.Logf("i %2d  pc %x  func: %v   file %v:%d", i, pc, f.Name(), file, line)
	}
}

//go:noinline
func f0(t *testing.T) {
	//	defer stack(uintptr(unsafe.Pointer(&t)), 0x111)
	//	pc := f1(0, 0x10)

	//	t.Logf("pc %x  caller %x", pc, Caller2(0))

	pc := uintptr(FastCaller(1))

	c := Caller(1)
	t.Logf("caller: %x  %v", uintptr(c), c)

	f := runtime.FuncForPC(pc)

	t.Logf("FastCaller SP %x", fastCallerSP)
	t.Logf("pc %x  func: %v", pc, f.Name())
}

//go:noinline
func f1(s int, x int) (c uintptr) {
	//	defer stack(uintptr(unsafe.Pointer(&c)), s)

	return f2(s, x+1)
}

//go:noinline
func f2(s int, x int) (c uintptr) {
	//	defer stack(uintptr(unsafe.Pointer(&c)), s)

	return uintptr(FastCaller(s))
}

//go:noinline
func e1(s int, x int) (c uintptr) {
	return e2(s, x+1)
}

//go:noinline
func e2(s int, x int) (c uintptr) {
	return uintptr(FastFuncentry(s))
}

var (
	stackdumpOff uintptr
	stackdump    [100]uintptr
	stackdumpSt  int
	stackdumpF   int
	stackdumpX   int

	stp0, stp1 unsafe.Pointer
	stv0, stv1 uintptr

	stl string

	fastCallerSP uintptr
)

////go:noinline
//func FastCaller(s int) (c uintptr) {
//	return uintptr(fastCaller(s))
//}

func Caller2(s int) (r uintptr) {
	caller2(1+s, &r, 1, 1)

	return
}

//go:noescape
//go:linkname caller2 runtime.callers
func caller2(skip int, pc *uintptr, len, cap int) int

//go:nosplit
func add(p, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func Benchmark3FastCaller(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		func() {
			func() {
				_ = FastCaller(3)
			}()
		}()
	}
}

func Benchmark3Caller(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		func() {
			func() {
				_ = Caller(3)
			}()
		}()
	}
}

func Benchmark3RuntimeCaller(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		func() {
			func() {
				_, _, _, _ = runtime.Caller(3)
			}()
		}()
	}
}
