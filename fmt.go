package loc

import (
	"fmt"
	"path"
	"path/filepath"
)

type (
	buf []byte

	locFmtState struct {
		buf
		flags string
	}
)

var spaces = []byte("                                                                                                                                                                ")

// String formats PC as base_name.go:line.
//
// Works only in the same binary where Caller of Funcentry was called.
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

	w2 := width(line)

	prec, ok := s.Precision()
	if !ok || prec < w2 {
		prec = w2
	}

	w, ok := s.Width()
	if !ok {
		w = len(file) + 1 + prec
	}

	w -= 1 + prec

	var bufdata [128]byte
	buf := bufdata[:0]

	if w > len(file) && s.Flag('-') {
		buf = append(buf, spaces[:w-len(file)]...)
	}

	w2 = w // reuse var
	if len(file) < w2 {
		w2 = len(file)
	}

	buf = append(buf, file[:w2]...)

	if w > len(file) && !s.Flag('-') {
		buf = append(buf, spaces[:w-len(file)]...)
	}

	w2 = len(buf) // reuse var
	buf = append(buf, ":        "[:1+prec]...)

	for q, j := line, prec; q != 0 && j >= 1; j-- {
		buf[w2+j] = byte(q%10) + '0'
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

	buf := l.appendStr(bufdata[:0], s, w, name)

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

	buf := l.appendStr(bufdata[:0], s, w, file)

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

	var bufdata [128]byte
	buf := append(bufdata[:0], "        "[:w]...)

	for q, j := line, w-1; q != 0 && j >= 0; j-- {
		buf[j] = byte(q%10) + '0'
		q /= 10
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
// Works only in the same binary where Caller of Funcentry was called.
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

// StringFlags formats PCs as list of type_name (file.go:line)
//
// Works only in the same binary where Caller of Funcentry was called.
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
			s.Write([]byte("at "))
			l.Format(s, c)
			s.Write([]byte("\n"))
		}
	default:
		for i, l := range t {
			if i != 0 {
				s.Write([]byte(" at "))
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
