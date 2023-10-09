package orm

import (
	"fmt"
	"testing"

	"github.com/jackc/pgconn"
)

func TestPgError(t *testing.T) {
	var err error = new(pgconn.PgError)

	fmt.Printf("~~~ err: %#+[1]v, type: %[1]T\n", err)

	e, ok := err.(*pgconn.PgError)
	fmt.Printf("~~~ ok: %t, e: %#+v\n", ok, e)
}
