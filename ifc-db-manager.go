package datarepo

import(
	"github.com/fbaube/datarepo/sqlite"
)

var DB_Manager DBManager

func init() {
     DB_Manager = sqlite.SQLite_DB_Manager
     }

// DBManager has methods to create and open databases.
//
// NOTE: The recommended action is to call OpenAtPath,
// which then selects one of the other two.
//
// It defines (for now: only) SQLite implementations.
// .
type DBManager interface {
	// OpenAtPath(string) (SimpleRepo, error) // recommended 
	OpenAtPath(string) (*sqlite.SqliteRepo, error) // recommended 
	// NewAtPath(string) (SimpleRepo, error)
	NewAtPath(string) (*sqlite.SqliteRepo, error)
	// OpenExistingAtPath(string) (SimpleRepo, error)
	OpenExistingAtPath(string) (*sqlite.SqliteRepo, error)
}
