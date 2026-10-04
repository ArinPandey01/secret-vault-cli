package main

import (
	"fmt"
	"os"

	"github.com/ArinPandey01/secret-vault-cli/create"
	"github.com/ArinPandey01/secret-vault-cli/db"
	"github.com/ArinPandey01/secret-vault-cli/retrieve"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`Secret Vault stores short-lived, encrypted secrets in a local database.

Usage:
  secret-vault-cli <command>

Commands:
  create  Enter a secret and time to live (default: 5m). Prints its ID and key.
  read    Retrieve a secret.
  help    Show this help message.

Keep the ID and key from create: both are needed to decrypt the secret.
`)
}

func run() error {
	if len(os.Args) != 2 {
		printHelp()
		return nil
	}

	command := os.Args[1]
	if command != "create" && command != "read" {
		printHelp()
		return nil
	}

	conn, err := db.InitDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	switch command {
	case "create":
		key, id, err := create.Create(conn)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d\n", id)
		fmt.Printf("Key: %x\n", key)
	case "read":
		secret, err := retrieve.RetrieveSecret(conn)
		if err != nil {
			return err
		}
		fmt.Printf("Secret: %s\n", secret)
	}

	return nil
}
