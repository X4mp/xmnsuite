package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func hashPass(pass []byte) []byte {
	hasher := curve.Hash()
	hasher.Write([]byte(pass))
	return hasher.Sum(nil)
}

func encrypt(pass []byte, msg []byte) (string, error) {
	block, blockErr := aes.NewCipher(hashPass(pass))
	if blockErr != nil {
		return "", blockErr
	}

	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], msg)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(pass []byte, encryptedText string) ([]byte, error) {
	cipherText, cipherTextErr := base64.StdEncoding.DecodeString(encryptedText)
	if cipherTextErr != nil {
		return nil, cipherTextErr
	}

	block, blockErr := aes.NewCipher(hashPass(pass))
	if blockErr != nil {
		return nil, blockErr
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("the encrypted text cannot be decoded using this password: ciphertext block size is too short")
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	// returns the decoded message:
	return cipherText, nil
}
