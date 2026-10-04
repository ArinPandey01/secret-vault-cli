package retrieve

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ArinPandey01/secret-vault-cli/db"
)

func decryptSecret(secret *db.DB, key []byte) ([]byte, error) {
	if !time.Now().Before(secret.ExpireAt) {
		return nil, fmt.Errorf("secret has expired")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decryptedSecret, err := gcm.Open(nil, secret.Nonce, secret.Ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decryptedSecret, nil
}

func scanInput(scanner *bufio.Scanner, prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(scanner.Text()), nil
}

func RetrieveSecret(conn *sql.DB) ([]byte, error) {
	scanner := bufio.NewScanner(os.Stdin)
	idInput, err := scanInput(scanner, "ID: ")
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(idInput, 10, 64)
	if err != nil || id <= 0 {
		return nil, fmt.Errorf("invalid ID: enter a positive number")
	}

	keyInput, err := scanInput(scanner, "Key (hex): ")
	if err != nil {
		return nil, err
	}
	key, err := hex.DecodeString(keyInput)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("invalid key: enter the 64-character hex key shown by create")
	}

	result, err := db.GetSecret(conn, id)
	if err != nil {
		return nil, err
	}
	decryptedSecret, err := decryptSecret(result, key)
	if err != nil {
		return nil, err
	}

	return decryptedSecret, nil
}
