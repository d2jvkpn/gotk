package orm

import (
	// "fmt"
	// "strings"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

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
