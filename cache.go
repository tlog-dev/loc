package loc

import "sync"

type (
	nfl struct {
		name string
		file string
		line int
	}
)

var (
	locmu sync.Mutex
	locc  = map[PC]nfl{}
)

// NameFileLine returns function name, file and line number for location.
//
// This works only in the same binary where location was captured.
//
// This functions is a little bit modified version of runtime.(*Frames).Next().
func (l PC) NameFileLine() (name, file string, line int) {
	if l == 0 {
		return
	}

	locmu.Lock()
	c, ok := locc[l]
	locmu.Unlock()
	if ok {
		return c.name, c.file, c.line
	}

	name, file, line = l.nameFileLine()

	if file != "" {
		file = cropFilename(file, name)
	}

	locmu.Lock()
	locc[l] = nfl{
		name: name,
		file: file,
		line: line,
	}
	locmu.Unlock()

	return
}

// SetCache sets name, file and line for the PC.
// It allows to work with PC in another binary the same as in original.
func SetCache(l PC, name, file string, line int) {
	locmu.Lock()
	if name == "" && file == "" && line == 0 {
		delete(locc, l)
	} else {
		locc[l] = nfl{
			name: name,
			file: file,
			line: line,
		}
	}
	locmu.Unlock()
}

func SetCacheBytes(l PC, name, file []byte, line int) {
	locmu.Lock()
	if name == nil && file == nil && line == 0 {
		delete(locc, l)
	} else {
		x := locc[l]

		if x.line != line || x.name != string(name) || x.file != string(file) {
			locc[l] = nfl{
				name: string(name),
				file: string(file),
				line: line,
			}
		}
	}
	locmu.Unlock()
}

func Cached(l PC) (ok bool) {
	locmu.Lock()
	_, ok = locc[l]
	locmu.Unlock()
	return
}
