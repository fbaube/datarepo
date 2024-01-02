package sqlite

import (
	"database/sql"
	D "github.com/fbaube/dsmnd"
	// "github.com/jmoiron/sqlx"
	// "github.com/pocketbase/dbx"
	// _ "github.com/mattn/go-sqlite3"
)

// RepoEntity provides lifecycle operations for databases.
//
// type RepoEntity interface {

/* RepoEntity ifc API
Handle() *sqlx.DB // (noun) the handle to the DB
Type() string     // "sqlite" (equiv.to "sqlite3")
Path() string     // file/URL (or dir/URL, if uses multiple files)
IsURL() bool
IsSingleFile() bool // true for sqlite
*/

// Handle (noun, not verb) returns the dbx handle to the DB.
func (pSR *SqliteRepo) Handle() *sql.DB {
	return pSR.DB
}

// Type returns "sqlite" (equiv.to "sqlite3").
func (pSR *SqliteRepo) Type() D.DB_type {
	return D.DB_SQLite // "sqlite"
}

// Path returns the filepath (for a multi-file
// DB type it could return dir/URL).
func (pSR *SqliteRepo) Path() string { // file/URL (or
	return pSR.filepath
}

// IsURL returns false for sqlite.
func (pSR *SqliteRepo) IsURL() bool {
	return false
}

// IsURL returns true for sqlite.
func (pSR *SqliteRepo) IsSingleFile() bool {
	return true
}
