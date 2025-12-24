package crypter

import (
	"context"
	"encoding/hex"
	"goph_keeper/cli/entity"
	"testing"
)

func TestIntegrationEncodeDecode(t *testing.T) {
	crypter, err := New()
	if err != nil {
		t.Errorf("New Crypter Error: %s", err)
	}
	ctx := context.Background()
	data := "sensitive string"
	user := &entity.User{
		Login:  "Fedor",
		Phrase: "some user phrase",
	}
	crypted, err := crypter.Encode(ctx, []byte(data), user)
	if err != nil {
		t.Errorf("Encode Error: %s", err)
	}
	decodeString, err := hex.DecodeString(crypted)
	if err != nil {
		t.Errorf("Decode Error: %s", err)
	}
	decrypted, err := crypter.Decode(ctx, decodeString, user)
	if err != nil {
		t.Errorf("Decode Error: %s", err)
	}
	if decrypted != data {
		t.Errorf("Decrypted string does not match data = %v, decrypted = %v", data, decrypted)
	}
}
