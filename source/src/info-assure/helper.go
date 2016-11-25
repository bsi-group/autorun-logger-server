package main

import (
	"io/ioutil"
	//"strings"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	//"unicode"
	"bytes"
	"strings"
)

//
func DecodeUtf16AsUtf8(input []byte) []byte {
	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(input), utf16bom)

	// decode and print:
	decoded, _ := ioutil.ReadAll(unicodeReader)
	return decoded //, err

	//reader, err := charset.NewReader("utf16", strings.NewReader(string(input)))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//output, err := ioutil.ReadAll(reader)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//return output
}

// Determines if the alert is a false positive
func IsFalsePositive(isSigned bool, signer string) bool {

	if isSigned == false {
		return false
	}

	switch strings.ToLower(signer) {
	case "microsoft windows hardware compatibility publisher",
		"microsoft windows",
		"microsoft windows publisher",
		"microsoft corporation":
		return true
	}

	return false
}
