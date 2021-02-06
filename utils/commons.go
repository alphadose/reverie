package utils

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/reverie/types"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GenerateRandomString returns a random string of fixed length
func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// MapToString convers map of type map[string]interface{} to string
func MapToString(m types.M) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s: \"%v\"\n", key, value)
	}
	return b.String()
}
