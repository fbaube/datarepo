// ==> REPLACE SqliteRepo w SimpleRepo !!!!

package sqlite

import (
	"fmt"
	"io"
	DRP "github.com/fbaube/datarepo"
	D "github.com/fbaube/dsmnd"
	S "strings"
	DRM "github.com/fbaube/datarepo/rowmodels"
	"database/sql"
	L "github.com/fbaube/mlog" // Bring in global var L
)

// DoInsertGeneric takes a generic RowModel and "does" a simple
// (not prepared) SQL INSERT statement, returning the inserted
// item's ID (i.e. primary key). A statement is created, then
// envalued & executed using regular parameter substitution
// with Exec(..).
// 
// API notes on this:
//  - [Exec](..) returns ([Result], error)
//  - ONLY Exec returns a Result
//  - ONLY a Result has LastInsertId() (is an int64) 
// // - QueryRow(..) returns (*Row)
// // - Query(..) returns (*Rows, error)
//
// API for [Row] (and [Rows]):
//  - func (r *Row/s) Scan(dest ...any) error
//  - Scan(..) can get the inserted-ID of INSERT...RETURNING 
//  - func (r *Row/s) Err() error 
//  - Use Err() to check for query errors without calling Scan(..)
//  - Err() returns any error hit when running the query
//  - If not nil, the error is also returned from Scan(..)
// .
// func DoInsertGeneric[T DRM.RowModel](pSR *SqliteRepo, pRM T) (int, error) {
func DoInsertGeneric[T DRM.RowModel](pSR DRP.SimpleRepo, pRM T) (int, error) {
	// [This is a bit OBS:] So, we have to 
	//  1) Check the struct-type of the RowModeler-instance
	//  2) Fetch the FieldPtrs
	//  3) Write out the fields to the stmt
	//  4) Write out the values to the stmt
	//  5) Add RETURNING
	//  6) (the stmt's user) Use that returned that ID
	
	var colPtrs []any
	var cp any 
	var iCol int 
	var pTD DRM.TableDetails
	// TMP var now = SU.Now() // time.Now().UTC().Format(time.RFC3339)

	// Here the generic argument is
	// resolved to a specific type 
	pTD = pRM.TableDetails()
	// false says do not include primary key 
	colPtrs = pTD.ColumnPtrsFunc(pRM, false)
	
	/* TMP - why was this disabled ? 
	pRM.T_Cre = now
	pRM.T_Imp = now
	pRM.T_Edt = now
	*/

	var w io.Writer
	// w = pSR.w
	w = pSR.LogWriter()

	// fmt.Fprintf(os.Stderr, "LENS: ColSpex<%d> ColPtrs<%d> \n",
	//	len(pTD.ColumnSpecs), len(colPtrs))
	// Add some log info 
	fmt.Fprintf(w, "=== %s.DoInsGenc.ColSpex ===\n", pTD.StorName)

	// ===========================
	//  1) Write table name and 
	//      all column names (CSV)
	// ===========================
	// colPtrs should NOT include
	// the primary key, D.SFT_PRKEY
	var sqlBldr S.Builder // simpl,e StringBuilder
	sqlBldr.WriteString("INSERT INTO ")
	sqlBldr.WriteString(pTD.TableSummary.StorName)
	sqlBldr.WriteString("(")
	sqlBldr.WriteString(pTD.ColumnNamesCSV)
	sqlBldr.WriteString(") VALUES(")

	// ===============================
	//  2) Add list of '$'-numbered 
	//     parameters (Postgres-style) 
	// ===============================
	// var sft D.SemanticFieldType
	var sn string
	var dt D.SemanticFieldType
	for iCol, _ = range colPtrs {
	    sqlBldr.WriteString(fmt.Sprintf("$%d, ", iCol))
	    // Add some log info 
	    sn = pTD.ColumnSpecs[iCol].StorName
	    dt = D.SemanticFieldType(pTD.ColumnSpecs[iCol].Datatype)
	    fmt.Fprintf(w, "[%d] %s / %s / %T \n", iCol, sn, dt, cp)
	}
	/* OBSOLETE
	// Stuff that was used when composing
	// the SQL stmt using Sprintf 
	for iCol, cp = range colPtrs {
	    switch cp.(type) {
	    	   case *string:
		   	var pS *string
			var sS string
			pS = cp.(*string)
			sS = *pS
		   	sqlBldr.WriteString(fmt.Sprintf("'%s', ", sS))
	    	   case *FU.AbsFilePath:
		   	sqlBldr.WriteString(fmt.Sprintf("'%s', ", AFPval(cp)))
	    	   case *SU.MarkupType:
		   	sqlBldr.WriteString(fmt.Sprintf("'%s', ", MTval(cp)))
	    	   case *CT.Raw:
		   	sqlBldr.WriteString(fmt.Sprintf("'%s', ", CTRval(cp)))
		   case *int:
		   	var pI *int
			pI = cp.(*int)
		   	sqlBldr.WriteString(fmt.Sprintf("%d, ", *pI))
	    }
	}
	*/
	// ===============================
	//  3) Trim off the last ", " and
	//     add RETURNING, which asks 
	//     for the added index
	// ===============================
	var theSQL string 
	theSQL = S.TrimSuffix(sqlBldr.String(), ", ")
	theSQL += ") RETURNING IDX_" + pTD.StorName + ";"
	fmt.Fprintf(w, "=== %s.DoInsGenc.SQL ===\n%s\n",
		pTD.StorName, theSQL)
		
	// return theSQL, nil
	// }

// (COPIED from...) (pSR *SqliteRepo) ExecInsertStmt(stmt string) (int, error) {

	var res sql.Result
	var id int64
	var e error 
	// ================================
	//  4) Call Exec(..) on the stmt, 
	//     with all columns (?pointers)
	// ================================
	res, e = pSR.Handle().Exec(theSQL, colPtrs...)
	if e != nil {
	     	L.L.Error("Exec.Ins failed: %w", e)
		return -1, e
		} 
	id, e = res.LastInsertId()
	if e != nil {
	     	L.L.Error("Exec.Ins.LastInsertId failed: %w", e)
		return -1, fmt.Errorf("ExecInsStmt.LastInsertId: %w", e)
		} 
	fmt.Fprintf(w, "=>> ExecInsStmt.RET.id: %d ===\n", id)
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

