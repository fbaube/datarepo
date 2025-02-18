package datarepo

// UniquerySpec specifies a query that keys on equality for a
// UNIQUE column. Generics instantiation requires the Go type
// of the keyed-on column `keyT`, e.g. `UniquerySpec[string]`.
//
// In principle the keyed-on column can be any UNIQUE column,
// but the most common case is 
//  - column name `ID` (or `{TBL}_ID`) 
//  - Go  type `int`
//  - SQL type `INT` (or, in SQLite, `INTEGER`)
// 
// For DBOp we can/could have
//  - sql INSERT / crud CREATE / http POST / "Add"  (pass in a record, get an ID)
//  - sql SELECT / crud RETRIV / http GET  / "Get"   (return a record)
//  - sql UPDATE / crud UPDATE / http PUT  / "Mod"  (pass in a record) 
//  - sql DELETE / crud DELETE / http DEL. / "Del"  (beware: FKey issues) 
//
// An error should be returned only if the SQL was rejected. If
// the unique row is not found, return a nil object and no error.
//
// This is meant to be passed to a query composer (not "builder")
// that is specific to a DB. Which means, for now, SQLite.
// 
// DBOp should probably be defined in package dsmnd.
// .
type UniquerySpec[keyT comparable] struct {
	DBOp   string
	Table  string
	Field  string
	Value  keyT
}

