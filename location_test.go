package loc

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// padding
// padding
// padding

func TestLocation(t *testing.T) {
	testLocationInside(t)
}

func testLocationInside(t *testing.T) {
	t.Helper()

	pc := Caller(0)
	name, file, line := pc.NameFileLine()
	assert.Equal(t, "loc.testLocationInside", path.Base(name))
	assert.Equal(t, "location_test.go", filepath.Base(file))
	assert.Equal(t, 24, line)
}

func TestLocationShort(t *testing.T) {
	pc := Caller(0)
	assert.Equal(t, "location_test.go:32", pc.String())
}

func TestLocation2(t *testing.T) {
	func() {
		func() {
			l := FuncEntry(0)

			ver := runtime.Version()
			exp := "location_test.go:38"
			if strings.HasPrefix(ver, "go1.24") {
				exp = "location_test.go:36"
			}

			assert.Equal(t, exp, l.String(), "ver: %v", ver)
		}()
	}()
}

func TestLocationOnce(t *testing.T) {
	var pc PC

	CallerOnce(-1, &pc)
	assert.Equal(t, "location.go:44", pc.String())

	pc++
	save := pc

	CallerOnce(-1, &pc)

	assert.Equal(t, save, pc) // not changed

	//
	pc = 0

	FuncEntryOnce(-1, &pc)
	assert.Equal(t, "location.go:51", pc.String())

	pc++
	save = pc

	FuncEntryOnce(-1, &pc)

	assert.Equal(t, save, pc) // not changed
}

func TestLocationCropFileName(t *testing.T) {
	assert.Equal(t, "github.com/nikandfor/tlog/sub/module/file.go",
		cropFilename("/path/to/src/github.com/nikandfor/tlog/sub/module/file.go", "github.com/nikandfor/tlog/sub/module.(*type).method"))
	assert.Equal(t, "github.com/nikandfor/tlog/sub/module/file.go",
		cropFilename("/path/to/src/github.com/nikandfor/tlog/sub/module/file.go", "github.com/nikandfor/tlog/sub/module.method"))
	assert.Equal(t, "github.com/nikandfor/tlog/root.go", cropFilename("/path/to/src/github.com/nikandfor/tlog/root.go", "github.com/nikandfor/tlog.type.method"))
	assert.Equal(t, "github.com/nikandfor/tlog/root.go", cropFilename("/path/to/src/github.com/nikandfor/tlog/root.go", "github.com/nikandfor/tlog.method"))
	assert.Equal(t, "root.go", cropFilename("/path/to/src/root.go", "github.com/nikandfor/tlog.method"))
	assert.Equal(t, "sub/file.go", cropFilename("/path/to/src/sub/file.go", "github.com/nikandfor/tlog/sub.method"))
	assert.Equal(t, "root.go", cropFilename("/path/to/src/root.go", "tlog.method"))
	assert.Equal(t, "subpkg/file.go", cropFilename("/path/to/src/subpkg/file.go", "subpkg.method"))
	assert.Equal(t, "subpkg/file.go", cropFilename("/path/to/src/subpkg/file.go", "github.com/nikandfor/tlog/subpkg.(*type).method"))
	assert.Equal(t, "errors/fmt_test.go",
		cropFilename("/home/runner/work/errors/errors/fmt_test.go", "tlog.app/go/error.TestErrorFormatCaller"))
	assert.Equal(t, "jq/object_test.go", cropFilename("/Users/nik/nikandfor/jq/object_test.go", "nikand.dev/go/jq.TestObject"))
}

func TestCaller(t *testing.T) {
	a, b := Caller(0),
		Caller(0)

	//	assert.False(t, a == b, "%x == %x", uintptr(a), uintptr(b))
	assert.NotEqual(t, a, b)
}

func TestSetCache(t *testing.T) {
	l := PC(0x1234567890)

	assert.False(t, Cached(l))

	SetCache(l, "", "", 0)

	assert.False(t, Cached(l))

	assert.NotEqual(t, "file.go:10", l.String())

	SetCache(l, "Name", "file.go", 10)

	assert.True(t, Cached(l))

	assert.Equal(t, "file.go:10", l.String())

	SetCacheBytes(l, []byte("name"), []byte("file"), 11)

	name, file, line := l.NameFileLine()
	assert.Equal(t, "name", name)
	assert.Equal(t, "file", file)
	assert.Equal(t, 11, line)

	SetCacheBytes(l, nil, nil, 12)

	name, file, line = l.NameFileLine()
	assert.Equal(t, "", name)
	assert.Equal(t, "", file)
	assert.Equal(t, 12, line)

	SetCacheBytes(l, nil, nil, 0)

	assert.False(t, Cached(l))
}

func BenchmarkLocationCaller(b *testing.B) {
	b.ReportAllocs()

	var l PC

	for i := 0; i < b.N; i++ {
		l = Caller(0)
	}

	_ = l
}

func BenchmarkLocationNameFileLine(b *testing.B) {
	b.ReportAllocs()

	var n, f string
	var line int

	l := Caller(0)

	for i := 0; i < b.N; i++ {
		n, f, line = l.nameFileLine()
	}

	_, _, _ = n, f, line //nolint:dogsled
}
