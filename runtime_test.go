package loc

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntimeCaller(t *testing.T) {
	rpc, rfile, rline, ok := runtime.Caller(0)

	pc := Caller(0)

	f := runtime.FuncForPC(rpc)
	rname := f.Name()
	name, file, line := pc.NameFileLine()

	require.True(t, ok)

	assert.Equal(t, rname, name)
	assert.True(t, strings.HasSuffix(rfile, file))
	assert.Equal(t, rline, line-2)

	assert.Equal(t, f.Entry(), uintptr(pc.FuncEntry()))

	assert.Equal(t, f.Entry(), uintptr(FuncEntryFromFunc(TestRuntimeCaller)))
}

func TestRuntimeCallers(t *testing.T) {
	var rpcs [2]uintptr
	n := runtime.Callers(1, rpcs[:])

	pcs := Callers(0, 2)

	var rsum string

	frames := runtime.CallersFrames(rpcs[:n])
	i := 0
	for {
		fr, ok := frames.Next()

		name, file, line := pcs[i].NameFileLine()

		rline := fr.Line

		if i != 0 {
			rsum += " at "
		} else {
			rline += 2 // we called them from different lines
		}

		rsum += filepath.Base(fr.File) + ":" + strconv.Itoa(rline)

		assert.Equal(t, fr.Function, name)
		assert.True(t, strings.HasSuffix(fr.File, file))
		assert.Equal(t, rline, line)

		i++
		if !ok {
			break
		}
	}

	if t.Failed() {
		t.Logf("callers: %v", pcs)
		t.Logf("runtime: %v", rsum)
	}
}

func BenchmarkRuntimeCallerNameFileLine(b *testing.B) {
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

func BenchmarkRuntimeCallerFileLine(b *testing.B) {
	b.ReportAllocs()

	var ok bool
	for i := 0; i < b.N; i++ {
		_, _, _, ok = runtime.Caller(0)
	}

	if !ok {
		b.Errorf("not ok")
	}
}

func gover() string {
	return runtime.Version()
}
