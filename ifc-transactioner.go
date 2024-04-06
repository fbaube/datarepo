package datarepo

import "database/sql"

// Transactioner methods come from the Go stdlib, but are modified 
// to work with a single possibly-active transaction. If it is active, 
// the trio of Exec/Query/QueryRow usw it, and if not, the calls go
// straight to the sql.DB .
//
// There are a couple of problems with this, mainly when there are
// multiple threads accessing. The fix would be that a call to Begin
// would return the shared sql.DB but a unique sql.Tx .
type Transactioner interface {
     GetTx() *sql.Tx
     IsInTx() bool 
     // func (db *DB) Begin() (*Tx, error)
     Begin() error // for re-entrancy return (Transactioner, error)
     Commit() error
     Rollback() error
     Exec(query string, args ...any) (sql.Result, error)
     Query(query string, args ...any) (*sql.Rows, error)
     QueryRow(query string, args ...any) *sql.Row
}
