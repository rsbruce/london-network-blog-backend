package resourcedata

import (
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DbConn *sqlx.DB
}
