package loc

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestFmt(t *testing.T) {
	t.Logf("[%d]", 1000)
	t.Logf("[%6d]", 1000)
	t.Logf("[%-6d]", 1000)

	t.Logf("%v", Caller(0))
}

func TestLocationFormat(t *testing.T) {
	l := PC(0x1234cd)

	//	name, file, line := l.nameFileLine()
	//	t.Logf("location: %v %v %v", name, file, line)

	SetCache(l, "github.com/nikandfor/loc.Caller", "github.com/nikandfor/loc/location.go", 26)

	var b bytes.Buffer

	fmt.Fprintf(&b, "%v", l)
	assert.Equal(t, "location.go:26", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%.3v", l)
	assert.Equal(t, "location.go: 26", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%18.3v", l)
	assert.Equal(t, "location.go   : 26", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%18.30v", l)
	assert.Len(t, b.String(), 18)

	b.Reset()

	fmt.Fprintf(&b, "%10.1v", l)
	assert.Equal(t, "locatio:26", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%18.1v", l)
	assert.Len(t, b.String(), 18)

	b.Reset()

	fmt.Fprintf(&b, "%-18.3v", l)
	assert.Equal(t, "   location.go: 26", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%+v", l)
	assert.True(t, regexp.MustCompile(`[\w./-]*location.go:26`).MatchString(b.String()), "got %v", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%n", l)
	assert.Equal(t, "loc.Caller", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%12n", l)
	assert.Equal(t, "loc.Caller  ", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%-12s", l)
	assert.Equal(t, "  loc.Caller", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%f", l)
	assert.Equal(t, "location.go", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%12f", l)
	assert.Equal(t, "location.go ", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%d", l)
	assert.Equal(t, "26", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%0100d", l)
	assert.Len(t, b.String(), 20)

	b.Reset()

	fmt.Fprintf(&b, "%4l", l)
	assert.Equal(t, "  26", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%x", l)
	assert.Equal(t, "0x1234cd", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%8x", l)
	assert.Equal(t, "  0x1234cd", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%010X", l)
	assert.Equal(t, "0x00001234CD", b.String())

	b.Reset()

	fmt.Fprintf(&b, "%010x", l)
	assert.Equal(t, fmt.Sprintf("%010p", unsafe.Pointer(uintptr(l))), b.String())

	b.Reset()

	fmt.Fprintf(&b, "%100x", l)
	assert.Len(t, b.String(), 20)
}

func BenchmarkLocationString(b *testing.B) {
	b.ReportAllocs()

	l := Caller(0)

	for i := 0; i < b.N; i++ {
		_ = l.String()
	}
}

func BenchmarkLocationFormat(b *testing.B) {
	b.ReportAllocs()

	var s formatter
	s.flags['+'] = true

	l := Caller(0)

	for i := 0; i < b.N; i++ {
		s.Reset()

		l.Format(&s, 'v')
	}
}

type formatter struct {
	bytes.Buffer
	flags   [128]bool
	prec    int
	width   int
	precok  bool
	widthok bool
}

func (f *formatter) Flag(c int) bool {
	return f.flags[c]
}

func (f *formatter) Precision() (int, bool) {
	return f.prec, f.precok
}

func (f *formatter) Width() (int, bool) {
	return f.width, f.widthok
}
