package orm

import (
	"errors"
	// "fmt"
	"database/sql"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	MYSQL_Datetime = "2006-01-02 15:04:05"
)

/*
MySQL initialize

	dsn: {USERANME}:{PASSWORD}@tcp({IP})/{DATABASE}?charset=utf8mb4&parseTime=True&loc=Local
*/
func MySQLConnect(vp *viper.Viper, debug bool) (db *gorm.DB, err error) {
	var (
		conf  *gorm.Config
		sqlDB *sql.DB
	)

	conf = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	if db, err = gorm.Open(mysql.Open(vp.GetString("dsn")), conf); err != nil {
		return nil, err
	}
	if debug {
		db = db.Debug()
	}

	// Get generic database object sql.DB to use its functions
	if sqlDB, err = db.DB(); err != nil {
		return nil, err
	}

	if v := vp.GetInt("max_idle_conns"); v > 0 {
		sqlDB.SetMaxIdleConns(v)
	}
	if v := vp.GetInt("max_open_conns"); v > 0 {
		sqlDB.SetMaxOpenConns(v)
	}
	if v := vp.GetDuration("conn_max_idle_time"); v > 0 {
		sqlDB.SetConnMaxLifetime(v)
	}
	if v := vp.GetDuration("conn_max_lifetime"); v > 0 {
		sqlDB.SetConnMaxLifetime(v)
	}

	return db, err
}

// errors
func MySQLIsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// MySQL Unique Violation
func MySQLUniqueViolation(err error) bool {
	sqlErr, ok := err.(*gomysql.MySQLError)
	return ok && sqlErr.Number == uint16(1062)
}

func MySQLDatetime(t ...time.Time) string {
	if len(t) == 0 {
		t = []time.Time{time.Now()}
	}

	return t[0].Format(MYSQL_Datetime)
}

func mysql_e01_json(tx *gorm.DB) {
	var (
		title string
	)

	tx = tx.Where("JSON_EXTRACT(data, '$.title') LIKE ?", "%"+title+"%")
	// SQL -- JSON_SET(translation, '$.Japanese', JSON_OBJECT('category', '分类', 'tpye', '类型'))
	tx = tx.Where("JSON_CONTAINS(field, ?, '$')", "xxx")
	// SQL -- JSON_CONTAINS('[1,2,3,4,5]','1','$')
	// SQL -- SELECT JSON_OVERLAPS("[1,3,5,7]", "[2,5,7]")
	// SQL -- SELECT JSON_SEARCH('["abc", [{"k": "10"}, "def"], {"x":"abc"}, {"y":"bcd"}]', 'one', 'abc');
	//     "$[0]"
	// SQL -- SELECT JSON_SEARCH(@j, 'all', 'abc');
	//     ["$[0]", "$[2].x"]
}
