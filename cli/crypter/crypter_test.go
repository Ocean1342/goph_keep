package crypter

import (
	"encoding/hex"
	"testing"
)

func Test_prepareKey(t *testing.T) {
	tests := []struct {
		name    string
		phrase  string
		wantErr bool
	}{
		{
			name:    "positive short phrase",
			phrase:  "short",
			wantErr: false,
		},
		{
			name:    "positive long phrase",
			phrase:  "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
			wantErr: false,
		},
		{
			name:    "empty phrase",
			phrase:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := prepareKey(tt.phrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			//если есть ошибка, то длина не важна
			if len(key) != 32 && !tt.wantErr {
				t.Errorf("generateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func Test_prepareKeyStability(t *testing.T) {
	key := "some-key"
	resOne, err := prepareKey(key)
	if err != nil {
		t.Errorf("prepareKey() error = %v", err)
	}
	resTwo, err := prepareKey(key)
	if err != nil {
		t.Errorf("prepareKey() error = %v", err)
	}
	if hex.EncodeToString(resOne) != hex.EncodeToString(resTwo) {
		t.Errorf("prepareKey() not same resOne = %v, resTwo = %v", resOne, resTwo)
	}
}

func Test_decryptStringStability(t *testing.T) {
	key, err := prepareKey("some-key")
	if err != nil {
		t.Fatalf("prepareKey() error = %v", err)
	}
	data := []byte("some sensitive data")
	cryptedData, err := encrypt(key, data)
	if err != nil {
		t.Errorf("encrypt() error = %v", err)
	}
	resOne, err := decrypt(cryptedData, hex.EncodeToString(key))
	if err != nil {
		t.Errorf("encrypt() error = %v", err)
	}
	resTwo, err := decrypt(cryptedData, hex.EncodeToString(key))
	if err != nil {
		t.Errorf("encrypt() error = %v", err)
	}
	if resOne != resTwo {
		t.Errorf("encrypt() not same resOne = %v, resTwo = %v", resOne, resTwo)
	}
}
