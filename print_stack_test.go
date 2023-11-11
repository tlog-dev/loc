package loc

import (
	"runtime"
	"testing"
)

func TestPrintStack(t *testing.T) {
	inline3var = inline3

	func() {
		inline2(t)
	}()
}

var inline3var func(*testing.T)

func inline2(t *testing.T) {
	inline3var(t)
}

func inline3(t *testing.T) {
	defer func() { //nolint:gocritic
		var pcsbuf [6]PC

		pcs := CallersFill(0, pcsbuf[:])

		for _, pc := range pcs {
			n, f, l := pc.NameFileLine()

			t.Logf("location %x  %v %v %v", uintptr(pc), n, f, l)
		}

		var fpcs [6]uintptr
		n := runtime.Callers(1, fpcs[:])

		fr := runtime.CallersFrames(fpcs[:n])

		for {
			f, more := fr.Next()

			t.Logf("runtime  %x  %v %v %v", f.PC, f.Function, f.File, f.Line)

			if !more {
				break
			}
		}
	}()
}
