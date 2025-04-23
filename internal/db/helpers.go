package db

import (
	"github.com/caarlos0/env/v6"
)

func getConfig() (config, error) {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}

func (db *DataBase) Close() {
	db.conn.Close()
}
