package sqlite

import(
	DRU "github.com/fbaube/datarepo/utils"
)

// RowModeler is implemented by ptrs to structs in 
// package [github.com/fbaube/datarepo/rowmodels] 
type RowModeler interface {
     TableDetails() DRU.TableDetails 
     ColumnNamesCSV() string 
     ColumnPtrs() []any
}

