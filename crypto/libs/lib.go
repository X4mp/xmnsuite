package main

/*

 */
import "C"

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/xmnservices/xmnsuite/crypto"
)

//export xGenPK
func xGenPK() *C.char {
	gen := crypto.SDKFunc.GenPK()
	return C.CString(gen.String())
}

//export xGenEncryptedPk
func xGenEncryptedPk(seedWords **C.char, amount C.int) *C.char {
	// retrieve the words from C:
	const arrayLen = 1<<30 - 1
	sliceWithCCHars := (*[arrayLen]*C.char)(unsafe.Pointer(seedWords))[:arrayLen:arrayLen]

	// convert the C string to Go strings:
	words := []string{}
	goAmount := int(amount)
	for i := 0; i < goAmount; i++ {
		words = append(words, C.GoString(sliceWithCCHars[i]))
	}

	//generate the PK:
	pk := crypto.SDKFunc.GenPK()

	// encrypt the PK:
	encPK := crypto.SDKFunc.Encrypt(crypto.EncryptParams{
		Pass: []byte(strings.Join(words, "|")),
		Msg:  []byte(pk.String()),
	})

	return C.CString(encPK)
}

//export xDecrypt
func xDecrypt(encryptedPK *C.char, passWords **C.char, amountPassWords C.int) (ret *C.char) {

	// returns an empty string if the PK can't be decoded using the given pass words:
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("xDecrypt panic: ", r)
			ret = C.CString("")
		}
	}()

	// retrieve the words from C:
	const arrayLen = 1<<30 - 1
	sliceWithCCHars := (*[arrayLen]*C.char)(unsafe.Pointer(passWords))[:arrayLen:arrayLen]

	// convert the C string to Go strings:
	passsWords := []string{}
	goAmount := int(amountPassWords)
	for i := 0; i < goAmount; i++ {
		passsWords = append(passsWords, C.GoString(sliceWithCCHars[i]))
	}

	// decrypt the pk:
	pkAsString := string(crypto.SDKFunc.Decrypt(crypto.DecryptParams{
		Pass:         []byte(strings.Join(passsWords, "|")),
		EncryptedMsg: C.GoString(encryptedPK),
	}))

	// make sure the decrypted pk is a real PK:
	decPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
		PKAsString: pkAsString,
	})

	return C.CString(decPK.String())
}

//export xPKGetPublicKey
func xPKGetPublicKey(pk *C.char) *C.char {
	lpk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
		PKAsString: C.GoString(pk),
	})

	pubKey := lpk.PublicKey().String()
	return C.CString(pubKey)
}

//export xPKSign
func xPKSign(pk *C.char, msg *C.char) *C.char {
	lpk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
		PKAsString: C.GoString(pk),
	})

	sig := lpk.Sign(C.GoString(msg))
	return C.CString(sig.String())
}

//export xPKRingSign
func xPKRingSign(pk *C.char, msg *C.char, ringPubKeys []*C.char) *C.char {
	lpk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
		PKAsString: C.GoString(pk),
	})

	pubKeys := []crypto.PublicKey{}
	for _, onePubKey := range ringPubKeys {
		pubKeys = append(pubKeys, crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
			PubKeyAsString: C.GoString(onePubKey),
		}))
	}

	ringSig, ringSigErr := lpk.RingSign(C.GoString(msg), pubKeys)
	if ringSigErr != nil {
		return nil
	}

	return C.CString(ringSig.String())
}

func main() {
}
