package datarepo

// RestQuery is meant to be an abstract query that
// can be interpreted by different SQL DBs as needed.
// DBOp is TBD but so far we have "get", which is
// actually an HTTP REST op, d'oh.
//
// Id1 and IdN should be mutually exclusive, such
// that if IdN is nil or len==0 then Id1 governs.
// Similarly(?), fetch "ALL" could be indicated
// by Id1 set to -1.
//
// It should probably include a datum to field 
// the query results.
//
// DBOp should probably be define din package dsmnd.
// .
type RestQuery struct {
	DBOp   string
	Table  string
	Id1    int
	IdN    []int
	// Conditions is AND of (e.g.)
	// { Name, { "!=", "Zork" } }.
	// For SELECT, is a map[string]nils 
	Conditions map[string][]string
}

