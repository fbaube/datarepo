package datarepo

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

