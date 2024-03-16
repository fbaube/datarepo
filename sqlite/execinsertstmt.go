package sqlite

import (
	"fmt"
	"database/sql"
	L "github.com/fbaube/mlog" // Brings in global var L
)

// ExecInsertStmt executes a simple (not prepared) SQL statement
// and returns the ID (primary key) of the added item.
// 
// Conceptually, it should use Exec() and not Query().
//
// Notes on this:
//  - [Exec](..) returns ([Result], error)
//  - ONLY Exec returns Result
//  - ONLY Result has LastInsertId() (is an int64) 
//  - Query(..) returns (*Rows, error)
//  - QueryRow(..) returns (*Row)
//
// API for [Row] (and [Rows]):
//  - func (r *Row/s) Scan(dest ...any) error
//  - Scan(..) can get the inserted-ID of INSERT...RETURNING 
//  - func (r *Row/s) Err() error 
//  - Use Err() to check for query errors without calling Scan(..)
//  - Err() returns any error hit when running the query
//  - If not nil, the error is also returned from Scan(..)
// .
func (pSR *SqliteRepo) ExecInsertStmt(stmt string) (int, error) {

	var res sql.Result
	var id int64
	var e error 
	// ==========
	// EXEC
	fmt.Fprintf(pSR.w, "=== ExecInsStmt.SQL ===\n%s\n", stmt)
	res, e = pSR.Exec(stmt)
	if e != nil {
	     	L.L.Error("Exec.Ins failed: %w", e)
		return -1, e
		} 
	id, e = res.LastInsertId()
	if e != nil {
	     	L.L.Error("Exec.Ins.LastInsertId failed: %w", e)
		return -1, fmt.Errorf("ExecInsStmt.LastInsertId: %w", e)
		} 
	fmt.Fprintf(pSR.w, "=== ExecInsStmt.RETval.id: %d ===\n", id)
	/*
	// ===================
	// or try QUERY + Scan
	// ===================
	var rowQ *sql.Row
	var idQ int 
	var errQ error 
	// ==========
	// QUERY
	L.L.Info("exec.L440+: Trying QUERY STMT: " + stmt)
	rowQ = pSR.QueryRow(stmt)
	errQ = rowQ.Scan(&idQ)
	// ==========
	*/

	return int(id), nil
}

