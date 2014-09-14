package mat

// #include <string.h>
// #include <mat.h>
import "C"

import (
	"errors"
	"reflect"
	"unsafe"
)

// Get reads an object from the file.
func (f *File) Get(name string, object interface{}) error {
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr {
		return errors.New("expected a pointer")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	array := C.matGetVariable(f.mat, cname)
	if array == nil {
		return errors.New("cannot find the variable")
	}
	defer C.mxDestroyArray(array)

	return f.readArray(array, value)
}

func (f *File) readArray(array *C.mxArray, value reflect.Value) error {
	ivalue := reflect.Indirect(value)
	switch ivalue.Kind() {
	case reflect.Struct:
		return f.readStruct(array, ivalue)
	default:
		return f.readMatrix(array, ivalue)
	}
}

func (f *File) readMatrix(array *C.mxArray, ivalue reflect.Value) error {
	parray := unsafe.Pointer(C.mxGetPr(array))
	if parray == nil {
		return errors.New("cannot read the variable")
	}

	count := C.mxGetM(array) * C.mxGetN(array)
	size, ok := classSizeMapping[C.mxGetClassID(array)]
	if !ok {
		return errors.New("unsupported data type")
	}

	if ivalue.Kind() == reflect.Slice {
		return readSlice(ivalue, parray, count, size)
	} else {
		if count != 1 {
			return errors.New("data size mismatch")
		}
		return readScalar(ivalue, parray)
	}
}

func (f *File) readStruct(array *C.mxArray, ivalue reflect.Value) error {
	if C.mxSTRUCT_CLASS != C.mxGetClassID(array) {
		return errors.New("data type mismatch")
	}

	if C.mxGetM(array)*C.mxGetN(array) != 1 {
		return errors.New("data size mismatch")
	}

	typo := ivalue.Type()
	count := typo.NumField()

	if count != int(C.mxGetNumberOfFields(array)) {
		return errors.New("data structure mismatch")
	}

	for i := 0; i < count; i++ {
		field := typo.Field(i)

		name := C.CString(field.Name)
		defer C.free(unsafe.Pointer(name))

		farray := C.mxGetField(array, 0, name)
		if farray == nil {
			return errors.New("data structure mismatch")
		}

		if err := f.readArray(farray, ivalue.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func readScalar(iv reflect.Value, p unsafe.Pointer) error {
	switch iv.Kind() {
	case reflect.Int8:
		*(*int8)(unsafe.Pointer(iv.UnsafeAddr())) = *(*int8)(p)
	case reflect.Uint8:
		*(*uint8)(unsafe.Pointer(iv.UnsafeAddr())) = *(*uint8)(p)
	case reflect.Int16:
		*(*int16)(unsafe.Pointer(iv.UnsafeAddr())) = *(*int16)(p)
	case reflect.Uint16:
		*(*uint16)(unsafe.Pointer(iv.UnsafeAddr())) = *(*uint16)(p)
	case reflect.Int32:
		*(*int32)(unsafe.Pointer(iv.UnsafeAddr())) = *(*int32)(p)
	case reflect.Uint32:
		*(*uint32)(unsafe.Pointer(iv.UnsafeAddr())) = *(*uint32)(p)
	case reflect.Int64:
		*(*int64)(unsafe.Pointer(iv.UnsafeAddr())) = *(*int64)(p)
	case reflect.Uint64:
		*(*uint64)(unsafe.Pointer(iv.UnsafeAddr())) = *(*uint64)(p)
	case reflect.Float32:
		*(*float32)(unsafe.Pointer(iv.UnsafeAddr())) = *(*float32)(p)
	case reflect.Float64:
		*(*float64)(unsafe.Pointer(iv.UnsafeAddr())) = *(*float64)(p)
	default:
		return errors.New("unsupported data type")
	}

	return nil
}

func readSlice(iv reflect.Value, p unsafe.Pointer, c C.size_t, s C.size_t) error {
	w := reflect.MakeSlice(iv.Type(), int(c), int(c))
	C.memcpy(unsafe.Pointer(w.Pointer()), p, c*s)

	iw := reflect.Indirect(reflect.New(iv.Type()))
	iw.Set(w)

	// FIXME: Bad, bad, bad! But how else to fill in unexported fields?
	src := (*reflect.SliceHeader)(unsafe.Pointer(iw.UnsafeAddr()))
	dst := (*reflect.SliceHeader)(unsafe.Pointer(iv.UnsafeAddr()))

	dst.Data, src.Data = src.Data, dst.Data
	dst.Len, src.Len = src.Len, dst.Len
	dst.Cap, src.Cap = src.Cap, dst.Cap

	return nil
}
