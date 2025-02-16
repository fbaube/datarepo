package datarepo

import DRM "github.com/fbaube/datarepo/rowmodels"

// Caller methods are a temp list.
type Caller[T DRM.RowModel] interface {
     BeginImmed() error 
     DoSelectByIdGeneric(SimpleRepo, int, T) (bool, error) 
     DoInsertGeneric(SimpleRepo, T) (int, error) 
}

