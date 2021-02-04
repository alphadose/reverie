package utils

import (
	"bytes"
	"fmt"

	"github.com/reverie/types"
)

// MapToString convers map of type map[string]interface{} to string
func MapToString(m types.M) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s: \"%v\"\n", key, value)
	}
	return b.String()
}
