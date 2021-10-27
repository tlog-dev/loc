package loc

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {
	t.Parallel()

	testLocationInside(t)
}

func testLocationInside(t *testing.T) {
	t.Helper()

	pc := Caller(0)
	name, file, line := pc.NameFileLine()
	assert.Equal(t, "loc.testLocationInside", path.Base(name))
	assert.Equal(t, "location_test.go", filepath.Base(file))
	assert.Equal(t, 20, line)
}

func TestLocationShort(t *testing.T) {
	t.Parallel()

	pc := Caller(0)
	assert.Equal(t, "location_test.go:30", pc.String())
}

func TestLocation2(t *testing.T) {
	t.Parallel()

	func() {
		func() {
			l := Funcentry(0)

			assert.Equal(t, "location_test.go:38", l.String())
		}()
	}()
}

func TestLocationOnce(t *testing.T) {
	t.Parallel()

	var pc PC

	CallerOnce(-1, &pc)
	assert.Equal(t, "location.go:44", pc.String())

	pc++
	save := pc

	CallerOnce(-1, &pc)

	assert.Equal(t, save, pc) // not changed

	//
	pc = 0

	FuncentryOnce(-1, &pc)
	assert.Equal(t, "location.go:51", pc.String())

	pc++
	save = pc

	FuncentryOnce(-1, &pc)

	assert.Equal(t, save, pc) // not changed
}

func TestLocationCropFileName(t *testing.T) {
	t.Parallel()

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
}

func TestCaller(t *testing.T) {
	t.Parallel()

	a, b := Caller(0),
		Caller(0)

	assert.False(t, a == b, "%x == %x", uintptr(a), uintptr(b))
}

func TestSetCache(t *testing.T) {
	t.Parallel()

	l := PC(0x1234567890)

	assert.False(t, Cached(l))

	SetCache(l, "", "", 0)

	assert.False(t, Cached(l))

	assert.NotEqual(t, "file.go:10", l.String())

	SetCache(l, "Name", "file.go", 10)

	assert.True(t, Cached(l))

	assert.Equal(t, "file.go:10", l.String())
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
