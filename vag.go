package main

import (
	"fmt"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/ripemd160"
	"os"
	"encoding/base32"
)

func main() {
	pub, priv, err := ed25519.GenerateKey(nil)
	fmt.Printf("Public: %x, Private: %x\n", pub, priv)
	if err != nil {
		os.Exit(-1)
	}
	h := sha3.Sum256(pub)
	fmt.Println("SHA3", h)

	md := ripemd160.New()
	md.Write(h[:])

	s := md.Sum(nil)
	fmt.Println("Ripemd", s)

	s = append([]byte {0x68}, s...)
	h = sha3.Sum256(s)

	address := append(s, h[:4]...)
	fmt.Printf("Address %v\n", address)
	account := base32.StdEncoding.EncodeToString(address)

	fmt.Printf("Account %v", account)
}
