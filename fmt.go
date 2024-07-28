package loc

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"unsafe"
)

type (
	buf []byte

	locFmtState struct {
		buf
		flags string
	}
)

var spaces = []byte("                                                                                                                                                                ") //nolint:lll

// String formats PC as base_name.go:line.
//
// Works only in the same binary where Caller of FuncEntry was called.
// Or if PC.SetCache was called.
func (l PC) String() string {
	_, file, line := l.NameFileLine()
	file = filepath.Base(file)

	b := append([]byte(file), ":        "...)

	i := len(file)
	n := 1 + width(line)

	b = b[:i+n]

	for q, j := line, n-1; j >= 1; j-- {
		b[i+j] = byte(q%10) + '0'
		q /= 10
	}

	return string(b)
}

// Format is fmt.Formatter interface implementation.
func (l PC) Format(s fmt.State, c rune) {
	switch c {
	default: // v
		l.formatV(s)
	case 'n', 's':
		l.formatName(s)
	case 'f':
		l.formatFile(s)
	case 'd', 'l':
		l.formatLine(s)
	case 'x', 'X', 'p', 'P':
		l.formatPC(s, c)
	}
}

func (l PC) formatV(s fmt.State) {
	name, file, line := l.NameFileLine()

	if !s.Flag('+') {
		file = filepath.Base(file)
		name = path.Base(name)
	}

	if s.Flag('#') {
		file = name
	}

	prec, ok := s.Precision()
	if lw := width(line); !ok || lw > prec {
		prec = lw
	}

	if prec > 20 {
		prec = 20
	}

	width, ok := s.Width()
	if !ok {
		width = len(file) + 1 + prec
	}

	fwidth := width - 1 - prec
	if fwidth < 0 {
		fwidth = 0
	}

	var bufdata [128]byte
	//	buf := bufdata[:0]
	buf := noescapeSlize(&bufdata[0], len(bufdata))

	buf = l.appendStr(buf, s, fwidth, file)

	if fwidth+1+prec > width {
		prec = width - fwidth - 1
	}

	buf = append(buf, ":                    "[:1+prec]...)

	for q, j := line, prec; q != 0 && j >= 0; j-- {
		buf[fwidth+j] = byte(q%10) + '0'
		q /= 10
	}

	_, _ = s.Write(buf)
}

func (l PC) formatName(s fmt.State) {
	name, _, _ := l.NameFileLine()

	if !s.Flag('+') {
		name = path.Base(name)
	}

	w, ok := s.Width()
	if !ok {
		w = len(name)
	}

	var bufdata [128]byte
	buf := noescapeSlize(&bufdata[0], len(bufdata))

	buf = l.appendStr(buf, s, w, name)

	_, _ = s.Write(buf)
}

func (l PC) formatFile(s fmt.State) {
	_, file, _ := l.NameFileLine()

	if !s.Flag('+') {
		file = filepath.Base(file)
	}

	w, ok := s.Width()
	if !ok {
		w = len(file)
	}

	var bufdata [128]byte
	buf := noescapeSlize(&bufdata[0], len(bufdata))

	buf = l.appendStr(buf, s, w, file)

	_, _ = s.Write(buf)
}

func (l PC) appendStr(buf []byte, s fmt.State, w int, name string) []byte {
	if w > len(name) {
		if s.Flag('-') {
			buf = append(buf, spaces[:w-len(name)]...)
		}

		buf = append(buf, name...)

		if !s.Flag('-') {
			buf = append(buf, spaces[:w-len(name)]...)
		}
	} else {
		buf = append(buf, name[:w]...)
	}

	return buf
}

func (l PC) formatLine(s fmt.State) {
	_, _, line := l.NameFileLine()

	lineW := width(line)

	w, ok := s.Width()
	if !ok || w < lineW {
		w = lineW
	}

	if w > 20 {
		w = 20
	}

	var bufdata [32]byte
	buf := noescapeSlize(&bufdata[0], len(bufdata))

	buf = append(buf, "                    "[:w]...)

	j := w - 1
	for q := line; q != 0 && j >= 0; j-- {
		buf[j] = byte(q%10) + '0'
		q /= 10
	}

	for j >= 1 && s.Flag('0') {
		buf[j] = '0'
		j--
	}

	_, _ = s.Write(buf)
}

func (l PC) formatPC(s fmt.State, c rune) {
	lineW := 1
	for x := l >> 4; x != 0; x >>= 4 {
		lineW++
	}

	w, ok := s.Width()
	if !ok || w < lineW {
		w = lineW
	}

	w += len("0x")

	if w > 20 {
		w = 20
	}

	var bufdata [32]byte
	buf := noescapeSlize(&bufdata[0], len(bufdata))

	buf = append(buf, "                    "[:w]...)

	const (
		hexc = "0123456789abcdef"
		hexC = "0123456789ABCDEF"
	)

	hex := hexc
	if c >= 'A' && c <= 'Z' {
		hex = hexC
	}

	j := w - 1
	for q := uint64(l); q != 0 && j >= 0; j-- {
		buf[j] = hex[q&0xf]
		q /= 16
	}

	for j > 1 && s.Flag('0') {
		buf[j] = '0'
		j--
	}

	if j > 0 {
		buf[j-1] = '0'
		buf[j] = 'x'
	}

	_, _ = s.Write(buf)
}

func width(v int) (n int) {
	n = 0
	for v != 0 {
		v /= 10
		n++
	}
	return
}

// String formats PCs as list of type_name (file.go:line)
//
// Works only in the same binary where Caller of FuncEntry was called.
// Or if PC.SetCache was called.
func (t PCs) String() string {
	var b buf

	for i, l := range t {
		if i != 0 {
			b = append(b, " at "...)
		}

		_, file, line := l.NameFileLine()
		file = filepath.Base(file)

		i := len(b) + len(file)
		n := 1 + width(line)

		b = append(b, file...)
		b = append(b, ":        "...)

		b = b[:i+n]

		for q, j := line, n-1; j >= 1; j-- {
			b[i+j] = byte(q%10) + '0'
			q /= 10
		}
	}

	return string(b)
}

// FormatString formats PCs as list of type_name (file.go:line)
//
// Works only in the same binary where Caller of FuncEntry was called.
// Or if PC.SetCache was called.
func (t PCs) FormatString(flags string) string {
	s := locFmtState{flags: flags}

	t.Format(&s, 'v')

	return string(s.buf)
}

func (t PCs) Format(s fmt.State, c rune) {
	switch {
	case s.Flag('+'):
		for _, l := range t {
			_, _ = s.Write([]byte("at "))
			l.Format(s, c)
			_, _ = s.Write([]byte("\n"))
		}
	default:
		for i, l := range t {
			if i != 0 {
				_, _ = s.Write([]byte(" at "))
			}
			l.Format(s, c)
		}
	}
}

func (s *locFmtState) Flag(c int) bool {
	for _, f := range s.flags {
		if f == rune(c) {
			return true
		}
	}

	return false
}

func (b *buf) Write(p []byte) (int, error) {
	*b = append(*b, p...)

	return len(p), nil
}

func (s *locFmtState) Width() (int, bool)     { return 0, false }
func (s *locFmtState) Precision() (int, bool) { return 0, false }

func noescapeSlize(b *byte, l int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{ //nolint:govet
		Data: uintptr(unsafe.Pointer(b)),
		Len:  0,
		Cap:  l,
	}))
}
