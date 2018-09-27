package seedwords

import (
	"reflect"
	"testing"
)

func TestSDK_createSeedWords_Success(t *testing.T) {
	sw := SDKFunc.Create()
	langs := sw.GetSupportedLanguages()

	expectedLanguages := []string{"cn", "en", "es", "fr", "jp", "pt"}
	if !reflect.DeepEqual(langs, expectedLanguages) {
		t.Errorf("the returned languages are invalid")
		return
	}
}
