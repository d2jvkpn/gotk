package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	// _ "github.com/lib/pq"
)

func MigratePsqlFromDir(dsn, migrations string) (err error) {
	var (
		db     *sql.DB
		driver database.Driver
		migr   *migrate.Migrate
	)

	if db, err = sql.Open("postgres", dsn); err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	if driver, err = postgres.WithInstance(db, &postgres.Config{}); err != nil {
		return fmt.Errorf("postgres.WithInstance: %w", err)
	}

	migr, err = migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrations),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %w", err)
	}

	// or m.Step(2) if you want to explicitly set the number of migrations to run
	if err = migr.Up(); err != nil {
		if err.Error() != "no change" {
			return fmt.Errorf("Migrate.Up: %w", err)
		}
	}

	e1, e2 := migr.Close()
	if err = errors.Join(e1, e2); err != nil {
		return fmt.Errorf("Migrate.Close: %w", err)
	}

	return nil
}

func MigratePsqlFromFs(dsn string, src fs.FS) (err error) {
	var (
		driver source.Driver
		migr   *migrate.Migrate
	)

	if driver, err = iofs.New(src, "/"); err != nil {
		return err
	}

	if migr, err = migrate.NewWithSourceInstance("iofs", driver, dsn); err != nil {
		return fmt.Errorf("migrate.NewWithSourceInstance: %w", err)
	}

	// or m.Step(2) if you want to explicitly set the number of migrations to run
	if err = migr.Up(); err != nil {
		if err.Error() != "no change" {
			return fmt.Errorf("Migrate.Up: %w", err)
		}
	}

	e1, e2 := migr.Close()
	if err = errors.Join(e1, e2); err != nil {
		return fmt.Errorf("Migrate.Close: %w", err)
	}

	return nil
}
