package wchar

import (
	"testing"
)

func TestStringConversion(t *testing.T) {
	testString := "Iñtërnâtiônàlizætiøn"
	expectedWcharString := WcharString{
		73, 241, 116, 235, 114,
		110, 226, 116, 105, 244,
		110, 224, 108, 105, 122,
		230, 116, 105, 248, 110, 0,
	}

	w, err := FromGoString(testString)
	if err != nil {
		t.Fatalf("Error on conversion. %s", err.Error())
	}

	if len(w) != len(expectedWcharString) {
		t.Fatal("Converted string did not match expected WcharString. Lengths are different.")
	}

	for i := 0; i < len(w); i++ {
		if w[i] != expectedWcharString[i] {
			t.Fatalf("Converted string did not match expected WcharString. Fault at position %d. %d!=%d\n", i, w[i], expectedWcharString[i])
		}
	}
}

func TestWcharStringConversion(t *testing.T) {
	testWcharString := WcharString{
		73, 241, 116, 235, 114,
		110, 226, 116, 105, 244,
		110, 224, 108, 105, 122,
		230, 116, 105, 248, 110, 0,
	}
	expectedGoString := "Iñtërnâtiônàlizætiøn"

	str, err := testWcharString.GoString()
	if err != nil {
		t.Fatalf("Error on conversion. %s", err.Error())
	}

	if len(str) != len(expectedGoString) {
		t.Fatal("Converted WcharString did not match expected string. Lengths are different.")
	}

	for i := 0; i < len(str); i++ {
		if str[i] != expectedGoString[i] {
			t.Fatalf("Converted WcharString did not match expected string. Fault at position %d. %d!=%d\n", i, str[i], expectedGoString[i])
		}
	}
}

func TestRuneConversion(t *testing.T) {
	testRunes := []rune{
		'I', 'ñ', 't', 'ë', 'r',
		'n', 'â', 't', 'i', 'ô',
		'n', 'à', 'l', 'i', 'z',
		'æ', 't', 'i', 'ø', 'n',
	}
	expectedWchars := []Wchar{
		73, 241, 116, 235, 114,
		110, 226, 116, 105, 244,
		110, 224, 108, 105, 122,
		230, 116, 105, 248, 110,
	}

	for i, testRune := range testRunes {
		w, err := FromGoRune(testRune)
		if err != nil {
			t.Fatalf("Error on conversion. %s", err.Error())
		}
		if w != expectedWchars[i] {
			t.Fatalf("Converted rune did not match expected Wchar. Fault at position %d. %d!=%d\n", i, w, expectedWchars[i])
		}
	}
}

func TestWcharConversion(t *testing.T) {
	testWchars := []Wchar{
		73, 241, 116, 235, 114,
		110, 226, 116, 105, 244,
		110, 224, 108, 105, 122,
		230, 116, 105, 248, 110,
	}
	expectedRunes := []rune{
		'I', 'ñ', 't', 'ë', 'r',
		'n', 'â', 't', 'i', 'ô',
		'n', 'à', 'l', 'i', 'z',
		'æ', 't', 'i', 'ø', 'n',
	}

	for i, testWchar := range testWchars {
		r, err := testWchar.GoRune()
		if err != nil {
			t.Fatalf("Error on conversion. %s", err.Error())
		}
		if r != expectedRunes[i] {
			t.Fatalf("Converted Wchar did not match expected rune. Fault at position %d. %d!=%d\n", i, r, expectedRunes[i])
		}
	}
}
