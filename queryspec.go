package datarepo

// QuerySpec specifies a cross-DB query.
//
// Id1 and IdN should be mutually exclusive, such
// that if IdN is nil or len==0 then Id1 governs.
// Similarly(?), fetch "ALL" could be indicated
// by Id1 set to -1.
//
// This does not handle a list of values to be equal
// to, except for multiple IDs (passed in as [IdN]). 
//
// DBOp should probably be defined in package [dsmnd].
// .
type QuerySpec struct {
	DBOp   string
	Table  string
	Id1    int
	IdN    []int
	IsAnd  bool // Irrelevant when only one condition
	// Conditions is AND-or-OR of (e.g.)
	// { FieldName, "<= Zork" } .
	// For SELECT, is a map[string]nils 
	Conditions map[string]string
}

