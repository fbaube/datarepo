package datarepo

// DBGetting has methods to create / find / open databases.
//
// NOTE: The recommended action is to call OpenAtPath,
// which then selects one of the other two.
//
// string retvals are fullpaths. 
//
// It defines (for now: only) SQLite implementations.
// .
type DBGetting interface {
	OpenAtPath(string) (string, error) // recommended 
	NewAtPath(string) (string, error)
	OpenExistingAtPath(string) (string, error)
}
