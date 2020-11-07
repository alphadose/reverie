package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/reverie/configs"
)

// The key and the nonce used for encryption
var (
	key   []byte
	nonce []byte
)

func init() {
	tmpKey, err := hex.DecodeString(configs.Project.Crypto.Key)
	if err != nil {
		LogError("", err)
		os.Exit(1)
	}
	tmpNonce, err := hex.DecodeString(configs.Project.Crypto.Nonce)
	if err != nil {
		LogError("", err)
		os.Exit(1)
	}
	key = tmpKey
	nonce = tmpNonce
}

// GenerateRandomKey returns a random 64 character key used for AES-256 encryption
func GenerateRandomKey() (string, error) {
	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil //encode key in bytes to string and keep as secret, put in a vault
}

// GenerateNonce returns a random 24 character nonce used for AES-256 encryption
func GenerateNonce() (string, error) {
	bytes := make([]byte, 12) //generate a random 12 byte nonce for AES-256
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil //encode key in bytes to string and keep as secret, put in a vault
}

// Encrypt encrypts a string using AES-256 encryption
func Encrypt(stringToEncrypt string) (string, error) {
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data
	// The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt decrypts a string using AES-256 encryption
func Decrypt(encryptedString string) (string, error) {
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}
