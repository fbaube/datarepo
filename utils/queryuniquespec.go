package utils

// QueryUniqueSpec is a convenience for executing a WHERE .. =" 
// for a value of a UNIQUE column, and simply return all columns
// ("SELECT *"). Specified: table name, column name, column value.
// 
// A caller using this spec might also include a (new!) value for 
// the row, i.e. an UPDATE, or if not, then it is a simple SELECT.
//
// An error should be returned only if the SQL was rejected. If
// the unique row is not found, return a nil object and no error. 
// 
// This is meant to be passed to a query composer (not "builder")
// that is specific to a DB. Which means, for now, SQLite.
// .
type QueryUniqueSpec struct {
	// DbOp
	// Table must not be empty; if it is treated as
	// case-insensitive then no validity checking
	// is done.
	Table string
	// If field is empty, it is UNIQUE ID. 
	Field string
	// Value has to be flexible.
	Value any
}
