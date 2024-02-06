package datarepo

import(
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// StatementRunner is DB-specific and
// implemented by *[sqlite.SqliteRepo] 
type StatementRunner interface {
	ExecInsertStmt(string) (int, error)
	ExecSelectOneStmt(stmt string) (DRM.RowModel, error) 
}

