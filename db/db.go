package db

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	id         int
	ciphertext []byte
	nonce      []byte
	expireAt   time.Time
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "secret-vault.db")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS secrets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ciphertext BLOB NOT NULL,
		nonce BLOB NOT NULL,
		expire_at DATETIME NOT NULL
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func InsertSecret(db *sql.DB, ciphertext []byte, nonce []byte, expireAt time.Time) (int64, error) {
	insertQuery := `
	INSERT INTO secrets (ciphertext, nonce, expire_at)
	VALUES (?, ?, ?);
	`
	result, err := db.Exec(insertQuery, ciphertext, nonce, expireAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetSecret(db *sql.DB, id int64) (*DB, error) {
	selectQuery := `
	SELECT id, ciphertext, nonce, expire_at
	FROM secrets
	WHERE id = ?;
	`
	row := db.QueryRow(selectQuery, id)

	var secret DB
	err := row.Scan(&secret.id, &secret.ciphertext, &secret.nonce, &secret.expireAt)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}
