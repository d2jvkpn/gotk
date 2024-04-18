package orm

import (
	// "fmt"
	// "strings"

	"github.com/jackc/pgx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	// "github.com/lib/pq"
)

/*
dsn format:
- {USERANME}:{PASSWORD}@tcp({IP})/{DATABASE}?sslmode=disable
- "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai"
*/
func PgConnect(dsn string, debugMode bool) (db *gorm.DB, err error) {
	conf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}
	// dsn: "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai"
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

// !!! works as expected only version(gorm.io/driver/postgres) <= v1.4.5
//
//	and errors.Is(e, gorm.ErrDuplicatedKey) doesn't work as expected
//
// PgUniqueViolation
func IsPgDuplicatedKey(err error) bool {
	// fmt.Printf("~~~ err: %[1]s\n    %#+[1]v, type: %[1]T\n", err)
	e, ok := err.(*pgx.PgError)
	// fmt.Printf("~~~ e: %v, ok: %t\n", e, ok)
	return ok && e.Code == "23505"

	/*
		switch err.(type) {
		case *pgx.PgError:
			return true
		default:
			return false
		}
	*/

	// return strings.Contains(err.Error(), "(SQLSTATE 23505)")
}

func IsPgRecordNotFound(err error) bool {
	return err.Error() == "record not found"
}
