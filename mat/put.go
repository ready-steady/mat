package mat

// #include <string.h>
// #include <mat.h>
import "C"

import (
	"errors"
	"reflect"
	"unsafe"
)

// Put writes an object into the file.
func (f *File) Put(name string, object interface{}) error {
	array, err := f.writeArray(reflect.ValueOf(object))
	if err != nil {
		return err
	}
	defer C.mxDestroyArray(array)

	return f.putVariable(name, array)
}

// Put writes a matrix into the file.
func (f *File) PutMatrix(name string, object interface{}, rows, cols uint32) error {
	value := reflect.ValueOf(object)
	if uint32(value.Len()) != rows*cols {
		return errors.New("dimension mismatch")
	}

	array, err := f.writeMatrix(value, rows, cols)
	if err != nil {
		return err
	}
	defer C.mxDestroyArray(array)

	return f.putVariable(name, array)
}

func (f *File) writeArray(value reflect.Value) (*C.mxArray, error) {
	switch value.Kind() {
	case reflect.Slice:
		return f.writeMatrix(value, 1, uint32(value.Len()))
	case reflect.Struct:
		return f.writeStruct(value)
	default:
		return f.writeMatrix(value, 1, 1)
	}
}

func (f *File) writeMatrix(value reflect.Value, rows, cols uint32) (*C.mxArray, error) {
	var classid C.mxClassID
	var write func(unsafe.Pointer)

	if value.Kind() == reflect.Slice {
		classid, write = writeSlice(value)
	} else {
		classid, write = writeScalar(value)
	}

	if classid == C.mxUNKNOWN_CLASS {
		errors.New("unsupported type")
	}

	array := C.mxCreateNumericMatrix(C.size_t(rows), C.size_t(cols), classid, C.mxREAL)
	if array == nil {
		return nil, errors.New("cannot create a matrix")
	}

	parray := unsafe.Pointer(C.mxGetPr(array))
	if parray == nil {
		C.mxDestroyArray(array)
		return nil, errors.New("cannot create a matrix")
	}

	write(parray)

	return array, nil
}

func (f *File) writeStruct(value reflect.Value) (*C.mxArray, error) {
	typo := value.Type()
	count := typo.NumField()
	names := make([]*C.char, 0, count)
	arrays := make([]*C.mxArray, 0, count)

	// NOTE: Should be called only when the function fails. If it succeeds, all
	// arrays will be disposed when the struct gets destroyed.
	cleanup := func() {
		for _, array := range arrays {
			C.mxDestroyArray(array)
		}
	}

	for i := 0; i < count; i++ {
		field := typo.Field(i)

		name := C.CString(field.Name)
		defer C.free(unsafe.Pointer(name))
		names = append(names, name)

		array, err := f.writeArray(value.Field(i))
		if err != nil {
			cleanup()
			return nil, err
		}
		arrays = append(arrays, array)
	}
	count = len(names)

	var pnames **C.char
	if count > 0 {
		pnames = (**C.char)(&names[0])
	}

	array := C.mxCreateStructMatrix(1, 1, C.int(count), pnames)
	if array == nil {
		cleanup()
		return nil, errors.New("cannot create a struct")
	}

	for i := 0; i < count; i++ {
		C.mxSetFieldByNumber(array, 0, C.int(i), arrays[i])
	}

	return array, nil
}

func writeScalar(v reflect.Value) (C.mxClassID, func(unsafe.Pointer)) {
	switch v.Kind() {
	case reflect.Int8:
		return C.mxINT8_CLASS, func(p unsafe.Pointer) {
			*((*int8)(p)) = int8(v.Int())
		}
	case reflect.Uint8:
		return C.mxUINT8_CLASS, func(p unsafe.Pointer) {
			*((*uint8)(p)) = uint8(v.Uint())
		}
	case reflect.Int16:
		return C.mxINT16_CLASS, func(p unsafe.Pointer) {
			*((*int16)(p)) = int16(v.Int())
		}
	case reflect.Uint16:
		return C.mxUINT16_CLASS, func(p unsafe.Pointer) {
			*((*uint16)(p)) = uint16(v.Uint())
		}
	case reflect.Int32:
		return C.mxINT32_CLASS, func(p unsafe.Pointer) {
			*((*int32)(p)) = int32(v.Int())
		}
	case reflect.Uint32:
		return C.mxUINT32_CLASS, func(p unsafe.Pointer) {
			*((*uint32)(p)) = uint32(v.Uint())
		}
	case reflect.Int64:
		return C.mxINT64_CLASS, func(p unsafe.Pointer) {
			*((*int64)(p)) = int64(v.Int())
		}
	case reflect.Uint64:
		return C.mxUINT64_CLASS, func(p unsafe.Pointer) {
			*((*uint64)(p)) = uint64(v.Uint())
		}
	case reflect.Float32:
		return C.mxSINGLE_CLASS, func(p unsafe.Pointer) {
			*((*float32)(p)) = float32(v.Float())
		}
	case reflect.Float64:
		return C.mxDOUBLE_CLASS, func(p unsafe.Pointer) {
			*((*float64)(p)) = float64(v.Float())
		}
	default:
		return C.mxUNKNOWN_CLASS, nil
	}
}

func writeSlice(v reflect.Value) (C.mxClassID, func(unsafe.Pointer)) {
	var c C.mxClassID
	var s C.size_t

	switch v.Type().Elem().Kind() {
	case reflect.Int8:
		c, s = C.mxINT8_CLASS, 1
	case reflect.Uint8:
		c, s = C.mxUINT8_CLASS, 1
	case reflect.Int16:
		c, s = C.mxINT16_CLASS, 2
	case reflect.Uint16:
		c, s = C.mxUINT16_CLASS, 2
	case reflect.Int32:
		c, s = C.mxINT32_CLASS, 4
	case reflect.Uint32:
		c, s = C.mxUINT32_CLASS, 4
	case reflect.Int64:
		c, s = C.mxINT64_CLASS, 8
	case reflect.Uint64:
		c, s = C.mxUINT64_CLASS, 8
	case reflect.Float32:
		c, s = C.mxSINGLE_CLASS, 4
	case reflect.Float64:
		c, s = C.mxDOUBLE_CLASS, 8
	default:
		return C.mxUNKNOWN_CLASS, nil
	}

	return c, func(p unsafe.Pointer) {
		C.memcpy(p, unsafe.Pointer(v.Pointer()), C.size_t(v.Len())*s)
	}
}

func (f *File) putVariable(name string, array *C.mxArray) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if C.matPutVariable(f.mat, cname, array) != 0 {
		return errors.New("cannot write a variable into the file")
	}

	return nil
}
