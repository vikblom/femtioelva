package femtioelva_test

import (
	"testing"

	"github.com/vikblom/femtioelva"
)

func TestRoundTrip(t *testing.T) {
	pass := "foo"
	msg := "lorem ipsum"

	key := femtioelva.GenerateKey(pass)

	cipher, err := femtioelva.Encrypt(msg, key)
	if err != nil {
		t.Fatal("Encryption error:", err)
	}
	out, err := femtioelva.Decrypt(cipher, key)
	if err != nil {
		t.Fatal("Decryption error:", err)
	}

	if out != msg {
		t.Fatalf("Expected '%s' got '%s'\n", msg, out)
	}

}
