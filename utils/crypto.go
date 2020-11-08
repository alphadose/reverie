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

// This module is used for encrypting/decrypting vendor emails
// Flow with notation:-
// vendorEmail -> utils.Encrypt -> offerKey
// offerKey -> utils.Decrypt -> vendorEmail

// The lengths of key and nonce
const (
	keyLength   = 32
	nonceLength = 12
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
	bytes := make([]byte, keyLength) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil //encode key in bytes to string and keep as secret, put in a vault
}

// GenerateNonce returns a random 24 character nonce used for AES-256 encryption
func GenerateNonce() (string, error) {
	bytes := make([]byte, nonceLength) //generate a random 12 byte nonce for AES-256
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil //encode key in bytes to string and keep as secret, put in a vault
}

// Encrypt encrypts a string using AES-256 encryption
func Encrypt(stringToEncrypt string) (string, error) {
	// Making it thread safe
	keyCopy := make([]byte, keyLength)
	nonceCopy := make([]byte, nonceLength)
	copy(keyCopy, key)
	copy(nonceCopy, nonce)

	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(keyCopy)
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
	ciphertext := aesGCM.Seal(nonceCopy, nonceCopy, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt decrypts a string using AES-256 encryption
func Decrypt(encryptedString string) (string, error) {
	// Making it thread safe
	keyCopy := make([]byte, keyLength)
	copy(keyCopy, key)

	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(keyCopy)
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
	nonceCopy, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonceCopy, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}
