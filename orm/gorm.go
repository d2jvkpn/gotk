package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GormDBClose(gdb *gorm.DB) (err error) {
	var db *sql.DB

	if gdb == nil {
		return nil
	}

	if db, err = gdb.DB(); err != nil {
		return err
	}

	return db.Close()
}

/*
MySQL initialize

	dsn format: {USERANME}:{PASSWORD}@tcp({IP})/{DATABASE}?charset=utf8mb4&parseTime=True&loc=Local
*/
func GormMySQLConnect(dsn string, debug bool) (db *gorm.DB, err error) {
	conf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	if db, err = gorm.Open(mysql.Open(dsn), conf); err != nil {
		return nil, err
	}
	if debug {
		db = db.Debug()
	}
	/*
		sqlDB, err := db.DB() // Get generic database object sql.DB to use its functions
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	*/

	return db, err
}

// errors
func GormMySQLIsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func GormMySQLUniqueViolation(err error) bool {
	sqlErr, ok := err.(*gomysql.MySQLError)
	return ok && sqlErr.Number == uint16(1062)
}

/*
dsn format:
- postgres://{USERANME}:{PASSWORD}@tcp({IP})/{DATABASE}?sslmode=disable
- "host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai"
*/
func GormPgConnect(dsn string, debugMode bool) (db *gorm.DB, err error) {
	conf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		TranslateError: true,
	}
	if db, err = gorm.Open(postgres.Open(dsn), conf); err != nil {
		return nil, err
	}

	if debugMode {
		db = db.Debug()
	}
	/*
		sqlDB, err := db.DB() // Get generic database object sql.DB to use its functions
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	*/

	return db, err
}

func GormPgNotFound(err error) bool {
	return err.Error() == "record not found"
	// return errors.Is(err, gorm.ErrDuplicatedKey) // !!! this doesn't work
}

// !!! works as expected only version(gorm.io/driver/postgres) <= v1.4.5
//
//	and errors.Is(e, gorm.ErrDuplicatedKey) doesn't work as expected
func GormPgUniqueViolation(err error) bool {
	// fmt.Printf("~~~ err: %[1]s\n    %#+[1]v, type: %[1]T\n", err)
	// e, ok := err.(*pgx.PgError)
	// fmt.Printf("~~~ e: %v, ok: %t\n", e, ok)
	// return ok && e.Code == "23505"

	// println("~~~", err.Error())
	// return err.Error() == "duplicated key not allowed"
	return errors.Is(err, gorm.ErrDuplicatedKey)
	// return strings.Contains(err.Error(), "(SQLSTATE 23505)")
}

func GormFlip(tx *gorm.DB, pageSize, pageNo int) *gorm.DB {
	if pageNo <= 0 || pageSize <= 0 {
		return tx.Limit(0)
	}

	return tx.Limit(pageSize).Offset((pageNo - 1) * pageSize)
}

func GormMultiEquals[T comparable](tx *gorm.DB, field string, values []T) *gorm.DB {
	wheres, items := make([]string, 0), make([]any, 0)
	unique := make(map[T]bool, len(values))

	for _, k := range values {
		if unique[k] {
			continue
		}
		wheres = append(wheres, fmt.Sprintf("%s = ?", field))
		items = append(items, k)
		unique[k] = true
	}

	if len(wheres) == 0 {
		tx.Where("0 = 1")
	}

	return tx.Where(strings.Join(wheres, " OR "), items...)
}
