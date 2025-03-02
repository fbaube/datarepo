package datarepo

import(
	DRM "github.com/fbaube/datarepo/rowmodels"
)
	
// DBEnginer is TBS.
type DBEnginer interface {
     	EngineUnique(dbOp string, tableName string, whereSpec *UniquerySpec, RM DRM.RowModel) (error, int)
}
