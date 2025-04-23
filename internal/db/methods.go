package db

import (
	"github.com/jmoiron/sqlx"
)

func New() (*DataBase, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	d, err := sqlx.Open("postgres", cfg.Dsn)
	if err != nil {
		return nil, err
	}

	err = d.Ping()
	if err != nil {
		return nil, err
	}

	db := new(DataBase)
	db.conn = d
	db.cfg = cfg

	return db, nil
}
