package datarepo

// DBGetting has methods to create / find /
// open databases.
//
// The default action is to call OpenAtPath,
// which then selects one of the other two.
//
// It defines methods on a (probably empty)
// singleton that selects (for now: only)
// SQLite implementations of this interface.
//
// .
type DBGetting interface {
	OpenAtPath(string) (string, error)
	NewAtPath(string) (string, error)
	OpenExistingAtPath(string) (string, error)
}
