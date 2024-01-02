package datarepo

import (
	"database/sql"
	D "github.com/fbaube/dsmnd"
	// "github.com/jmoiron/sqlx"
	// "github.com/pocketbase/dbx"
	// _ "github.com/mattn/go-sqlite3"
)

// Entity provides realization-related
// operations for databases.
//
// It defines methods on a (probably empty)
// singleton that selects (for now: only)
// SQLite implementations of this interface.
// .
type Entity interface {
	Handle() *sql.DB    // (noun) the handle to the DB
	Type() D.DB_type    // DB_SQLite ("sqlite", equiv.to "sqlite3")
	Path() string       // file/URL (or dir/URL, if uses multiple files)
	IsURL() bool        // false for a local SQLite file
	IsSingleFile() bool // true for SQLite
}
