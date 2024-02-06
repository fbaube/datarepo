package datarepo

import(
	// DRU "github.com/fbaube/datarepo/utils"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// StatementBuilder is DB-specific and
// implemented by *[sqlite.SqliteRepo] 
type StatementBuilder interface {
	// BuildQueryStmt(*DRU.QuerySpec) (string, error)
	NewCreateTableStmt(*DRM.TableDetails) (string, error)
	NewSelectByIdStmt(*DRM.TableDetails, int) (string, error)
	NewInsertStmt(DRM.RowModel) (string, error)
}
