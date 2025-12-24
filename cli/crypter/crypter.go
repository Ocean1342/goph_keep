package crypter

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"goph_keeper/cli/entity"
	"io"
)

//TODO: возможные ошибки

type Crypter struct {
	storage *storage
}

func New() (*Crypter, error) {
	return &Crypter{
		storage: &storage{storage: make(map[string][]byte)},
	}, nil
}

// Encode - кодирует. Возвращает строки, чтбы удобно было работать в бд.
func (c Crypter) Encode(_ context.Context, data []byte, user *entity.User) (string, error) {
	keyBytes, err := c.getCryptKey(user)
	if err != nil {
		return "", fmt.Errorf("could not get crypt key: %w", err)
	}
	res, err := encrypt(keyBytes, data)
	if err != nil {
		return "", fmt.Errorf("error encrypting data: %v", err)
	}
	return res, nil
}

func (c Crypter) Decode(_ context.Context, data []byte, user *entity.User) (string, error) {
	keyBytes, err := c.getCryptKey(user)
	if err != nil {
		return "", fmt.Errorf("could not get crypt key: %w", err)
	}
	res, err := decrypt(hex.EncodeToString(data), hex.EncodeToString(keyBytes))
	if err != nil {
		return "", fmt.Errorf("error decrypting data: %v", err)
	}
	return res, nil
}

func (c Crypter) getCryptKey(user *entity.User) ([]byte, error) {
	keyBytes, err := c.storage.Get(user.Login)
	if err != nil {
		if errors.Is(ErrKeyNotFound, err) {
			//generate
			keyBytes, err = prepareKey(user.Phrase)
			if err != nil {
				return nil, fmt.Errorf("error preparing crypto key: %v", err)
			}
			//store
			c.storage.Set(user.Login, keyBytes)
		} else {
			return nil, fmt.Errorf("error getting key: %v", err)
		}
	}
	if keyBytes == nil {
		return nil, fmt.Errorf("key does not exist")
	}
	return keyBytes, nil
}

// prepareKey - подготавливает 32 байтный слайс для дальнейшего преобразования
func prepareKey(phrase string) ([]byte, error) {
	if phrase == "" {
		return nil, fmt.Errorf("phrase is empty")
	}
	hash := sha256.Sum256([]byte(phrase))
	return hash[:], nil
}

func encrypt(key []byte, data []byte) (string, error) {
	if len(key) != 32 {
		return "", fmt.Errorf("error encrypting data: key length must be 32. given len: %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %v", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %v", err)
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error creating nonce: %v", err)
	}
	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(encryptedString string, keyString string) (string, error) {
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %v", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %v", err)
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting data: %v", err)
	}
	return fmt.Sprintf("%s", plaintext), nil
}
