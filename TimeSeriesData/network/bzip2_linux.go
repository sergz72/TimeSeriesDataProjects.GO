package network

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lbz2
// #include "bzlib.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func bzipData(in []byte) ([]byte, error) {
	l := len(in)
	var ol C.uint = C.uint(l * 2)
	outbuf := make([]byte, l*2)
	rc := C.BZ2_bzBuffToBuffCompress((*C.char)(unsafe.Pointer(&outbuf[0])), &ol, (*C.char)(unsafe.Pointer(&in[0])),
		C.uint(l), 9, 0, 30)
	if rc != C.BZ_OK {
		return nil, fmt.Errorf("compression failure %v", rc)
	}
	return outbuf[0:ol], nil
}
