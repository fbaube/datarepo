package rowmodels

/*
import(
	DRU "github.com/fbaube/datarepo/utils"
)
*/

// RowModeler is implemented by ptrs to structs in 
// package [github.com/fbaube/datarepo/rowmodels] 
type RowModeler interface {
     TableDetails() TableDetails 
     // ColumnNamesCSV() string 
     ColumnPtrs(bool) []any
     // When generic, include *T
     // *T
}

/*
Here is another type of constraint:

// EndpointOf applies structural type constraints to T and makes sure
// it implements the unmarshaler interface.
type EndpointOf[T any] interface {
	*T
	encoding.BinaryUnmarshaler
}

*/