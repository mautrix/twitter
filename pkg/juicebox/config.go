package juicebox

/*
#include "./juicebox-sdk-ffi.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Configuration represents a Juicebox realm configuration.
type Configuration struct {
	ptr *C.JuiceboxConfiguration
}

// ConfigurationFromJSON creates a Configuration from a JSON string.
func ConfigurationFromJSON(json string) (*Configuration, error) {
	cJSON := C.CString(json)
	defer C.free(unsafe.Pointer(cJSON))

	ptr := C.juicebox_configuration_create_from_json(cJSON)
	if ptr == nil {
		return nil, ErrAssertion
	}

	return &Configuration{ptr: ptr}, nil
}

// Destroy releases the configuration resources.
func (c *Configuration) Destroy() {
	if c.ptr != nil {
		C.juicebox_configuration_destroy(c.ptr)
		c.ptr = nil
	}
}
