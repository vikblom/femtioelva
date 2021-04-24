package femtioelva

import "testing"

func TestRounTrip(t *testing.T) {
	pass := "foo"
	msg := "lorem ipsum"

	key := GenerateKey(pass)

	cipher, err := Encrypt(msg, key)
	if err != nil {
		t.Fatal("Encryption error:", err)
	}
	out, err := Decrypt(cipher, key)
	if err != nil {
		t.Fatal("Decryption error:", err)
	}

	if (out != msg) {
		t.Fatalf("Expected '%s' got '%s'\n", msg, out)
	}

}
