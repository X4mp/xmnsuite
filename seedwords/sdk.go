package seedwords

// SeedWords represents the seedwords interface
type SeedWords interface {
	GetSupportedLanguages() []string
	GetWords(lang string, amount int) []string
}

// SDKFunc represents the seed words SDK func
var SDKFunc = struct {
	Create func() SeedWords
}{
	Create: func() SeedWords {
		return createSeedWords()
	},
}
