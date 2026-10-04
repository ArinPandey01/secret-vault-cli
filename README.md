# Secret Vault CLI

Secret Vault CLI is a small Go program for storing short-lived secrets on your computer. It encrypts each secret, saves the encrypted data in a local SQLite database, and gives you an ID and key to read it later.

## Why I made it

This is one of my small projects to brush up on Go.

## Getting started

You need Go 1.27.1 or newer. From the project directory, run:

```sh
go run . create
```

Enter a time to live, such as `30s`, `5m`, or `1h`. Press Enter to use the default of five minutes. Then enter your secret. The command prints its ID and a hexadecimal key. Keep both: the key is not saved in the database.

To read the secret before it expires, run:

```sh
go run . read
```

Enter the ID and key printed by `create` when prompted. Run `go run . help` to see the available commands.

## How it works

Each secret gets a new random key and is encrypted with AES-GCM. The SQLite database stores the ciphertext, a nonce, and an expiration time. Reading a secret requires its ID and key; expired secrets cannot be read. Expired records are not automatically deleted.

The database file, `secret-vault.db`, is created in the directory where you run the command. This is a learning project and has not been security audited.

## Privacy and access

The program runs locally and has no server component. It does not send your secrets to a remote server. Someone who gets a secret's ID and key **and can access the database file** can decrypt that secret, so keep the key private.

The CLI refuses to read expired secrets, but it does not delete their encrypted records. Expiry is therefore a rule enforced by the CLI, not a way to erase the stored data.
