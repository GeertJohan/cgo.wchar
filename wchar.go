package wchar

// #include <wchar.h>
import "C"
import (
	"log"
	"unsafe"
)

// go representation of a wchar
type Wchar int32

// go representation of a wchar string (array)
type WcharString []Wchar

// create a new Wchar string
// FIXME: why hardcoded length?? Isn't there a better way to do this?
func NewWcharString(length int) WcharString {
	return make(WcharString, length)
}

// create a Go string from a WcharString
func NewWcharStringFromGoString(input string) (WcharString, error) {
	return convertGoStringToWcharString(input)
}

// convert a *C.wchar_t to a WcharString
func NewWcharStringFromWcharPtr(first unsafe.Pointer) WcharString {
	if uintptr(first) == 0x0 {
		return NewWcharString(0)
	}

	// Get uintptr from first wchar_t
	wcharPtr := uintptr(first)

	// allocate new WcharString to fill with data. Cap is unknown
	ws := make(WcharString, 0)

	// append data using pointer arithmic
	var x Wchar
	for {
		// get Wchar value
		x = *((*Wchar)(unsafe.Pointer(wcharPtr)))

		// check for null byte terminator
		if x == 0 {
			break
		}

		// append Wchar to WcharString
		ws = append(ws, x)

		//++ increment pointer
		//++ FIXME: doing this properly??
		wcharPtr += 4
	}

	return ws
}

// convert a *C.wchar_t and length int to a WcharString
func NewWcharStringFromWcharPtrInt(first unsafe.Pointer, length int) WcharString {
	wcharPtr := uintptr(first)

	// allocate new WcharString to fill with data. Only set cap, later use append
	ws := make(WcharString, 0, length)

	// append data using pointer arithmic
	var x Wchar
	for i := 0; i < length; i++ {
		// get Wchar
		x = *((*Wchar)(unsafe.Pointer(wcharPtr)))

		// append Wchar to WcharString
		ws = append(ws, x)

		//++ increment pointer
		//++ FIXME: doing this properly??
		wcharPtr += 4
	}

	return ws
}

// return pointer to first element
func (ws WcharString) Pointer() *Wchar {
	return &ws[0]
}

// convert and get WcharString as Go string
// might return an error when conversion failed.
func (ws WcharString) GoString() (string, error) {
	log.Println("Going to GoString this ws:")
	str, err := convertWcharStringToGoString(ws)
	if err != nil {
		log.Printf("Error at convertWcharStringToGoString(ws): %s\n", err)
		return "", err
	}
	return str, nil
}

// convert a null terminated *C.wchar_t to a Go string
// convenient wrapper for WcharPtrToWcharString(first).GoString()
func WcharPtrToGoString(first unsafe.Pointer) (string, error) {
	if uintptr(first) == 0x0 {
		return "", nil
	}
	return convertWcharStringToGoString(NewWcharStringFromWcharPtr(first))
}

// convert a *C.wchar_t and length int to a Go string
// convenient wrapper for WcharPtrIntToWcharString(first, length).GoString()
func WcharPtrIntToGoString(first unsafe.Pointer, length int) (string, error) {
	if uintptr(first) == 0x0 {
		return "", nil
	}
	return convertWcharStringToGoString(NewWcharStringFromWcharPtrInt(first, length))
}
