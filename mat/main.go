// Package mat is an adapter for the MATLAB MAT-file API.
//
// http://www.mathworks.com/help/pdf_doc/matlab/apiext.pdf
package mat

// #cgo LDFLAGS: -lmat -lmx
// #include <string.h>
// #include <mat.h>
import "C"

import (
	"errors"
	"unsafe"
)

type File struct {
	path string
	mat *C.MATFile
}

func Open(path string, mode string) (*File, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cmode := C.CString(mode)
	defer C.free(unsafe.Pointer(cmode))

	file := &File{ path: path, mat: C.matOpen(cpath, cmode) }
	if file.mat == nil {
		goto error
	}

	return file, nil

error:
	return nil, errors.New("Cannot open the file.")
}

func (f *File) Close() {
	C.matClose(f.mat)
	f.mat = nil
}

func (f *File) PutMatrix(name string, rows, cols uint32,
	data []float64) error {

	var cname *C.char
	var pmatrix unsafe.Pointer

	matrix := C.mxCreateDoubleMatrix(C.size_t(rows), C.size_t(cols), C.mxREAL)
	if matrix == nil {
		goto error
	}

	defer C.mxDestroyArray(matrix)

	pmatrix = unsafe.Pointer(C.mxGetPr(matrix))
	if pmatrix == nil {
		goto error
	}

	C.memcpy(pmatrix, unsafe.Pointer(&data[0]), C.size_t(8 * len(data)))

	cname = C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if C.matPutVariable(f.mat, cname, matrix) != 0 {
		goto error
	}

	return nil

error:
	return errors.New("Cannot write into the file.")
}
