package datarepo

import (
	"fmt"
	S "strings"
	"database/sql"
	L "github.com/fbaube/mlog" 
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// DoSelectByIdGeneric returns true/false indicating whether the ID
// was found.
//
// In order to process the ID, the func needs a generic argument that can
// be used in the body of the func. A [RowModel] is a pointer - having 
// write access - so we kill two birds with one stone by passing in a 
// RowModel pointing to a destination buffer to Scan the DB fetch into.
// . 
// func DoSelectByIdGeneric[T DRM.RowModel](pSR *SqliteRepo, anID int, pDest T) (bool, error) {
func DoSelectByIdGeneric[T DRM.RowModel](pSR SimpleRepo, anID int, pDest T) (bool, error) {

	// So, we have to
	//  1) Check the struct-type of the RowModeler-instance
	//  2) Fetch the FieldPtrs
	//  3) Write out the fields to the stmt
	//  4) Write out the values to the stmt
	//  5) Add RETURNING
	//  6) (the stmt's user) Use that returned that ID

	var pTD DRM.TableDetails
	pTD = pDest.TableDetails()
	var sb S.Builder
	sb.WriteString("SELECT " + pTD.PKname + ", ")
	sb.WriteString(pTD.ColumnNamesCSV + " FROM ")
	sb.WriteString(pTD.TableSummary.StorName)
	sb.WriteString(" WHERE " + pTD.PKname + " = ")
	// We don't worry about SQL injection here (tho we should?) 
	sb.WriteString(fmt.Sprintf("%d;", anID))
	
	var theStmt string 
	theStmt = sb.String()
	w := pSR.LogWriter()
	fmt.Fprintf(w, "== %s.DoSelectByIdGeneric.SQL ===\n%s\n",
		pTD.StorName, theStmt)

	var row *sql.Row
	var e error 
	// ==========
	// QUERY (not EXEC)
	// func (db *DB) QueryRow(query string, args ...any) *Row 
	L.L.Info("DoSelectByIdGeneric.SQL: " + theStmt)
	row = pSR.Handle().QueryRow(theStmt)
	// ==========
	// What if there is no row in the result, and .Scan() can't
	// scan a value. What then? The error constant sql.ErrNoRows
	// is returned by QueryRow() when the result is empty. 
	// This needs to be handled as a special case in most cases.
	// You should only see this error when you’re using QueryRow().
	// If you see this error elsewhere, you’re doing something wrong.

	var colPtrs []any
	colPtrs = pDest.ColumnPtrsMethod(true)
	e = row.Scan(colPtrs...)
	switch e {
	  case sql.ErrNoRows:
	       return false, nil
	  case nil:
	       return true, nil
	  default:
		return false, fmt.Errorf("DoSelectByIdGeneric" +
		       "(%s:%d) failed: %w", anID, pTD.StorName, e)
	}
}

