package datarepo

import(
	DRS "github.com/fbaube/datarepo/sqlite"
)

// DB_Manager is a global, maybe for SQLite, maybe for ebberyting, 
var DB_Manager DBManager

func init() {
     DB_Manager = DRS.SQLite_DB_Manager
     }

// DBManager has methods to create, open, and configure databases.
//
// NOTE: The recommended action is to call OpenAtPath,
// which then selects one of the other two.
// .
type DBManager interface {
	// OpenAtPath should be
	// OpenAtPath(string) (SimpleRepo, error) // recommended 
	OpenAtPath(string) (*DRS.SqliteRepo, error) // recommended 
	// NewAtPath should be
	// NewAtPath(string) (SimpleRepo, error)
	NewAtPath(string) (*DRS.SqliteRepo, error)
	// OpenExistingAtPath should be
	// OpenExistingAtPath(string) (SimpleRepo, error)
	OpenExistingAtPath(string) (*DRS.SqliteRepo, error)
	// InitznPragmas is assumed to be multiline 
	InitznPragmas() string
	// ReadonlyPragma is assumed to be a single pragma that
	// returns some sort of status message
	ReadonlyPragma() string 
}
