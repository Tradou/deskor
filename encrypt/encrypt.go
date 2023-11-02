package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

//openssl rand -out encrypt.key 32

type Encrypter interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
	readKey(keyPath string) ([]byte, error)
}

type EncrypterImpl struct {
	//KeyPath string
}

func (encrypt *EncrypterImpl) readKey() ([]byte, error) {
	key, err := os.ReadFile("./cert/encrypt.key")
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (encrypt *EncrypterImpl) Encrypt(plaintext string) (string, error) {
	key, err := encrypt.readKey()
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	ciphertextString := base64.StdEncoding.EncodeToString(ciphertext)

	return ciphertextString, nil
}

func (encrypt *EncrypterImpl) Decrypt(ciphertext string) (string, error) {
	key, err := encrypt.readKey()
	if err != nil {
		return "", err
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext is too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	mode := cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(ciphertextBytes, ciphertextBytes)

	plaintext := string(ciphertextBytes)

	return plaintext, nil
}
