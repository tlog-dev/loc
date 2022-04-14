package loc

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type eface struct {
	_ unsafe.Pointer
	_ unsafe.Pointer
}

func TestFuncFunc(t *testing.T) {
	var f interface{}

	f = TestSetCache
	r := reflect.ValueOf(f)
	t.Logf("reflect: %v %v %v %x", r.Kind(), r.Type(), r, *(*eface)(unsafe.Pointer(&f)))

	f = TestFuncFunc
	r = reflect.ValueOf(f)
	t.Logf("reflect: %v %v %v %x", r.Kind(), r.Type(), r, *(*eface)(unsafe.Pointer(&f)))

	pc := FuncentryFromFunc(TestFuncFunc)

	name, file, line := pc.NameFileLine()
	t.Logf("pc: %v %v %v", name, file, line)

	assert.Equal(t, Funcentry(0), pc)

	assert.Equal(t, Funcentry(0), PC(reflect.ValueOf(TestFuncFunc).Pointer()))

	name, file, line = FuncentryFromFunc(nil).NameFileLine()
	t.Logf("pc: %v %v %v", name, file, line)

	var q func()

	name, file, line = FuncentryFromFunc(q).NameFileLine()
	t.Logf("pc: %v %v %v", name, file, line)

	var e PC
	q = func() {
		t.Logf("closure func")

		e = Funcentry(0)
	}

	q()

	rt := reflect.ValueOf(q).Type()
	for i := 0; i < rt.NumIn(); i++ {
		t.Logf("q in  %v", rt.In(i))
	}
	for i := 0; i < rt.NumOut(); i++ {
		t.Logf("q out %v", rt.Out(i))
	}

	pc = FuncentryFromFunc(q)

	name, file, line = pc.NameFileLine()
	t.Logf("pc: %v %v %v", name, file, line)

	assert.Equal(t, e, pc)

	assert.Panics(t, func() {
		pc = FuncentryFromFunc(3)
	})
}
