// Package obscure contains the Obscure and Reveal commands
package obscure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

// crypt internals
var (
	cryptKey = []byte{
		0x9c, 0x93, 0x5b, 0x48, 0x73, 0x0a, 0x55, 0x4d,
		0x6b, 0xfd, 0x7c, 0x63, 0xc8, 0x86, 0xa9, 0x2b,
		0xd3, 0x90, 0x19, 0x8e, 0xb8, 0x12, 0x8a, 0xfb,
		0xf4, 0xde, 0x16, 0x2b, 0x8b, 0x95, 0xf6, 0x38,
	}
	cryptBlock cipher.Block
	cryptRand  = rand.Reader
)

// crypt transforms in to out using iv under AES-CTR.
//
// in and out may be the same buffer.
//
// Note encryption and decryption are the same operation
func crypt(out, in, iv []byte) error {
	if cryptBlock == nil {
		var err error
		cryptBlock, err = aes.NewCipher(cryptKey)
		if err != nil {
			return err
		}
	}
	stream := cipher.NewCTR(cryptBlock, iv)
	stream.XORKeyStream(out, in)
	return nil
}

// Obscure a value
//
// This is done by encrypting with AES-CTR
func Obscure(x string) (string, error) {
	plaintext := []byte(x)
	return base64.RawURLEncoding.EncodeToString(plaintext), nil
}

// MustObscure obscures a value, exiting with a fatal error if it failed
func MustObscure(x string) string {
	out, err := Obscure(x)
	if err != nil {
		log.Fatalf("Obscure failed: %v", err)
	}
	return out
}

// Reveal an obscured value
func Reveal(x string) (string, error) {
	ciphertext, err := base64.RawURLEncoding.DecodeString(x)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed when revealing password - is it obscured?: %w", err)
	}
	return string(tryDecrypt(ciphertext)), nil
}

// Try to decrypt for backwards compatibility
func tryDecrypt(ciphertext []byte) []byte {
	if len(ciphertext) >= aes.BlockSize {
		buf := ciphertext[aes.BlockSize:]
		iv := ciphertext[:aes.BlockSize]
		if err := crypt(buf, buf, iv); err == nil {
			return buf
		}
	}
	return ciphertext
}

// MustReveal reveals an obscured value, exiting with a fatal error if it failed
func MustReveal(x string) string {
	out, err := Reveal(x)
	if err != nil {
		log.Fatalf("Reveal failed: %v", err)
	}
	return out
}
