package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 32
	keySize    = 32
	iterations = 100000
)

// Encryptor provides methods for encrypting and decrypting sensitive data
type Encryptor struct {
	masterKey []byte
}

// NewEncryptor creates a new encryptor with the provided master key
func NewEncryptor(masterKey string) *Encryptor {
	// Derive a key from the master key string
	hash := sha256.Sum256([]byte(masterKey))
	return &Encryptor{
		masterKey: hash[:],
	}
}

// GenerateSalt generates a random salt for key derivation
func GenerateSalt() (string, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// EncryptWithSalt encrypts data using AES-GCM with a salt-derived key
func (e *Encryptor) EncryptWithSalt(plaintext []byte, salt string) ([]byte, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	// Derive key from master key and salt
	key := pbkdf2.Key(e.masterKey, saltBytes, iterations, keySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptWithSalt decrypts data using AES-GCM with a salt-derived key
func (e *Encryptor) DecryptWithSalt(ciphertext []byte, salt string) ([]byte, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	// Derive key from master key and salt
	key := pbkdf2.Key(e.masterKey, saltBytes, iterations, keySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns base64 encoded ciphertext
func (e *Encryptor) EncryptString(plaintext string, salt string) (string, error) {
	encrypted, err := e.EncryptWithSalt([]byte(plaintext), salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptString decrypts a base64 encoded ciphertext and returns the plaintext string
func (e *Encryptor) DecryptString(ciphertext string, salt string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	decrypted, err := e.DecryptWithSalt(data, salt)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// EncryptJSON encrypts JSON data
func (e *Encryptor) EncryptJSON(data []byte, salt string) ([]byte, error) {
	return e.EncryptWithSalt(data, salt)
}

// DecryptJSON decrypts JSON data
func (e *Encryptor) DecryptJSON(ciphertext []byte, salt string) ([]byte, error) {
	return e.DecryptWithSalt(ciphertext, salt)
}