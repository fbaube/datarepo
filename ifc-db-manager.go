package datarepo

// DB_Manager is a global, maybe for SQLite, maybe for
// ebberyting. It probably needs some sort of mutex. 
var DB_Manager DBManager

// func init() { DB_Manager = DRS.SQLite_DB_Manager }

// DBManager has methods to create, open, and configure databases.
//
// NOTE: The recommended action is to call OpenAtPath,
// which then selects one of the other two.
// .
type DBManager interface {
	OpenAtPath(string) (SimpleRepo, error) 
	 NewAtPath(string) (SimpleRepo, error)
	OpenExistingAtPath(string) (SimpleRepo, error)
	// InitznPragmas is assumed to be multiline 
	InitznPragmas() string
	// ReadonlyPragma is assumed to be a single pragma 
	// that returns some sort of status message
	ReadonlyPragma() string 
}
