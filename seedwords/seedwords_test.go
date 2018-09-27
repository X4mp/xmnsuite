package seedwords

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestCreateSeedWords_Success(t *testing.T) {
	sw := createSeedWords()
	langs := sw.GetSupportedLanguages()

	expectedLanguages := []string{"cn", "en", "es", "fr", "jp", "pt"}
	if !reflect.DeepEqual(langs, expectedLanguages) {
		t.Errorf("the returned languages are invalid")
		return
	}

	for _, oneLang := range langs {
		amount := rand.Int() % 100
		words := sw.GetWords(oneLang, amount)
		if len(words) != amount {
			t.Errorf("the returned amount of words were expected to be %d, returned: %d", amount, len(words))
			return
		}
	}

	// unsupported language:
	unsupportedWords := sw.GetWords("us", 20)
	if len(unsupportedWords) != 0 {
		t.Errorf("the returned words were expected to be empty because the language is not supported.")
		return
	}

	// amount too big:
	tooBigWords := sw.GetWords("fr", 5000000000000000)
	if len(tooBigWords) >= 5000000000000 {
		t.Errorf("the returned words should not exceed the total amount of words of the language")
		return
	}
}
