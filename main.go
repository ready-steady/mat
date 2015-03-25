// Package mat provides a reader and writer of MATLAB MAT-files.
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

const (
	is64bit = uint64(^uint(0)) == ^uint64(0)
)

// File represents a file.
type File struct {
	mat *C.MATFile
}

// Open opens a file for reading and writing.
//
// http://www.mathworks.com/help/matlab/apiref/matopen.html
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

// Close closes the file.
func (f *File) Close() {
	C.matClose(f.mat)
	f.mat = nil
}
