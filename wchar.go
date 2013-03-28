package wchar

// #include <wchar.h>
import "C"
import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"log"
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
func WcharPtrToWcharString(first interface{}) (WcharString, error) {
	wcharPtr := first.(*C.wchar_t)
	//++ do stuff
	spew.Dump(wcharPtr)
	return nil, errors.New("not implemented yet")
}

// convert a *C.wchar_t to a Go string
func WcharPtrToGoString(first interface{}) (string, error) {
	ws, err := WcharPtrToWcharString(first)
	if err != nil {
		return "", err
	}
	return ws.GoString(), nil
}

// convert a *C.wchar_t and length int to a Go string
func WcharPtrIntToGoString(first interface{}, length int) string {
	// asert the pointer to be *C.wchar_t (would direct assert to Wchar also work?)
	wcharPtr := first.(*C.wchar_t)

	// allocate new WcharString to fill with data. Only set cap, later use append
	ws := make(WcharString, 0, length)

	// append data using pointer arithmic
	for i := 0; i < length; i++ {
		ws = append(ws, (Wchar)(*wcharPtr))
		break //++ remove this when pointer arithmic is fixed.
	}

	// Convert and return Go string
	return ws.GoString()
}
