package crypter

import (
	"context"
	"encoding/hex"
	"goph_keeper/cli/entity"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const path = "./test"
const newPath = "./test-generated"

func TestCryptDecryptBin(t *testing.T) {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("read file err: %v", err)
	}
	crypter, err := New()
	if err != nil {
		t.Errorf("New Crypter Error: %s", err)
	}
	ctx := context.Background()
	user := &entity.User{
		Login:  "Fedor",
		Phrase: "some user phrase",
	}

	crypted, err := crypter.Encode(ctx, data, user)
	decodeString, err := hex.DecodeString(crypted)
	if err != nil {
		t.Errorf("Decode Error: %s", err)
	}
	decrypted, err := crypter.Decode(ctx, decodeString, user)
	if err != nil {
		t.Errorf("Decode Error: %s", err)
	}
	err = os.WriteFile(newPath, []byte(decrypted), 0777)
	if err != nil {
		t.Errorf("WriteFile Error: %s", err)
		return
	}

	out, err := exec.Command(newPath).Output()
	if err != nil {
		t.Error(err)
		return
	}
	if len(out) == 0 {
		t.Errorf("empty output")
		return
	}
	if strings.EqualFold(string(out), "Hello World") {
		t.Errorf("wrong output: %s", out)
	}
	err = os.Remove(newPath)
	if err != nil {
		t.Errorf("Remove Error: %s", err)
	}
}
