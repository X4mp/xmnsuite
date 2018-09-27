package main

import "C"
import (
	"unsafe"

	seedwords "github.com/xmnservices/xmnsuite/seedwords"
)

//export xGetSupportedLangs
func xGetSupportedLangs(ptr **C.char) {
	const arrayLen = 1<<30 - 1
	slice := (*[arrayLen]*C.char)(unsafe.Pointer(ptr))[:arrayLen:arrayLen]

	langs := seedwords.SDKFunc.Create().GetSupportedLanguages()
	for index, oneLang := range langs {
		slice[index] = C.CString(oneLang)
	}
}

//export xGetWords
func xGetWords(ptr **C.char, lang string, amount int) {
	const arrayLen = 1<<30 - 1
	slice := (*[arrayLen]*C.char)(unsafe.Pointer(ptr))[:arrayLen:arrayLen]

	index := 0
	words := seedwords.SDKFunc.Create().GetWords(lang, amount)
	for _, oneWord := range words {
		slice[index] = C.CString(oneWord)
		index++
	}
}

func main() {}
