package sqlite

/*

import (
	"fmt"
	S "strings"
	"database/sql"
	//D "github.com/fbaube/dsmnd"
	DRP "github.com/fbaube/datarepo"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

/*
func (p *SqliteRepo) RunQuery0(*DRU.QuerySpec) (any, error) { // ie. Exec()
	fmt.Fprintf(p.w, "NOT IMPL'D: RunQuery0 \n")
	return nil, nil
}
* /

// RunUniquerySpec uses strings (table, field, value) to
// retrieve a single row based on the value of a UNIQUE column.
// If Field is named "ID" it gets enhanced to "{table}_ID".
// 
// Value is a string no matter what comparable type keyT is,
// but that's OK because we are only using integers and strings.
// .
// func RunUniquerySpec[keyT comparable](pSR *SqliteRepo, pUS *DRP.UniquerySpec[keyT]) (any, error) { // One row, like by_ID
func RunUniquerySpec(pSR *SqliteRepo, pUS *DRP.UniquerySpec) (any, error) { 
	// fmt.Fprintf(pSR.LogWriter(), "NOT IMPL'D: RunQuery1 \n")
	var pTD *DRM.TableDetails
        pTD = GetTableDetailsByCode(pUS.Table)
	// We might have to modfiy the Value string we were passed.
	// But for now let's now worry about it. 
	/* switch pUS.FVtypes {
	case SQLITE_ERR:
	case SQLITE_INTEGER: // 64-bit signed integer
	case SQLITE_FLOAT: // 64-bit IEEE FP number
	case SQLITE_TEXT: // string; incl JSON 
	case SQLITE_BLOB: // incl JSONB 
	case SQLITE_NULL: 
	case SQLYT_DATETIME: 
	} * /
	var sQuery string
	if S.EqualFold("ID", pUS.Field) {
	   	pUS.Field = pUS.Table + "_" + pUS.Field 
	   	println("Fixed ID to: ", pUS.Field)
	}
	sQuery = "SELECT * from " + pUS.Table + " WHERE " +
	     	pUS.Field + " = " + pUS.Value + ";"
	w := pSR.LogWriter()
        fmt.Fprintf(w, "== %s.RunUniquerySpec.SQL ===\n%s\n",
                pUS.Table, sQuery)
	var row *sql.Row
        var e error 
        // ==========
        // QUERY (not EXEC)
        // func (db *DB) QueryRow(query string, args ...any) *Row 
        // L.L.Info("RunUniquerySpec.SQL: " + theStmt)
        row = pSR.Handle().QueryRow(sQuery)
        // ==========
        // What if there is no row in the result, and .Scan() can't
        // scan a value. What then? The error constant sql.ErrNoRows
        // is returned by QueryRow() when the result is empty. 
        // This needs to be handled as a special case in most cases.
        // You should only see this error when you’re using QueryRow().
        // If you see this error elsewhere, you’re doing something wrong.
	var colPtrs []any
        colPtrs = pTD.ColumnPtrsFunc(pTD.BlankInstance, true)
        e = row.Scan(colPtrs...)
        switch e {
          case sql.ErrNoRows:
               return false, nil
          case nil:
               return true, nil
          default:
                return false, fmt.Errorf("RunUniquerySpec" +
                       "(%s:%d) failed: %w", pUS.Value, pTD.StorName, e)
        }
	return nil, nil
}

func RunQuerySpec(pSR *SqliteRepo, pQS *DRP.QuerySpec) ([]any, error) { // Multiple rows
	fmt.Fprintf(pSR.LogWriter(), "NOT IMPL'D: RunQueryN \n")
	return nil, nil
}

*/

