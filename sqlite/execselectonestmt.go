package sqlite

import (
	"fmt"
	"database/sql"
	L "github.com/fbaube/mlog" // Brings in global var L
	RM "github.com/fbaube/datarepo/rowmodels"
)

// ExecInsertStmt executes a simple (not prepared) SQL statement.
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
// func (pSR *SqliteRepo) ExecSelectOneStmt[T RowModeler](stmt string) (T, error) {
func ExecSelectOneStmt[T RM.RowModeler](pSR *SqliteRepo, stmt string) (T, error) {

	var row *sql.Row
	var e error 
	// ==========
	// QUERY (not EXEC)
	// func (db *DB) QueryRow(query string, args ...any) *Row 
	L.L.Info("Trying QUERY SELECT ONE: " + stmt)
	// res, e = pSR.Exec(stmt)
	row = pSR.QueryRow(stmt)
	// ==========
	// What if there was no user with id = 1? Then there would be 
	// no row in the result, and .Scan() would not scan a value 
	// into name. What happens then?
	// Go defines a special error constant, called sql.ErrNoRows, 
	// which is returned from QueryRow() when the result is empty. 
	// This needs to be handled as a special case in most cases.
	// An empty result might not be an error in application code, 
	// and if you don’t check whether an error is this special 
	// constant, you’ll cause app-level errors you didn’t expect.
	// You might ask: Why is an empty result set an error ? There's
	// nothing erroneous about an empty set. The reason is that the 
	// method QueryRow() needs to use this special-case in order to
	// let the caller distinguish whether QueryRow() in fact found 
	// a row; without it, Scan() wouldn’t do anything and you might 
	// not realize that your variable didn’t get any value from the 
	// database after all.
	// You should only see this error when you’re using QueryRow().
	// If you see this error elsewhere, you’re doing something wrong.

	// func (r *Row) Scan(dest ...any) error
	var colPtrs []any
	var anInstance T
	// var paI *T
	// paI = &anInstance
	colPtrs = /*paI*/anInstance.ColumnPtrs(true)
	// scanErr := row.Scan(colPtrs)

	// ===

        if e = row.Scan(colPtrs); e != nil {
            return anInstance, fmt.Errorf("ExecSelectOneStmt(scan:%s): %w", stmt, e)
        }
	if e = row.Err(); e != nil { // rows
           return anInstance, fmt.Errorf("ExecSelectOneStmt(%s):err:  %w", stmt, e)
	   }

	return anInstance, nil
}

