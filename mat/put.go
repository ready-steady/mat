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
	array, err := f.writeObject(reflect.ValueOf(object))
	if err != nil {
		return err
	}
	defer C.mxDestroyArray(array)

	return f.putVariable(name, array)
}

// PutArray writes a multidimensional array into the file.
func (f *File) PutArray(name string, object interface{}, dimensions ...uint) error {
	value := reflect.ValueOf(object)
	length := uint(value.Len())

	switch len(dimensions) {
	case 0:
		dimensions = append(dimensions, length, 1)
	case 1:
		dimensions = append(dimensions, 1)
	}

	size := dimensions[0]
	for i := 1; i < len(dimensions); i++ {
		size *= dimensions[i]
	}

	if length != size {
		return errors.New("dimension mismatch")
	}

	array, err := f.writeArray(value, dimensions...)
	if err != nil {
		return err
	}
	defer C.mxDestroyArray(array)

	return f.putVariable(name, array)
}

// PutMatrix writes a matrix into the file.
func (f *File) PutMatrix(name string, object interface{}, rows, cols uint) error {
	return f.PutArray(name, object, rows, cols)
}

func (f *File) writeObject(value reflect.Value) (*C.mxArray, error) {
	switch value.Kind() {
	case reflect.Slice:
		return f.writeArray(value, uint(value.Len()), 1)
	case reflect.Struct:
		return f.writeStruct(value)
	default:
		return f.writeArray(value, 1, 1)
	}
}

func (f *File) writeArray(value reflect.Value, dimensions ...uint) (*C.mxArray, error) {
	var kind reflect.Kind

	if value.Kind() == reflect.Slice {
		kind = value.Type().Elem().Kind()
	} else {
		kind = value.Kind()
	}

	classid, ok := kindClassMapping[kind]
	if !ok {
		return nil, errors.New("unsupported data type")
	}

	// NOTE: Do we need a proper conversion from uint to C.size_t?
	array := C.mxCreateNumericArray(C.size_t(len(dimensions)),
		(*C.size_t)(unsafe.Pointer(&dimensions[0])), classid, C.mxREAL)
	if array == nil {
		return nil, errors.New("cannot create an array")
	}

	parray := unsafe.Pointer(C.mxGetPr(array))
	if parray == nil {
		C.mxDestroyArray(array)
		return nil, errors.New("cannot create an array")
	}

	size, ok := classSizeMapping[classid]
	if !ok {
		return nil, errors.New("unsupported data type")
	}

	if value.Kind() == reflect.Slice {
		if err := writeSlice(value, parray, size); err != nil {
			C.mxDestroyArray(array)
			return nil, err
		}
	} else {
		if err := writeScalar(value, parray); err != nil {
			C.mxDestroyArray(array)
			return nil, err
		}
	}

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

		array, err := f.writeObject(value.Field(i))
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

func writeScalar(v reflect.Value, p unsafe.Pointer) error {
	switch v.Kind() {
	case reflect.Int8:
		*((*int8)(p)) = int8(v.Int())
	case reflect.Uint8:
		*((*uint8)(p)) = uint8(v.Uint())
	case reflect.Int16:
		*((*int16)(p)) = int16(v.Int())
	case reflect.Uint16:
		*((*uint16)(p)) = uint16(v.Uint())
	case reflect.Int32:
		*((*int32)(p)) = int32(v.Int())
	case reflect.Uint32:
		*((*uint32)(p)) = uint32(v.Uint())
	case reflect.Int64:
		*((*int64)(p)) = int64(v.Int())
	case reflect.Uint64:
		*((*uint64)(p)) = uint64(v.Uint())
	case reflect.Float32:
		*((*float32)(p)) = float32(v.Float())
	case reflect.Float64:
		*((*float64)(p)) = float64(v.Float())
	default:
		return errors.New("unsupported data type")
	}

	return nil
}

func writeSlice(v reflect.Value, p unsafe.Pointer, s C.size_t) error {
	C.memcpy(p, unsafe.Pointer(v.Pointer()), C.size_t(v.Len())*s)

	return nil
}

func (f *File) putVariable(name string, array *C.mxArray) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if C.matPutVariable(f.mat, cname, array) != 0 {
		return errors.New("cannot write a variable into the file")
	}

	return nil
}
