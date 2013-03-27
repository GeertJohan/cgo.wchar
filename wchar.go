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

	// calculate bufferSizes in bytes
	bufferSizeIn := len([]byte(input)) // count exact amount of bytes
	bufferSizeOut := len(input) * 4    // wide char seems to be 4 bytes for every single- or multi-byte character. Not very sure though.

	// bufferSizes for C (copies, because iconv will touch them)
	bytesLeftIn := bufferSizeIn
	bytesLeftInCSize := C.size_t(bytesLeftIn)
	bytesLeftOut := bufferSizeOut
	bytesLeftOutCSize := C.size_t(bytesLeftOut)

	// input for C
	inputCString := C.CString(input)
	defer C.free(unsafe.Pointer(inputCString))

	// create output buffer
	outputChars := make([]int8, bufferSizeOut)

	// output for C
	outputCString := (*C.char)(&outputChars[0])

	x, errno := C.iconv(iconv, &inputCString, &bytesLeftInCSize, &outputCString, &bytesLeftOutCSize)
	spew.Dump(x)
	spew.Dump(errno)

	output = make([]C.wchar_t, 0, len(input))
	for len(outputChars) >= 4 {
		// create 4 position byte slice
		b := make([]byte, 4)
		b[0] = byte(outputChars[0])
		b[1] = byte(outputChars[1])
		b[2] = byte(outputChars[2])
		b[3] = byte(outputChars[3])
		uchar := binary.LittleEndian.Uint32(b)
		if uchar == 0 { // find null terminator (doing this right?)
			break
		}
		// Combine 4 position byte slice into uint32, and append uint32 to outputUint32
		output = append(output, C.wchar_t(uchar))
		// reslice the outputChars
		outputChars = outputChars[4:]
	}

	return output, nil
}

func ToGoString(input *C.wchar_t) (output string, err error) {
	// open iconv
	iconv, errno := C.iconv_open(strUtf8, strWchar)
	if iconv == nil || errno != nil {
		return "", fmt.Errorf("Could not create iconv instance: %s", errno)
	}
	defer C.iconv_close(iconv)

	inputAsChars := make([]C.char, 0)
	for {
		nextWchar := *input
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(nextWchar))
		spew.Dump(b)
		if b[0] == 0 { // find null terminator (doing this right?)
			break
		}
		//++ do something with b
		//++ split b into seperate bytes, make those int8's.. make those int8's C.char again. Add the C.char to inputAsChars

		break // remove when pointer arithmic is working and actual loop can be done
	}

	// input for C
	inputAsCharsFirst := &inputAsChars[0]

	// calculate buffer size for input
	bufferSizeIn := len(inputAsChars)
	bufferLeftIn := bufferSizeIn
	bytesLeftInCSize := C.size_t(bufferLeftIn)

	// calculate buffer size for output
	bufferSizeOut := len(inputAsChars)
	bytesLeftOut := bufferSizeOut
	bytesLeftOutCSize := C.size_t(bytesLeftOut)

	// create output buffer
	outputChars := make([]int8, bufferSizeOut)

	// output for C
	outputCString := (*C.char)(&outputChars[0])

	x, errno := C.iconv(iconv, &inputAsCharsFirst, &bytesLeftInCSize, &outputCString, &bytesLeftOutCSize)
	spew.Dump(x)
	spew.Dump(errno)

	//++ do processing on outputChars

	return "todo", nil
}
