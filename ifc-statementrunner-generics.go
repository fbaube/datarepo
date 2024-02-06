package datarepo

import(
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// StatementRunner_generics is DB-specific 
// and implemented by *[sqlite.SqliteRepo] 
type StatementRunner_generics [T DRM.RowModel] interface {
	ExecInsertStmt(string) (int, error)
	ExecSelectOneStmt(string) (T, error) 
}

