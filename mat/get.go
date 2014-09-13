package mat

// #include <string.h>
// #include <mat.h>
import "C"

import (
	"errors"
	"reflect"
	"unsafe"
)

func (f *File) Get(name string, object interface{}) error {
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr {
		return errors.New("expected a pointer")
	}
	value = reflect.Indirect(value)

	classid, _ := mapFromMATLAB(value.Type())
	if classid == C.mxUNKNOWN_CLASS {
		return errors.New("unsupported data type")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	array := C.matGetVariable(f.mat, cname)
	if array == nil {
		return errors.New("cannot get the variable")
	}
	defer C.mxDestroyArray(array)

	if classid != C.mxGetClassID(array) {
		return errors.New("data type missmatch")
	}

	return nil
}

func mapFromMATLAB(typo reflect.Type) (C.mxClassID, C.size_t) {
	switch typo.Kind() {
	case reflect.Struct:
		return C.mxSTRUCT_CLASS, 0
	case reflect.Slice:
		typo = typo.Elem()
	}

	switch typo.Kind() {
	case reflect.Int8:
		return C.mxINT8_CLASS, 1
	case reflect.Uint8:
		return C.mxUINT8_CLASS, 1
	case reflect.Int16:
		return C.mxINT16_CLASS, 2
	case reflect.Uint16:
		return C.mxUINT16_CLASS, 2
	case reflect.Int32:
		return C.mxINT32_CLASS, 4
	case reflect.Uint32:
		return C.mxUINT32_CLASS, 4
	case reflect.Int64:
		return C.mxINT64_CLASS, 8
	case reflect.Uint64:
		return C.mxUINT64_CLASS, 8
	case reflect.Float32:
		return C.mxSINGLE_CLASS, 4
	case reflect.Float64:
		return C.mxDOUBLE_CLASS, 8
	default:
		return C.mxUNKNOWN_CLASS, 0
	}
}
