package infrastructure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io"
	mathrand "math/rand"
	"time"

	"github.com/yakob-abada/delfare/deamon-service/domain"
)

type EventFactory interface {
	CreateEvent() domain.Event
}

type SecurityEventFactory struct {
	encryptionKey string
}

func NewSecurityEventFactory(encryptionKey string) *SecurityEventFactory {
	return &SecurityEventFactory{
		encryptionKey: encryptionKey,
	}
}

func (f *SecurityEventFactory) CreateEvent() domain.Event {
	criticality := mathrand.Intn(10)
	message := fmt.Sprintf("Security alert %d", mathrand.Intn(1000))

	encryptedMessage, err := encrypt(message, f.encryptionKey) // Replace with a secure key
	if err != nil {
		panic(fmt.Sprintf("Failed to encrypt message: %v", err))
	}

	return domain.Event{
		RequestID:   uuid.New().String(),
		Criticality: criticality,
		Timestamp:   time.Now().Format(time.RFC3339),
		Message:     encryptedMessage,
	}
}

func encrypt(text string, key string) (string, error) {
	// Validate key length
	keyLength := len(key)
	if keyLength != 16 && keyLength != 24 && keyLength != 32 {
		return "", fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
