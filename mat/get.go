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
	var classid C.mxClassID
	var read func(unsafe.Pointer, C.size_t)
	var scalar bool

	if ivalue.Kind() == reflect.Slice {
		classid, read = readSlice(ivalue)
		scalar = false
	} else {
		classid, read = readScalar(ivalue)
		scalar = true
	}

	if classid == C.mxUNKNOWN_CLASS {
		return errors.New("unsupported type")
	}

	if classid != C.mxGetClassID(array) {
		return errors.New("data type mismatch")
	}

	count := C.mxGetM(array) * C.mxGetN(array)
	if scalar && count != 1 {
		return errors.New("data size mismatch")
	}

	parray := unsafe.Pointer(C.mxGetPr(array))
	if parray == nil {
		return errors.New("cannot read the variable")
	}

	read(parray, count)

	return nil
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

func readScalar(iv reflect.Value) (C.mxClassID, func(unsafe.Pointer, C.size_t)) {
	switch iv.Kind() {
	case reflect.Int8:
		return C.mxINT8_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetInt(int64(*(*int8)(p)))
		}
	case reflect.Uint8:
		return C.mxUINT8_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetUint(uint64(*(*uint8)(p)))
		}
	case reflect.Int16:
		return C.mxINT16_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetInt(int64(*(*int16)(p)))
		}
	case reflect.Uint16:
		return C.mxUINT16_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetUint(uint64(*(*uint16)(p)))
		}
	case reflect.Int32:
		return C.mxINT32_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetInt(int64(*(*int32)(p)))
		}
	case reflect.Uint32:
		return C.mxUINT32_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetUint(uint64(*(*uint32)(p)))
		}
	case reflect.Int64:
		return C.mxINT64_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetInt(int64(*(*int64)(p)))
		}
	case reflect.Uint64:
		return C.mxUINT64_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetUint(uint64(*(*uint64)(p)))
		}
	case reflect.Float32:
		return C.mxSINGLE_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetFloat(float64(*(*float32)(p)))
		}
	case reflect.Float64:
		return C.mxDOUBLE_CLASS, func(p unsafe.Pointer, _ C.size_t) {
			iv.SetFloat(float64(*(*float64)(p)))
		}
	default:
		return C.mxUNKNOWN_CLASS, nil
	}
}

func readSlice(iv reflect.Value) (C.mxClassID, func(unsafe.Pointer, C.size_t)) {
	read := func(w interface{}, p unsafe.Pointer, s C.size_t) {
		iw := reflect.Indirect(reflect.ValueOf(w))
		C.memcpy(unsafe.Pointer(iw.Pointer()), p, s)
		iv.Set(iw)
	}

	switch iv.Type().Elem().Kind() {
	case reflect.Int8:
		return C.mxINT8_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]int8, c)
			read(&w, p, 1*c)
		}
	case reflect.Uint8:
		return C.mxUINT8_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]uint8, c)
			read(&w, p, 1*c)
		}
	case reflect.Int16:
		return C.mxINT16_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]int16, c)
			read(&w, p, 2*c)
		}
	case reflect.Uint16:
		return C.mxUINT16_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]uint16, c)
			read(&w, p, 2*c)
		}
	case reflect.Int32:
		return C.mxINT32_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]int32, c)
			read(&w, p, 4*c)
		}
	case reflect.Uint32:
		return C.mxUINT32_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]uint32, c)
			read(&w, p, 4*c)
		}
	case reflect.Int64:
		return C.mxINT64_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]int64, c)
			read(&w, p, 8*c)
		}
	case reflect.Uint64:
		return C.mxUINT64_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]uint64, c)
			read(&w, p, 8*c)
		}
	case reflect.Float32:
		return C.mxSINGLE_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]float32, c)
			read(&w, p, 4*c)
		}
	case reflect.Float64:
		return C.mxDOUBLE_CLASS, func(p unsafe.Pointer, c C.size_t) {
			w := make([]float64, c)
			read(&w, p, 8*c)
		}
	default:
		return C.mxUNKNOWN_CLASS, nil
	}
}
