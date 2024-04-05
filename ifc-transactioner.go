package datarepo

import "database/sql"

// Transactioner methods come from the Go stdlib, with the exception that
// Begin() returns a new Transactioner that has an active Transaction in it. 
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
