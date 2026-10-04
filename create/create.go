package create

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ArinPandey01/secret-vault-cli/db"
)

func randomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func encrypt(key []byte, secret string) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce, err := randomBytes(gcm.NonceSize())
	if err != nil {
		return nil, nil, err
	}

	plaintext := []byte(secret)
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

func Create(conn *sql.DB) ([]byte, int64, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Fprintln(os.Stderr, "Enter time to live (default: 5m):")

	var ttl time.Duration
	var err error

	ok := scanner.Scan()
	if !ok {
		if err := scanner.Err(); err != nil {
			return nil, 0, err
		}
		return nil, 0, io.EOF
	}

	input := scanner.Text()

	if len(input) == 0 {
		ttl = 5 * time.Minute
	} else {
		ttl, err = time.ParseDuration(strings.TrimSpace(input))
		if err != nil {
			fmt.Fprintln(os.Stderr, "For TTL, use a duration such as 30s, 5m, or 1h")
			return nil, 0, err
		}
	}

	fmt.Fprint(os.Stderr, "Enter your secret: ")

	ok = scanner.Scan()
	if !ok {
		if err := scanner.Err(); err != nil {
			return nil, 0, err
		}
		return nil, 0, io.EOF
	}

	input = scanner.Text()

	key, err := randomBytes(32)
	if err != nil {
		return nil, 0, err
	}

	ciphertext, nonce, err := encrypt(key, input)
	if err != nil {
		return nil, 0, err
	}

	id, err := db.InsertSecret(conn, ciphertext, nonce, time.Now().Add(ttl))
	if err != nil {
		return nil, 0, err
	}

	return key, id, nil
}
