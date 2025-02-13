package datarepo

import DRU "github.com/fbaube/datarepo/utils"

// QueryRunner runs queries!
//
// https://go.dev/wiki/SQLInterface#executing-queries
//  - ExecContext gets NO rows:
//    result, err := db.ExecContext(ctx,
//   "INSERT INTO users (name, age) VALUES ($1, $2)", "gopher", 27)
//  - type [sql.Result] interface {
//	LastInsertId() (int64, error) ; 
//	RowsAffected() (int64, error) }
//   - QueryRowContext gets 1 row
//   - QueryContext gets N rows
// .
type QueryRunner interface {
     	// RunQuery0 is Exec() is get NO rows
	RunQuery0(*DRU.QuerySpec) (any, error) 
	// RunQuery1 is QueryRowContext() is get ONE row (like: by_ID)
	RunQuery1(*DRU.QuerySpec) (any, error)
	// RunQueryN is QueryContext is get N rows 
	RunQueryN(*DRU.QuerySpec) ([]any, error) 
}
