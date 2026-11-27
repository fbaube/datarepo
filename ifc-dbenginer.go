package datarepo

import(
    "uuid"
	DRM "github.com/fbaube/datarepo/rowmodels"
)
	
// DBEnginer is TBS.
type DBEnginer interface {
     // EngineUnique works on/with a single DB record. 
     	EngineUnique(dbOp string, u *uuid.UUID, RM DRM.RowModel) error
}
