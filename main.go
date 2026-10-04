package main

import (
	"fmt"
	"os"

	"github.com/ArinPandey01/secret-vault-cli/create"
	"github.com/ArinPandey01/secret-vault-cli/db"
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

	switch os.Args[1] {
	case "create":
		conn, err := db.InitDB()
		if err != nil {
			return err
		}
		defer conn.Close()

		key, _, err := create.Create(conn)
		if err != nil {
			return err
		}
		fmt.Printf("Key: %x\n", key)
	case "read":
	case "help":
		printHelp()
	default:
		printHelp()
	}

	return nil
}
