package seedwords

import (
	"math/rand"
	"sort"
	"time"
)

type seedWords struct {
	words map[string][]string
}

func createSeedWords() SeedWords {
	out := seedWords{
		words: map[string][]string{
			"cn": getChineseWords(),
			"en": getEnglishWords(),
			"fr": getFrenchWords(),
			"jp": getJapaneseWords(),
			"pt": getPortugeseWords(),
			"es": getSpanishWords(),
		},
	}

	return &out
}

// GetSupportedLanguages returns the supported languages
func (app *seedWords) GetSupportedLanguages() []string {
	langs := []string{}
	for oneLang := range app.words {
		langs = append(langs, oneLang)
	}

	sort.Strings(langs)
	return langs
}

// GetWords retrieves an amount of words of a specific langauge.  If the language is not supported, the list returns empty
func (app *seedWords) GetWords(lang string, amount int) []string {
	if words, ok := app.words[lang]; ok {
		return app.getAmount(words, amount)
	}

	return []string{}
}

func (app *seedWords) getAmount(words []string, amount int) []string {
	if len(words) <= amount {
		return words
	}

	out := []string{}
	for i := 0; i < amount; i++ {
		ts := time.Now().UnixNano()
		source := rand.NewSource(ts)
		rnd := rand.New(source)

		index := rnd.Intn(len(words) - 1)
		out = append(out, words[index])
	}

	return out
}
