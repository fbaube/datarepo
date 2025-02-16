package datarepo

import (
	"database/sql"
	D "github.com/fbaube/dsmnd"
)

// Entity provides operations for database entities (i.e. instances).
type Entity interface {
     	// Handle (noun) is the handle to the DB.
	Handle() *sql.DB
	// Type has value [dsmnd.DB_SQLite] ("sqlite", equiv.to "sqlite3").
	Type() D.DB_type
	// Path is the file/URL (or dir/URL, if uses multiple files) to the DB.
	Path() string
	// IsURL is false for a local SQLite file.
	IsURL() bool
	// IsSingleFile is true for SQLite. 
	IsSingleFile() bool 
}
