package loc

import (
	"reflect"
	"unsafe"
)

type (
	// rtype is the common implementation of most values.
	// It is embedded in other struct types.
	//
	// rtype must be kept in sync with ../runtime/type.go:/^type._type.
	//nolint:structcheck,unused
	rtype struct {
		size       uintptr
		ptrdata    uintptr // number of bytes in the type that can contain pointers
		hash       uint32  // hash of type; avoids computation in hash tables
		tflag      uint8   // extra type information flags
		align      uint8   // alignment of variable with this type
		fieldAlign uint8   // alignment of struct field with this type
		kind       uint8   // enumeration for C
		// function for comparing objects of this type
		// (ptr to object A, ptr to object B) -> ==?
		equal  func(unsafe.Pointer, unsafe.Pointer) bool
		gcdata *byte // garbage collection data
		//	str       nameOff // string form
		//	ptrToThis typeOff // type for pointer to this type, may be zero
	}

	fface struct {
		t *rtype
		e *uintptr
	}
)

//nolint:deadcode,varcheck
const (
	kindDirectIface = 1 << 5
	kindGCProg      = 1 << 6 // Type.gc points to GC program
	kindMask        = (1 << 5) - 1
)

func FuncEntryFromFunc(f interface{}) PC {
	ff := (*fface)(unsafe.Pointer(&f))

	if ff.t == nil {
		return 0
	}

	if reflect.Kind(ff.t.kind&kindMask) != reflect.Func {
		panic("not a function")
	}

	if ff.e == nil {
		return 0
	}

	return PC(*ff.e)
}
