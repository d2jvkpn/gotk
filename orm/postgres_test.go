package orm

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx"
)

func TestPgError(t *testing.T) {
	var err error = new(pgx.PgError)

	fmt.Printf("~~~ err: %#+[1]v, type: %[1]T\n", err)

	e, ok := err.(*pgx.PgError)
	fmt.Printf("~~~ ok: %t, e: %#+v\n", ok, e)
}
