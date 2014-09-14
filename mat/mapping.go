package mat

// #include <mat.h>
import "C"

import (
	"reflect"
)

var classSizeMapping = map[C.mxClassID]C.size_t{
	C.mxINT8_CLASS:   1,
	C.mxUINT8_CLASS:  1,
	C.mxINT16_CLASS:  2,
	C.mxUINT16_CLASS: 2,
	C.mxINT32_CLASS:  4,
	C.mxUINT32_CLASS: 4,
	C.mxINT64_CLASS:  8,
	C.mxUINT64_CLASS: 8,
	C.mxSINGLE_CLASS: 4,
	C.mxDOUBLE_CLASS: 8,
}

var kindClassMapping = map[reflect.Kind]C.mxClassID{
	reflect.Int8:    C.mxINT8_CLASS,
	reflect.Uint8:   C.mxUINT8_CLASS,
	reflect.Int16:   C.mxINT16_CLASS,
	reflect.Uint16:  C.mxUINT16_CLASS,
	reflect.Int32:   C.mxINT32_CLASS,
	reflect.Uint32:  C.mxUINT32_CLASS,
	reflect.Int64:   C.mxINT64_CLASS,
	reflect.Uint64:  C.mxUINT64_CLASS,
	reflect.Float32: C.mxSINGLE_CLASS,
	reflect.Float64: C.mxDOUBLE_CLASS,
}
