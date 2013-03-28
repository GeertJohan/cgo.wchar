package wchar

// #include <wchar.h>
import "C"
import (
	"errors"
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

// return pointer to first element
func (ws WcharString) Pointer() *Wchar {
	return &ws[0]
}

// convert and get WcharString as Go string
func (ws WcharString) GoString() string {
	str, err := convertWcharStringToGoString(ws)
	if err != nil {
		log.Printf("Error at convertWcharStringToGoString(ws): %s\n", err)
	}
	return str
}

// create a Go string from a WcharString
func GoStringToWcharString(input string) (WcharString, error) {
	return convertGoStringToWcharString(input)
}

// convert a *C.wchar_t to a WcharString
func WcharPtrToWcharString(first unsafe.Pointer) (WcharString, error) {
	wcharPtr := uintptr(first)

	// allocate new WcharString to fill with data. Cap is unknown
	ws := make(WcharString, 0)

	// append data using pointer arithmic
	var x Wchar
	for {
		// get Wchar
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

		break //++ remove this when pointer arithmic is fixed.
	}

	return nil, errors.New("not implemented yet") //++ TODO: is return error requirerd?
}

// convert a *C.wchar_t and length int to a WcharString
func WcharPtrIntToWcharString(first unsafe.Pointer, length int) (WcharString, error) {
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

		break //++ remove this when pointer arithmic is fixed.
	}

	return nil, errors.New("not imlpemented yet") //++ TODO: is return error requirerd?
}

// convert a *C.wchar_t to a Go string
func WcharPtrToGoString(first unsafe.Pointer) (string, error) {
	ws, err := WcharPtrToWcharString(first)
	if err != nil {
		return "", err
	}
	return ws.GoString(), nil
}

// convert a *C.wchar_t and length int to a Go string
func WcharPtrIntToGoString(first unsafe.Pointer, length int) (string, error) {
	ws, err := WcharPtrIntToWcharString(first, length)
	if err != nil {
		return "", err
	}

	// Convert and return Go string
	return ws.GoString(), nil
}
