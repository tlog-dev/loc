package loc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationFillCallers(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	var st PCs
	func() {
		func() {
			st = testLocationsInside()
		}()
	}()

	assert.Len(t, st, 3)
	assert.Equal(t, "location_stack_test.go:26", st[0].String())
	assert.Equal(t, "location_stack_test.go:27", st[1].String())
	assert.Equal(t, "location_stack_test.go:38", st[2].String())

	re := `location_stack_test.go:26 at location_stack_test.go:27 at location_stack_test.go:38`

	assert.Equal(t, re, st.String())
}

func TestLocationPCsFormat(t *testing.T) {
	t.Parallel()

	var st PCs
	func() {
		func() {
			st = testLocationsInside()
		}()
	}()

	assert.Equal(t, "location_stack_test.go:26 at location_stack_test.go:27 at location_stack_test.go:58", fmt.Sprintf("%v", st))

	assert.Equal(t, "loc.testLocationsInside.func1:26 at loc.testLocationsInside:27 at loc.TestLocationPCsFormat.func1:58", fmt.Sprintf("%#v", st))

	re := `at [\w.-/]*location_stack_test.go:26
at [\w.-/]*location_stack_test.go:27
at [\w.-/]*location_stack_test.go:58
`
	v := fmt.Sprintf("%+v", st)
	assert.True(t, regexp.MustCompile(re).MatchString(v), "expected:\n%vgot:\n%v", re, v)
}

func TestLocationPCsFormatString(t *testing.T) {
	t.Parallel()

	var st PCs
	func() {
		func() {
			st = testLocationsInside()
		}()
	}()

	assert.Equal(t, "location_stack_test.go:26 at location_stack_test.go:27 at location_stack_test.go:80", st.FormatString(""))

	assert.Equal(t, "loc.testLocationsInside.func1:26 at loc.testLocationsInside:27 at loc.TestLocationPCsFormatString.func1:80", st.FormatString("#"))

	re := `at [\w.-/]*location_stack_test.go:26
at [\w.-/]*location_stack_test.go:27
at [\w.-/]*location_stack_test.go:80
`

	v := st.FormatString("+")
	assert.True(t, regexp.MustCompile(re).MatchString(v), "expected:\n%vgot:\n%v", re, v)
}
