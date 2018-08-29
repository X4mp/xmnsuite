package tendermint

type someDataForTests struct {
	Title string `json:"title"`
	Desc  string `json:"description"`
}

func createSomeDataForTests(title string, desc string) *someDataForTests {
	out := someDataForTests{
		Title: title,
		Desc:  desc,
	}

	return &out
}
