package db

import (
	"database/sql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type DB struct {
	db *sql.DB
}

func NewDB(db_ *sql.DB) *DB {
	return &DB{
		db: db_,
	}
}

func (db *DB) Close() {
	err := db.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("could not close underlying DB")
	}
}

func debugLog(sql_ string, args []interface{}, msg string) {
	if log.Debug().Enabled() {
		la := zerolog.Arr()
		for _, arg := range args {
			la = la.Interface(arg)
		}
		log.Debug().CallerSkipFrame(1).Str("sql", sql_).Array("args", la).Msg(msg)
	}
}
