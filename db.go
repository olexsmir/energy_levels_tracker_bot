package main

import (
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type EnergyLevel struct {
	ID        int       `json:"id"`
	Value     string    `json:"value"`
	Hour      int       `json:"hour"`
	CreatedAt time.Time `json:"created_at"`
}

type DB struct {
	db *sql.DB
}

func NewDB(dbstr string) (*DB, error) {
	sdb, err := sql.Open("libsql", dbstr)
	if err != nil {
		return nil, err
	}

	if err = sdb.Ping(); err != nil {
		return nil, err
	}

	return &DB{db: sdb}, nil
}

func (d *DB) Insert(val string, hour int, createdAt time.Time) error {
	slog.Debug("writing value", "val", val)
	_, err := d.db.Exec(
		"INSERT INTO levels (value, hour, created_at) VALUES (?, ?, ?)",
		val,
		hour,
		createdAt,
	)
	return err
}

func (d *DB) GetAll() ([]EnergyLevel, error) {
	rows, err := d.db.Query("SELECT id, value, hour, created_at FROM levels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []EnergyLevel
	for rows.Next() {
		var el EnergyLevel
		if err := rows.Scan(&el.ID, &el.Value, &el.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, el)
	}

	return res, rows.Err()
}
