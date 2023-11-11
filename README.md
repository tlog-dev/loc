[![Documentation](https://pkg.go.dev/badge/github.com/nikandfor/loc)](https://pkg.go.dev/github.com/nikandfor/loc?tab=doc)
[![Go workflow](https://github.com/nikandfor/json/actions/workflows/go.yml/badge.svg)](https://github.com/nikandfor/json/actions/workflows/go.yml)
[![CircleCI](https://circleci.com/gh/nikandfor/loc.svg?style=svg)](https://circleci.com/gh/nikandfor/loc)
[![codecov](https://codecov.io/gh/nikandfor/loc/tags/latest/graph/badge.svg)](https://codecov.io/gh/nikandfor/loc)
[![GolangCI](https://golangci.com/badges/github.com/nikandfor/loc.svg)](https://golangci.com/r/github.com/nikandfor/loc)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikandfor/loc)](https://goreportcard.com/report/github.com/nikandfor/loc)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/nikandfor/loc?sort=semver)

# loc

It's a fast, alloc-free and convinient version of `runtime.Caller`.

Performance benefits are available when using the `nikandfor_loc_unsafe` build tag. This relies on the internal runtime implementation, which means that older versions of the package may not compile with future versions of Go. Without the tag, it is safe to use with any version of Go.

It was born from [tlog](https://tlog.app/go/tlog).

## What is similar

Caller

```go
pc := loc.Caller(1)
ok := pc != 0
name, file, line := pc.NameFileLine()
e := pc.FuncEntry()

// is similar to

pc, file, line, ok := runtime.Caller(1) 
f := runtime.FuncForPC(pc)
name := f.Name()
e := f.Entry()
```

Callers

```go
pcs := loc.Callers(1, 3)
// or
var pcsbuf [3]loc.PC
pcs := loc.CallersFill(1, pcsbuf[:])

for _, pc := range pcs {
    name, file, file := pc.NameFileLine()
}

// is similar to
var pcbuf [3]uintptr
n := runtime.Callers(2, pcbuf[:])

frames := runtime.CallersFrames(pcbuf[:n])
for {
    frame, more := frames.Next()
    
    _, _, _ = frame.Function, frame.File, frame.Line

    if !more {
        break
    }
}
```

## What is different

### Normalized path

`loc` returns cropped filepath.
```
github.com/nikandfor/loc/func_test.go

# vs

/home/nik/nikandfor/loc/github.com/nikandfor/loc/func_test.go
```

### Performance

In `loc` the full cycle (get pc than name, file and line) takes 360=200+160 (200+0 when you repeat) ns whereas runtime takes 690 (640 without func name) + 2 allocs per each frame.

It's up to 3.5x improve.

```
BenchmarkLocationCaller-8              	 5844801	       201.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkLocationNameFileLine-8        	 7313388	       156.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkRuntimeCallerNameFileLine-8   	 1709940	       689.1 ns/op	     216 B/op	       2 allocs/op
BenchmarkRuntimeCallerFileLine-8       	 1917613	       642.1 ns/op	     216 B/op	       2 allocs/op
```
