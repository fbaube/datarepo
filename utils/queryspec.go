package utils

// QuerySpec is the basic strings that get plugged into a query:
// table name, column names, "where" clause. It is meant to be
// passed to query composers (not "builders") that are specific
// to various DBs. Which means, for now, SQLite.
//
// A query composer might also be a query executor, which would
// mean distinguiishing at the API level among return values of
// int64, []int64, Row, and []Row.
//
// There is redundancy built in, to help ensure against errors in
// usage: the DB op, the presence or absence of a WHERE clause,
// the number of IDs passed (0, 1, N).
// .
type QuerySpec struct {
	DbOp
	// Table must not be empty; if it is treated as
	// case-insensitive then no validity checking
	// is done.
	Table string
	// Fields must not be empty, for consistency
	// in technique re. using ptrs everywhere in
	// order to enable generics and Scan(..).
	Fields string
	// ID can safely be ignored if
	// (IDs != nil) && (ID == 0 || ID == -1)
	ID int64
	// IDs != nil is an error if ID is valid.
	IDs   []int64
	Where string
}
