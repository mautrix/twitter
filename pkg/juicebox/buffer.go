package juicebox

/*
#include "./juicebox-sdk-ffi.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// emptyDataPtr is C-allocated memory for non-null empty array pointers.
// The Rust SDK asserts that array data pointers are non-null, even for zero-length arrays.
// We use C.malloc because Go pointers cannot be passed to C code.
var emptyDataPtr *C.uint8_t

func init() {
	// Allocate C memory once at startup (never freed, lives for program duration)
	emptyDataPtr = (*C.uint8_t)(C.malloc(1))
}

func bytesToUnmanagedArray(data []byte) C.JuiceboxUnmanagedDataArray {
	if len(data) == 0 {
		return C.JuiceboxUnmanagedDataArray{
			data:   emptyDataPtr,
			length: 0,
		}
	}
	return C.JuiceboxUnmanagedDataArray{
		data:   (*C.uint8_t)(unsafe.Pointer(&data[0])),
		length: C.size_t(len(data)),
	}
}

func unmanagedArrayToBytes(arr C.JuiceboxUnmanagedDataArray) []byte {
	if arr.length == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(arr.data), C.int(arr.length))
}

func httpMethodToString(method C.JuiceboxHttpRequestMethod) string {
	switch method {
	case C.JuiceboxHttpRequestMethodGet:
		return "GET"
	case C.JuiceboxHttpRequestMethodPut:
		return "PUT"
	case C.JuiceboxHttpRequestMethodPost:
		return "POST"
	case C.JuiceboxHttpRequestMethodDelete:
		return "DELETE"
	default:
		return "GET"
	}
}
