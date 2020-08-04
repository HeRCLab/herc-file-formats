package main

//#include <stdlib.h>
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/herclab/wavegen/pkg/wavegen"
)

var nexthandle C.int = 0
var handles map[C.int]*wavegen.WaveFile = map[C.int]*wavegen.WaveFile{}

var lastError string = ""

//export WGGetError
func WGGetError() *C.char {
	return C.CString(lastError)
}

//export WGOpen
func WGOpen(path *C.char, handle *C.int) C.int {
	h := nexthandle
	nexthandle++

	wf, err := wavegen.ReadJSON(C.GoString(path))
	if err != nil {
		lastError = fmt.Sprintf("%v", err)
		return 1
	}

	handles[h] = wf
	*handle = h
	return 0
}

//export WGClose
func WGClose(handle C.int) C.int {
	if _, ok := handles[handle]; ok {
		delete(handles, handle)
		return 0
	} else {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return 1
	}
}

//export WGReadS
func WGReadS(handle C.int, index C.int, result *C.double) C.int {
	if _, ok := handles[handle]; ok {
		err := handles[handle].Signal.ValidateIndex(int(index))

		if err != nil {
			lastError = fmt.Sprintf("%v", err)
			return 1
		}

		*result = C.double(handles[handle].Signal.S[int(index)])

		return 0
	} else {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return 1
	}
}

//export WGReadT
func WGReadT(handle C.int, index C.int, result *C.double) C.int {
	if _, ok := handles[handle]; ok {
		err := handles[handle].Signal.ValidateIndex(int(index))

		if err != nil {
			lastError = fmt.Sprintf("%v", err)
			return 1
		}

		*result = C.double(handles[handle].Signal.T[int(index)])

		return 0
	} else {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return 1
	}
}

//export WGSize
func WGSize(handle C.int, result *C.int) C.int {
	if _, ok := handles[handle]; ok {
		*result = C.int(handles[handle].Signal.Size())
		return 0
	} else {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return 1
	}
}

//export WGCopyS
func WGCopyS(handle C.int, buf *C.double) C.int {
	if _, ok := handles[handle]; !ok {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return 1
	}

	var slice []C.double
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = handles[handle].Signal.Size()
	header.Len = handles[handle].Signal.Size()
	header.Data = uintptr(unsafe.Pointer(buf))

	for i, v := range handles[handle].Signal.S {
		slice[i] = C.double(v)
	}

	return 0
}

//export WGCopyT
func WGCopyT(handle C.int, buf *C.double) C.int {
	if _, ok := handles[handle]; !ok {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return 1
	}

	var slice []C.double
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = handles[handle].Signal.Size()
	header.Len = handles[handle].Signal.Size()
	header.Data = uintptr(unsafe.Pointer(buf))

	for i, v := range handles[handle].Signal.T {
		slice[i] = C.double(v)
	}

	return 0
}

//export WGSampleRate
func WGSampleRate(handle C.int, result *C.double) C.int {
	if _, ok := handles[handle]; !ok {
		lastError = fmt.Sprintf("unknown handle %d", int(handle))
		return 1
	}

	*result = C.double(handles[handle].Signal.SampleRate)

	return 0
}

func main() {
}
