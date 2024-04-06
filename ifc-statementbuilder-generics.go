package datarepo

import(
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// StatementBuilder_generics is DB-specific and
// implemented by *[sqlite.SqliteRepo] 
type StatementBuilder_generics [T DRM.RowModel] interface {
	// BuildQueryStmt(*DRU.QuerySpec) (string, error)
	NewCreateTableStmt(T) (string, error)
	NewSelectByIdStmt(T, int) (string, error)
	// NewInsertStmt(T) (string, error)
}
