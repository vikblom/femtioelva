package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vikblom/femtioelva"
)



func main() {
	var pass string
	flag.StringVar(&pass, "pass", "", "passphrase")
	flag.Parse()

	key := femtioelva.GenerateKey(pass)

	var payload string
	if len(flag.Args()) == 1 {
		payload = flag.Arg(0)
	} else {
		payloadBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Failed to read payload from stdin:", err)
		}
		payload = string(payloadBytes)
	}

	cipher, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		log.Fatal("Base64 encoding error:", err)
	}

	plain, err := femtioelva.Decrypt(cipher, key)
	if err != nil {
		log.Fatal("Decryption error:", err)
	}
	fmt.Println(plain)
}
