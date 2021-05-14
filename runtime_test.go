package loc

import (
	"runtime"
	"testing"
)

func BenchmarkRuntime(b *testing.B) {
	b.ReportAllocs()

	var pc uintptr
	var ok bool
	for i := 0; i < b.N; i++ {
		pc, _, _, ok = runtime.Caller(0)

		f := runtime.FuncForPC(pc)

		_ = f.Name()
	}

	if !ok {
		b.Errorf("not ok")
	}
}
