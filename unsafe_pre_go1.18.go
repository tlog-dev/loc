//go:build nikandfor_loc_unsafe && !go1.18
// +build nikandfor_loc_unsafe,!go1.18

package loc

func funcInfoEntry(f funcInfo) PC {
	return PC(*f.entry)
}
