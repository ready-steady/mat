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
		return nil, errors.New("unsupported data type")
	}
}

func (f *File) createMatrix(data reflect.Value, rows, cols uint32) (*C.mxArray, error) {
	if data.Kind() != reflect.Slice {
		return nil, errors.New("expected a slice")
	}

	var classid C.mxClassID
	var size int

	switch data.Type().Elem().Kind() {
	case reflect.Float64:
		classid, size = C.mxDOUBLE_CLASS, 8
	default:
		return nil, errors.New("unsupported slice type")
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

	C.memcpy(parray, unsafe.Pointer(data.Pointer()), C.size_t(size*data.Len()))

	return array, nil
}

func (f *File) createStruct(data reflect.Value) (*C.mxArray, error) {
	if data.Kind() != reflect.Struct {
		return nil, errors.New("expected a struct")
	}

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

	array := C.mxCreateStructMatrix(1, 1, C.int(count), (**C.char)(&names[0]))
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
