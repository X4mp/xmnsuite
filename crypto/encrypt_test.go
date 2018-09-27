package crypto

import (
	"reflect"
	"testing"
)

func TestEncrypt_decrypt_Success(t *testing.T) {
	// variables:
	pass := []byte("this is the password used to generate the pk")
	textToEncrypt := []byte("this is some text to encrypt... this is even longer text!")

	encryptedText, encryptedTextErr := encrypt(pass, textToEncrypt)
	if encryptedTextErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", encryptedTextErr.Error())
		return
	}

	decrypted, decryptedErr := decrypt(pass, encryptedText)
	if decryptedErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", decryptedErr.Error())
		return
	}

	if !reflect.DeepEqual(textToEncrypt, decrypted) {
		t.Errorf("the decrypted text was excpected to be the same as the original text")
		return
	}
}
