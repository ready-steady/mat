// Package mat provides an adapter to the MATLAB MAT-file API.
//
// http://www.mathworks.com/help/pdf_doc/matlab/apiext.pdf
package mat

// #cgo LDFLAGS: -lmat -lmx
// #include <string.h>
// #include <mat.h>
import "C"

import (
	"errors"
	"reflect"
	"unsafe"
)

// File represents a MAT file.
type File struct {
	mat *C.MATFile
}

// Open opens a MAT file for reading and writing.
//
// http://www.mathworks.se/help/matlab/apiref/matopen.html
func Open(path string, mode string) (*File, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cmode := C.CString(mode)
	defer C.free(unsafe.Pointer(cmode))

	file := &File{mat: C.matOpen(cpath, cmode)}
	if file.mat == nil {
		return nil, errors.New("cannot open the file")
	}

	return file, nil
}

// Close closes the MAT file.
func (f *File) Close() {
	C.matClose(f.mat)
	f.mat = nil
}

// Put writes an array into the MAT file.
func (f *File) Put(name string, data interface{}) error {
	array, err := f.createArray(reflect.ValueOf(data))
	if err != nil {
		return err
	}
	defer C.mxDestroyArray(array)

	return f.putVariable(name, array)
}

// Put writes a matrix into the MAT file.
func (f *File) PutMatrix(name string, data interface{}, rows, cols uint32) error {
	array, err := f.createMatrix(reflect.ValueOf(data), rows, cols)
	if err != nil {
		return err
	}
	defer C.mxDestroyArray(array)

	return f.putVariable(name, array)
}

func (f *File) createArray(data reflect.Value) (*C.mxArray, error) {
	switch data.Kind() {
	case reflect.Slice:
		return f.createMatrix(data, 1, uint32(data.Len()))
	case reflect.Struct:
		return f.createStruct(data)
	default:
		return f.createScalar(data)
	}
}

func (f *File) createScalar(data reflect.Value) (*C.mxArray, error) {
	classid, size, write := mapToMATLAB(data.Kind())
	if size == 0 {
		return nil, errors.New("unsupported scalar type")
	}

	array := C.mxCreateNumericMatrix(1, 1, classid, C.mxREAL)
	if array == nil {
		return nil, errors.New("cannot create a scalar")
	}

	parray := unsafe.Pointer(C.mxGetPr(array))
	if parray == nil {
		C.mxDestroyArray(array)
		return nil, errors.New("cannot create a scalar")
	}

	write(parray, data)

	return array, nil
}

func (f *File) createMatrix(data reflect.Value, rows, cols uint32) (*C.mxArray, error) {
	classid, size, _ := mapToMATLAB(data.Type().Elem().Kind())
	if size == 0 {
		errors.New("unsupported slice type")
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

	C.memcpy(parray, unsafe.Pointer(data.Pointer()), C.size_t(data.Len())*size)

	return array, nil
}

func (f *File) createStruct(data reflect.Value) (*C.mxArray, error) {
	typo := data.Type()
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
		if field.PkgPath != "" { // not exported
			continue
		}

		name := C.CString(field.Name)
		defer C.free(unsafe.Pointer(name))
		names = append(names, name)

		array, err := f.createArray(data.Field(i))
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

func (f *File) putVariable(name string, array *C.mxArray) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if C.matPutVariable(f.mat, cname, array) != 0 {
		return errors.New("cannot write a variable into the file")
	}

	return nil
}

func mapToMATLAB(kind reflect.Kind) (C.mxClassID, C.size_t, func(unsafe.Pointer, reflect.Value)) {
	switch kind {
	case reflect.Int8:
		return C.mxINT8_CLASS, 1, func(p unsafe.Pointer, v reflect.Value) {
			*((* int8)(p)) = int8(v.Int())
		}
	case reflect.Uint8:
		return C.mxUINT8_CLASS, 1, func(p unsafe.Pointer, v reflect.Value) {
			*((* uint8)(p)) = uint8(v.Uint())
		}
	case reflect.Int16:
		return C.mxINT16_CLASS, 2, func(p unsafe.Pointer, v reflect.Value) {
			*((* int16)(p)) = int16(v.Int())
		}
	case reflect.Uint16:
		return C.mxUINT16_CLASS, 2, func(p unsafe.Pointer, v reflect.Value) {
			*((* uint16)(p)) = uint16(v.Uint())
		}
	case reflect.Int32:
		return C.mxINT32_CLASS, 4, func(p unsafe.Pointer, v reflect.Value) {
			*((* int32)(p)) = int32(v.Int())
		}
	case reflect.Uint32:
		return C.mxUINT32_CLASS, 4, func(p unsafe.Pointer, v reflect.Value) {
			*((* uint32)(p)) = uint32(v.Uint())
		}
	case reflect.Int64:
		return C.mxINT64_CLASS, 8, func(p unsafe.Pointer, v reflect.Value) {
			*((* int64)(p)) = int64(v.Int())
		}
	case reflect.Uint64:
		return C.mxUINT64_CLASS, 8, func(p unsafe.Pointer, v reflect.Value) {
			*((* uint64)(p)) = uint64(v.Uint())
		}
	case reflect.Float32:
		return C.mxSINGLE_CLASS, 4, func(p unsafe.Pointer, v reflect.Value) {
			*((* float32)(p)) = float32(v.Float())
		}
	case reflect.Float64:
		return C.mxDOUBLE_CLASS, 8, func(p unsafe.Pointer, v reflect.Value) {
			*((* float64)(p)) = float64(v.Float())
		}
	default:
		return 0, 0, nil
	}
}
