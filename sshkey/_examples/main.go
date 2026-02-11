package main

import (
	"fmt"

	"github.com/purpose168/charm-experimental-packages-cn/sshkey"
)

func main() {
	// 密码是 "asd"。
	signer, err := sshkey.Open("./key")
	if err != nil {
		panic(err)
	}

	if signer != nil {
		fmt.Println("Key opened!")
	}
}
