package utils

// TODO These should be done with sprintf(..).

/*
// Dooo NOT include ID here !
// TODO Use FieldPtrs here !!
var Q_create = "INSERT into %s (%s) values (?, ?, ?);"

var Q_getByID = "SELECT ID, %s from %s where ID = ? limit 1;"

var Q_getCount = "SELECT COUNT(*) from %s;"

var Q_getByIDs = "SELECT ID, %s from %s where ID in (%s) order by id;"

var Q_getByIDsWithBio = "SELECT ID from %s where Bio != ''; -- cmt"
*/

var QueriesByOp = map[DbOp]string{
	// Odd ops
	OpCount:       "SELECT COUNT(*) from %s;",
	OpCreateTable: "!TBS!",
	// One item
	OpAdd:     "INSERT into %s (%s) values (?, ?, ?);",
	OpGetByID: "SELECT ID, %s from %s where ID = ? limit 1;",

	OpGetByIDs:    "SELECT ID, %s from %s where ID in (%s) order by id;",
	OpGetIDsWhere: "SELECT ID from %s where %s;",
}
