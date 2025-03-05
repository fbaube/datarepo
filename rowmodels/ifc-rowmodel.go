package rowmodels

// RowModel is implemented by ptrs-to-structs 
// in package [github.com/fbaube/m5db] 
type RowModel interface {
     TableDetails() TableDetails 
     ColumnNamesCsv(bool) string
     // Note that ColumnPtrsMethod is a method,
     // and/but while there is also a func with 
     // signature ColumnPtrsFunc(RowModel, bool) 
     ColumnPtrsMethod(bool) []any
     // We should constrain it to be a ptr: 
     // *T
     // The generic REST interface takes a ptr to a
     // buffer, and each such ptr is a RowModel, so
     // ideally this interface would include a way
     // to use a table name to get a typed memory 
     // allocator, such as: FuncNew() RowModel
     // But instead, file tabledetails.go has this: 
     // type NewInstanceFunc func() RowModel
}

/*
Here is another type of constraint:
// EndpointOf says T is a ptr and it implements interface BinaryUnmarshaler.
type EndpointOf[T any] interface {
	*T
	encoding.BinaryUnmarshaler
}
*/