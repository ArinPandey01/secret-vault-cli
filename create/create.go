package create

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"time"
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

func Create() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter time to live(Default=5m)")

	var ttl time.Duration
	var err error

	ok := scanner.Scan()
	if !ok {
		fmt.Printf("Error has occured: %v", scanner.Err())
		return "", scanner.Err()
	}

	input := scanner.Text()

	if len(input) == 0 {
		ttl = 5 * time.Minute
	} else {
		ttl, err = time.ParseDuration(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("For ttl use one of these format: <n>s, <n>m, <n>h")
			return "", err
		}
	}

	fmt.Println(ttl)

	fmt.Print("Enter your secret: ")

	ok = scanner.Scan()
	if !ok {
		fmt.Printf("Error has occured: %v", scanner.Err())
		return "", scanner.Err()
	}

	input = scanner.Text()

	key, err := randomBytes(32)
	if err != nil {
		fmt.Printf("Error has occured: %v", err)
		return "", err
	}

	ciphertext, nonce, err := encrypt(key, input)
	if err != nil {
		fmt.Printf("Error has occured: %v", err)
		return "", err
	}

	fmt.Printf("Ciphertext: %x\n", ciphertext)
	fmt.Printf("Nonce: %x\n", nonce)

	return input, nil // later this becomes ID/key output
}