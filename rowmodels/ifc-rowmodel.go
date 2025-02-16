package rowmodels

// RowModel is implemented by ptrs-to-structs 
// in package [github.com/fbaube/m5db] 
type RowModel interface {
     TableDetails() TableDetails 
     ColumnNamesCsv(bool) string 
     ColumnPtrsMethod(bool) []any
     // The generic REST interface takes a ptr to 
     // a buffer, and each such ptr is a RowModel, 
     // so this interface must include a way to use 
     // a table name to get a typed memory allocator. 
     // FuncNew() RowModel
     // When generic, include *T
     // *T
}

/*
Here is another type of constraint:

// EndpointOf applies structural type constraints to T and 
// makes sure it implements the unmarshaler interface.
type EndpointOf[T any] interface {
	*T
	encoding.BinaryUnmarshaler
}
*/