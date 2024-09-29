package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/coral/midiname/ai"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func New(path string) (*DB, error) {
	os.Remove(path)

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	newDB := &DB{db: db}
	err = newDB.SetupTables()
	if err != nil {
		return nil, err
	}

	return newDB, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Add(resp ai.Response, filename string) error {
	hash := sha256.Sum256([]byte(filename))
	hashString := hex.EncodeToString(hash[:])

	genresJSON, err := json.Marshal(resp.Genres)
	if err != nil {
		return err
	}
	_, err = d.db.Exec("INSERT INTO midi (title, artist, genres, comments, decade,  confident, filename, hash) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", resp.Title, resp.Artist, genresJSON, resp.Comments, resp.Decade, resp.Confident, filename, hashString)
	return err
}

func (d *DB) SetupTables() error {
	sql, err := os.ReadFile("init.sql")
	if err != nil {
		return err
	}
	_, err = d.db.Exec(string(sql))
	return err
}
