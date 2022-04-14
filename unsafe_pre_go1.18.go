//go:build !go1.18
// +build !go1.18

package loc

func funcInfoEntry(f funcInfo) PC {
	return PC(*f.entry)
}
