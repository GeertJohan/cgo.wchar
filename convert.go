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

// iconv charset strings
var (
	iconvCharsetWchar = C.CString("wchar_t//TRANSLIT")
	iconvCharsetChar  = C.CString("//TRANSLIT")
	iconvCharsetAscii = C.CString("ascii//TRANSLIT")
	iconvCharsetUtf8  = C.CString("utf-8//TRANSLIT")
)

// iconv documentation:
// Use iconv. It seems to support conversion between char and wchar_t
// http://www.gnu.org/savannah-checkouts/gnu/libiconv/documentation/libiconv-1.13/iconv_open.3.html
// http://www.gnu.org/savannah-checkouts/gnu/libiconv/documentation/libiconv-1.13/iconv.3.html
// http://www.gnu.org/savannah-checkouts/gnu/libiconv/documentation/libiconv-1.13/iconv_close.3.html

func findSize(first *C.wchar_t) int {
	//++ find size for C.wchar_t string
	return 0
}

// Internal helper function, wrapped by several other functions
func convertGoStringToWcharString(input string) (output WcharString, err error) {
	// open iconv
	iconv, errno := C.iconv_open(iconvCharsetWchar, iconvCharsetUtf8)
	if iconv == nil || errno != nil {
		return nil, fmt.Errorf("Could not create iconv instance: %s", errno)
	}
	defer C.iconv_close(iconv)

	// calculate bufferSizes in bytes
	bufferSizeIn := len([]byte(input)) // count exact amount of bytes from input
	bufferSizeOut := len(input) * 4    // wide char seems to be 4 bytes for every single- or multi-byte character. Not very sure though.

	// bufferSizes for C
	bytesLeftIn := bufferSizeIn
	bytesLeftInCSize := C.size_t(bytesLeftIn)
	bytesLeftOut := bufferSizeOut
	bytesLeftOutCSize := C.size_t(bytesLeftOut)

	// input for C. makes a copy using C malloc and therefore should be free'd.
	inputCString := C.CString(input)
	defer C.free(unsafe.Pointer(inputCString))

	// create output buffer
	outputChars := make([]int8, bufferSizeOut)

	// output for C
	outputCString := (*C.char)(&outputChars[0])

	x, errno := C.iconv(iconv, &inputCString, &bytesLeftInCSize, &outputCString, &bytesLeftOutCSize)
	spew.Dump(x)
	spew.Dump(errno)
	if errno != nil {
		return nil, errno
	}

	// Convert []int8 to WcharString
	output = make(WcharString, 0, len(input))
	wcharAsByteAry := make([]byte, 4)
	for len(outputChars) >= 4 {
		// create 4 position byte slice
		wcharAsByteAry[0] = byte(outputChars[0])
		wcharAsByteAry[1] = byte(outputChars[1])
		wcharAsByteAry[2] = byte(outputChars[2])
		wcharAsByteAry[3] = byte(outputChars[3])
		whcarAsUint32 := binary.LittleEndian.Uint32(wcharAsByteAry)
		// find null terminator (doing this right?)
		if whcarAsUint32 == 0 {
			break
		}
		// Combine 4 position byte slice into uint32, and append uint32 to outputUint32
		output = append(output, Wchar(whcarAsUint32))
		// reslice the outputChars
		outputChars = outputChars[4:]
	}

	return output, nil
}

func convertWcharStringToGoString(input WcharString) (output string, err error) {
	// open iconv
	iconv, errno := C.iconv_open(iconvCharsetUtf8, iconvCharsetWchar)
	if iconv == nil || errno != nil {
		return "", fmt.Errorf("Could not create iconv instance: %s", errno.Error())
	}
	defer C.iconv_close(iconv)

	inputAsCChars := make([]C.char, 0)
	wcharAsBytes := make([]byte, 4)
	for _, nextWchar := range input {
		// find null terminator (doing this right?)
		if nextWchar == 0 {
			// Return empty string if there are no chars in buffer
			if len(inputAsCChars) == 0 {
				return "", nil
			}
			break
		}

		// split Wchar into bytes
		binary.LittleEndian.PutUint32(wcharAsBytes, uint32(nextWchar))

		//++ split b into seperate bytes, make those int8's.. make those int8's C.char again. Add the C.char to inputAsCChars
		for i := 0; i < 4; i++ {
			inputAsCChars = append(inputAsCChars, C.char(wcharAsBytes[i]))
		}
	}

	// input for C
	inputAsCharsFirst := &inputAsCChars[0]

	// calculate buffer size for input
	bufferSizeIn := len(inputAsCChars)
	bufferLeftIn := bufferSizeIn
	bytesLeftInCSize := C.size_t(bufferLeftIn)

	// calculate buffer size for output
	bufferSizeOut := len(inputAsCChars)
	bytesLeftOut := bufferSizeOut
	bytesLeftOutCSize := C.size_t(bytesLeftOut)

	// create output buffer
	outputChars := make([]int8, bufferSizeOut)

	// output for C
	outputCString := (*C.char)(&outputChars[0])

	_, errno = C.iconv(iconv, &inputAsCharsFirst, &bytesLeftInCSize, &outputCString, &bytesLeftOutCSize)
	if errno != nil {
		return "", errno
	}

	// conver output buffer to go string
	str := C.GoString((*C.char)(&outputChars[0]))

	return str, nil
}
