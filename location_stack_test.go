package loc

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationFillCallers(t *testing.T) {
	st := make(PCs, 1)

	st = CallersFill(0, st)

	assert.Len(t, st, 1)
	assert.Equal(t, "location_stack_test.go:16", st[0].String())
}

func testLocationsInside() (st PCs) {
	func() {
		func() {
			st = Callers(1, 3)
		}()
	}()

	return
}

func TestLocationPCsString(t *testing.T) {
	var st PCs
	func() {
		func() {
			st = testLocationsInside()
		}()
	}()

	assert.Len(t, st, 3)
	assert.Equal(t, "location_stack_test.go:26", st[0].String())
	assert.Equal(t, "location_stack_test.go:27", st[1].String())
	assert.Equal(t, "location_stack_test.go:36", st[2].String())

	re := `location_stack_test.go:26 at location_stack_test.go:27 at location_stack_test.go:36`

	assert.Equal(t, re, st.String())
}

func TestLocationPCsFormat(t *testing.T) {
	var st PCs
	func() {
		func() {
			st = testLocationsInside()
		}()
	}()

	assert.Equal(t, "location_stack_test.go:26 at location_stack_test.go:27 at location_stack_test.go:54", fmt.Sprintf("%v", st))

	addAllSubs := innerFuncName(Caller(0), 2)
	t.Logf("go version: %q: %q", gover(), addAllSubs)

	assert.Equal(t, "loc.testLocationsInside.func1:26 at loc.testLocationsInside:27 at loc.TestLocationPCsFormat"+addAllSubs+":54", fmt.Sprintf("%#v", st))

	re := `at [\w.-/]*location_stack_test.go:26
at [\w.-/]*location_stack_test.go:27
at [\w.-/]*location_stack_test.go:54
`
	v := fmt.Sprintf("%+v", st)
	assert.True(t, regexp.MustCompile(re).MatchString(v), "expected:\n%vgot:\n%v", re, v)
}

func TestLocationPCsFormatString(t *testing.T) {
	var st PCs
	func() {
		func() {
			st = testLocationsInside()
		}()
	}()

	assert.Equal(t, "location_stack_test.go:26 at location_stack_test.go:27 at location_stack_test.go:77", st.FormatString(""))

	addAllSubs := innerFuncName(Caller(0), 2)
	t.Logf("all sub funs suffix (go ver %q): %q", gover(), addAllSubs)

	assert.Equal(t, "loc.testLocationsInside.func1:26 at loc.testLocationsInside:27 at loc.TestLocationPCsFormatString"+addAllSubs+":77", st.FormatString("#"))

	re := `at [\w.-/]*location_stack_test.go:26
at [\w.-/]*location_stack_test.go:27
at [\w.-/]*location_stack_test.go:77
`

	v := st.FormatString("+")
	assert.True(t, regexp.MustCompile(re).MatchString(v), "expected:\n%vgot:\n%v", re, v)
}

func innerFuncName(fn PC, n int) string {
	var s string

	switch {
	//	case regexp.MustCompile("go1.16.*").MatchString(gover()):
	//		return ".func1"
	case regexp.MustCompile("go1.21.*").MatchString(gover()):
		name, _, _ := fn.NameFileLine()
		name = path.Base(name)
		name = name[strings.IndexByte(name, '.')+1:]

		s = "." + name

		for i := 0; i < n; i++ {
			s += fmt.Sprintf(".func%v", i+1)
		}
	default:
		s = ".func"

		for i := 0; i < n; i++ {
			if i != 0 {
				s += "."
			}

			s += fmt.Sprintf("%v", 1)
		}
	}

	return s
}
