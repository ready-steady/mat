package mat

// #include <string.h>
// #include <mat.h>
import "C"

import (
	"errors"
	"reflect"
	"unsafe"
)

func (f *File) Get(name string, object interface{}) error {
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr {
		return errors.New("expected a pointer")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return f.getArray(cname, value)
}

func (f *File) getArray(name *C.char, value reflect.Value) error {
	ivalue := reflect.Indirect(value)
	switch ivalue.Kind() {
	case reflect.Slice:
		return f.copyArray(name, ivalue, ivalue.Type().Elem().Kind(), false)
	case reflect.Struct:
		return f.getStruct(name, ivalue)
	default:
		return f.copyArray(name, ivalue, ivalue.Kind(), true)
	}
}

func (f *File) getStruct(name *C.char, ivalue reflect.Value) error {
	return nil
}

func (f *File) copyArray(name *C.char, ivalue reflect.Value, kind reflect.Kind, scalar bool) error {
	classid, writeScalar, writeSlice := mapFromMATLAB(kind)
	if classid == C.mxUNKNOWN_CLASS {
		return errors.New("unsupported type")
	}

	array := C.matGetVariable(f.mat, name)
	if array == nil {
		return errors.New("cannot find the variable")
	}
	defer C.mxDestroyArray(array)

	if classid != C.mxGetClassID(array) {
		return errors.New("data type missmatch")
	}

	count := C.mxGetM(array)*C.mxGetN(array)
	if scalar && count != 1 {
		return errors.New("data size missmatch")
	}

	parray := unsafe.Pointer(C.mxGetPr(array))
	if parray == nil {
		return errors.New("cannot read the variable")
	}

	if scalar {
		writeScalar(ivalue, parray)
	} else {
		writeSlice(ivalue, parray, count)
	}

	return nil
}

func mapFromMATLAB(kind reflect.Kind) (C.mxClassID, func(reflect.Value, unsafe.Pointer), func(reflect.Value, unsafe.Pointer, C.size_t)) {
	switch kind {
	case reflect.Int8:
		return C.mxINT8_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetInt(int64(*(*int8)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]int8, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 1*c)
				}
			}
	case reflect.Uint8:
		return C.mxUINT8_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetUint(uint64(*(*uint8)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]uint8, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 1*c)
				}
			}
	case reflect.Int16:
		return C.mxINT16_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetInt(int64(*(*int16)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]int16, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 2*c)
				}
			}
	case reflect.Uint16:
		return C.mxUINT16_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetUint(uint64(*(*uint16)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]uint16, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 2*c)
				}
			}
	case reflect.Int32:
		return C.mxINT32_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetInt(int64(*(*int32)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]int32, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 4*c)
				}
			}
	case reflect.Uint32:
		return C.mxUINT32_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetUint(uint64(*(*uint32)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]uint32, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 4*c)
				}
			}
	case reflect.Int64:
		return C.mxINT64_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetInt(int64(*(*int64)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]int64, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 8*c)
				}
			}
	case reflect.Uint64:
		return C.mxUINT64_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetUint(uint64(*(*uint64)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]uint64, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 8*c)
				}
			}
	case reflect.Float32:
		return C.mxSINGLE_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetFloat(float64(*(*float32)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]float32, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 4*c)
				}
			}
	case reflect.Float64:
		return C.mxDOUBLE_CLASS,
			func(iv reflect.Value, p unsafe.Pointer) {
				iv.SetFloat(float64(*(*float64)(p)))
			},
			func(iv reflect.Value, p unsafe.Pointer, c C.size_t) {
				w := make([]float64, c)
				iv.Set(reflect.Indirect(reflect.ValueOf(&w)))
				if c > 0 {
					C.memcpy(unsafe.Pointer(&w[0]), p, 8*c)
				}
			}
	case reflect.Struct:
		return C.mxSTRUCT_CLASS, nil, nil
	default:
		return C.mxUNKNOWN_CLASS, nil, nil
	}
}
