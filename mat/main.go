package mat

// #cgo LDFLAGS: -lmat
// #include <mat.h>
import "C"

import "unsafe"

type File struct {
	path string
}

func Open(path string) (*File, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cmode := C.CString("r")
	defer C.free(unsafe.Pointer(cmode))

	file := C.matOpen(cpath, cmode)
	C.matClose(file)

	return &File{ path: path }, nil
}
