package datarepo

import DRU "github.com/fbaube/datarepo/utils"

// QueryRunner runs queries!
//
// https://github.com/golang/go/wiki/SQLInterface
//   - ExecContext is used when no rows are returned ("0")
//   - QueryContext is used for retrieving rows ("N")
//   - QueryRowContext is used where only a single row is expected ("1")
//
// .
type QueryRunner interface {
	RunQuery0(*DRU.QuerySpec) (any, error)   // ie. Exec()
	RunQuery1(*DRU.QuerySpec) (any, error)   // One row, like by_ID
	RunQueryN(*DRU.QuerySpec) ([]any, error) // Multiple rows
}
