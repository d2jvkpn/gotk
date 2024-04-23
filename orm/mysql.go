package orm

import (
	"errors"
	// "fmt"

	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/*
MySQL initialize

	dsn format: {USERANME}:{PASSWORD}@tcp({IP})/{DATABASE}?charset=utf8mb4&parseTime=True&loc=Local
*/
func GormMySQConnectL(dsn string, debug bool) (db *gorm.DB, err error) {
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
