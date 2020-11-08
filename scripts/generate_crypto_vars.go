package main

import (
	"fmt"

	"github.com/reverie/utils"
)

func main() {
	key, _ := utils.GenerateRandomKey()
	nonce, _ := utils.GenerateNonce()
	fmt.Println("key: ", key)
	fmt.Println("nonce: ", nonce)
}
