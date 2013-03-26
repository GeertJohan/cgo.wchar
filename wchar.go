package wchar

// #include <stdlib.h>
// #include <wchar.h>
// #include <iconv.h>
import "C"
import (
	"encoding/binary"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"unsafe"
)

var (
	strWchar = C.CString("wchar_t//TRANSLIT")
	strChar  = C.CString("//TRANSLIT")
	strAscii = C.CString("ascii//TRANSLIT")
	strUtf8  = C.CString("utf-8//TRANSLIT")
)

// Use iconv. It seems to support conversion between char and wchar_t
// http://www.gnu.org/savannah-checkouts/gnu/libiconv/documentation/libiconv-1.13/iconv_open.3.html
// http://www.gnu.org/savannah-checkouts/gnu/libiconv/documentation/libiconv-1.13/iconv.3.html
// http://www.gnu.org/savannah-checkouts/gnu/libiconv/documentation/libiconv-1.13/iconv_close.3.html

func FromGoString(input string) (output []C.wchar_t, err error) {
	// open iconv
	iconv, errno := C.iconv_open(strWchar, strUtf8)
	if iconv == nil || errno != nil {
		return nil, fmt.Errorf("Could not create iconv instance: %s", errno)
	}
	defer C.iconv_close(iconv)

	// calculate bufferSize in bytes
	bufferSize := len([]byte(input))

	// bufferSizes for C
	bytesLeftIn := bufferSize
	bytesLeftInCSize := C.size_t(bytesLeftIn)
	bytesLeftOut := bufferSize * 4 //FIXME: wide chars seems to be 4 bytes for 1 char. Not very sure though.
	bytesLeftOutCSize := C.size_t(bytesLeftOut)

	// input for C
	inputCString := C.CString(input)
	defer C.free(unsafe.Pointer(inputCString))

	// create output buffer
	outputChars := make([]int8, bytesLeftOut)

	// output for C
	outputCString := (*C.char)(&outputChars[0])

	x, errno := C.iconv(iconv, &inputCString, &bytesLeftInCSize, &outputCString, &bytesLeftOutCSize)
	spew.Dump(x)
	spew.Dump(errno)

	output = make([]C.wchar_t, 0, bufferSize)
	for len(outputChars) > bytesLeftOut/4 {
		// create 4 position byte slice
		b := make([]byte, 4)
		b[0] = byte(outputChars[0])
		b[1] = byte(outputChars[1])
		b[2] = byte(outputChars[2])
		b[3] = byte(outputChars[3])
		// Combine 4 position byte slice into uint32, and append uint32 to outputUint32
		output = append(output, C.wchar_t(binary.LittleEndian.Uint32(b)))
		// reslice the outputChars
		outputChars = outputChars[4:]
	}

	a := int32(12)
	b := (C.wchar_t)(a)
	spew.Dump(b)

	return output, nil
}
