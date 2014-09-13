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
