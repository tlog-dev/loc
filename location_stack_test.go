package loc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationFillCallers(t *testing.T) {
	st := make(PCs, 1)

	st = CallersFill(0, st)

	assert.Len(t, st, 1)
	assert.Equal(t, "location_stack_test.go:14", st[0].String())
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
	assert.Equal(t, "location_stack_test.go:24", st[0].String())
	assert.Equal(t, "location_stack_test.go:25", st[1].String())
	assert.Equal(t, "location_stack_test.go:34", st[2].String())

	re := `location_stack_test.go:24 at location_stack_test.go:25 at location_stack_test.go:34`

	assert.Equal(t, re, st.String())
}

func TestLocationPCsFormat(t *testing.T) {
	var st PCs
	func() {
		func() {
			st = testLocationsInside()
		}()
	}()

	assert.Equal(t, "location_stack_test.go:24 at location_stack_test.go:25 at location_stack_test.go:52", fmt.Sprintf("%v", st))

	addAllSubs := ""
	if regexp.MustCompile("go1.1[^6].*").MatchString(gover()) {
		addAllSubs = ".1"
	}

	t.Logf("go version: %q: %q", gover(), addAllSubs)

	assert.Equal(t, "loc.testLocationsInside.func1:24 at loc.testLocationsInside:25 at loc.TestLocationPCsFormat.func1"+addAllSubs+":52", fmt.Sprintf("%#v", st))

	re := `at [\w.-/]*location_stack_test.go:24
at [\w.-/]*location_stack_test.go:25
at [\w.-/]*location_stack_test.go:52
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

	assert.Equal(t, "location_stack_test.go:24 at location_stack_test.go:25 at location_stack_test.go:79", st.FormatString(""))

	addAllSubs := ""
	if regexp.MustCompile("go1.1[^6].*").MatchString(gover()) {
		addAllSubs = ".1"
	}

	t.Logf("all sub funs suffix (go ver %q): %q", gover(), addAllSubs)

	assert.Equal(t, "loc.testLocationsInside.func1:24 at loc.testLocationsInside:25 at loc.TestLocationPCsFormatString.func1"+addAllSubs+":79", st.FormatString("#"))

	re := `at [\w.-/]*location_stack_test.go:24
at [\w.-/]*location_stack_test.go:25
at [\w.-/]*location_stack_test.go:79
`

	v := st.FormatString("+")
	assert.True(t, regexp.MustCompile(re).MatchString(v), "expected:\n%vgot:\n%v", re, v)
}
