wavegen(3) -- WaveGen C API
===========================


Wavegen's C API is a Cgo shared library. It is recommended to simply copy
[`client.go`](./client.go) into your project, and add appropriate Make rules to
generate `client.so` and `client.h`.

An example of using the client library is provided in
[`example.c`](./example.c).

## Running the example

```
$ wavegen generate -s 100 -d 0.2 -o out.json
$ go build -o client.so -buildmode=c-shared client.go
$ cc example.c client.so
$ ./a.out out.json
Opened handle 0 on file 'out.json'
Handle contains 20 samples
Retrieved sample rate 100.000000
First 10 samples... S[0]=0.000000 T[0]=0.000000
S[1]=0.062791 T[1]=0.010000
S[2]=0.125333 T[2]=0.020000
S[3]=0.187381 T[3]=0.030000
S[4]=0.248690 T[4]=0.040000
S[5]=0.309017 T[5]=0.050000
S[6]=0.368125 T[6]=0.060000
S[7]=0.425779 T[7]=0.070000
S[8]=0.481754 T[8]=0.080000
S[9]=0.535827 T[9]=0.090000
```

## Library Documentation

### Conventions

All client functions return a C integer, which will be `0` on success, and `1`
on failure. If you wish to learn more about the error which occurred, the
function `WGGetError()` will return the string describing the most recently
encountered error message. As a result, any functions which return results do
so by reference.

The library is organized around the concept of a "handle". The object
representing the file is stored in Go memory, and not exposed to the C client
at all. Instead, the library returns an integer handle (a reference into an
internal hash table) that is used to refer to an entry in said table.

### `char* WGGetError(void)`

Returns the most recently encountered error text. The string returned in
allocated in C memory, and should be `free()`-ed by the caller.

### `int WGOpen(char* path, int* handle)`

Opens a new WaveGen JSON file, writing the handle into the `handle` parameter.

### `int WGClose(int handle)`

Close a previously opened WaveGen file handle.

### `int WGReadS(int handle, int index, double* result)`

Reads the signal value at the specified index.

### `int WGReadT(int handle, int index, double* result)`

Reads the time value at the specified.

### `int WGSize(int handle, int* result)`

Retrieves the number of samples in the file corresponding to the handle.

### `int WGCopyS(int handle, double* buf)`

Copies all sample values associated with the given handle into the array `buf`.
`buf` needs to be pre-allocated by the caller to be of the proper size, which
can be determined using `WGSize()`.

### `int WGCopyT(int handle, double* buf)`

As with `WGCopyS()`, but for the time values.

### `int WGSampleRate(int handle, double* result)`

Reads the sample rate value associated with the handle.


