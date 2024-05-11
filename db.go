package main

import (
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type EnergyLevel struct {
	ID        int
	Value     string
	CreatedAt time.Time
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

	return &DB{
		db: sdb,
	}, nil
}

func (d *DB) Insert(val string) error {
	slog.Debug("writting value", "val", val)
	_, err := d.db.Exec("INSERT INTO levels (value) VALUES (?)", val)
	return err
}

func (d *DB) GetAll() ([]EnergyLevel, error) {
	rows, err := d.db.Query("SELECT id, value, created_at FROM levels")
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
