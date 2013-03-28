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

// Create a new Wchar string
// FIXME: why hardcoded length?? Isn't there a better way to do this?
func NewWcharString(length int) WcharString {
	return make(WcharString, length)
}

// Return pointer to first element
func (ws WcharString) Pointer() *Wchar {
	return &ws[0]
}

// Convert and get WcharString as Go string
func (ws WcharString) GoString() string {
	str, err := convertWcharStringToGoString(ws)
	if err != nil {
		log.Printf("Error at convertWcharStringToGoString(ws): %s\n", err)
	}
	return str
}

// Create a Go string from a WcharString
func GoStringToWcharString(input string) (WcharString, error) {
	//++
	return nil, errors.New("not implemented yet")
}

// Convert a *C.wchar_t to a WcharString
func WcharPtrToWcharString(first interface{}) (WcharString, error) {
	wcharPtr := first.(*C.wchar_t)
	//++ do stuff
	spew.Dump(wcharPtr)
	return nil, errors.New("not implemented yet")
}

// Convert a *C.wchar_t to a Go string
func WcharPtrToGoString(first interface{}) (string, error) {
	ws, err := WcharPtrToWcharString(first)
	if err != nil {
		return "", err
	}
	return ws.GoString(), nil
}

// Convert a *C.wchar_t and length int to a Go string
func WcharPtrIntToGoString(first interface{}, length int) string {
	// asert the pointer to be *C.wchar_t (would direct assert to Wchar also work?)
	wcharPtr := first.(*C.wchar_t)

	// allocate new WcharString to fill with data
	ws := NewWcharString(length)

	//++ do pointer arithmic and store values in array
	ws[0] = (Wchar)(*wcharPtr)

	// Convert and return Go string
	return ws.GoString()
}
