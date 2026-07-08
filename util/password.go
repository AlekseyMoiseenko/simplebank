package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Config struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

// DefaultConfig returns default config
func DefaultConfig() Argon2Config {
	// owasp Cheat Sheet minimum configuration
	return Argon2Config{
		Time:    2,
		Memory:  19 * 1024,
		Threads: 1,
		KeyLen:  32,
		SaltLen: 16,
	}
}

// returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	cfg := DefaultConfig()
	salt := make([]byte, cfg.SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, cfg.Time, cfg.Memory, cfg.Threads, cfg.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, cfg.Memory, cfg.Time, cfg.Threads, b64Salt, b64Hash)
	return encoded, nil
}

// checks if the provided password is correct
func CheckPassword(password, encodedHash string) error {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return fmt.Errorf("invalid hash format: expected 6 parts, got %d", len(parts))
	}
	if parts[1] != "argon2id" {
		return fmt.Errorf("unsupported algorithm: %s (only argon2id)", parts[1])
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return errors.New("invalid hash format")
	}
	if version != argon2.Version {
		return fmt.Errorf("incompatible version: expected %d, got %d", argon2.Version, version)
	}

	var memory, time uint32
	var threads uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return errors.New("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("failed to decode salt: %w", err)
	}
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(storedHash)))

	if subtle.ConstantTimeCompare(computedHash, storedHash) != 1 {
		return fmt.Errorf("invalid password")
	}
	return nil
}
