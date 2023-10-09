package orm

import (
	// "fmt"
	"database/sql"

	"gorm.io/gorm"
)

func CloseDB(gdb *gorm.DB) (err error) {
	var db *sql.DB

	if gdb == nil {
		return nil
	}

	if db, err = gdb.DB(); err != nil {
		return err
	}

	return db.Close()
}
