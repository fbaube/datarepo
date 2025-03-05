package sqlite

// check ~/go/src/github.com/simukti/sqldb-logger

import(
	"fmt"
	"errors"
	"io"
	"database/sql"
	S "strings"
	D "github.com/fbaube/dsmnd"
	DRP "github.com/fbaube/datarepo"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// EngineUnique acts on a single DB record, based on the value of 
// of a column that is specified as UNIQUE (for example, a row ID). 
// One of four basic actions is performed (listed as SQL/CRUD/HTTP):
//  - SELECT / Retrieve / GET  (returns 0/1 + buffer) 
//  - INSERT / Create / POST   (returns new-ID + nil) 
//  - UPDATE / Update / PUT    (returns 0/1 + nil) 
//  - DELETE / Delete / DELETE (returns 0/1 + nil) 
// 
// It takes four input arguments:
//  - The DB operation, one of the four listed above; the dbOp is
//    specified by the first letter (only!) of the argument dbOp
//  - The name of the DB table (case-insensitive)
//  - A "whereSpec" of column name and column value (not used for INSERT; as
//    a convenience, if the name is "ID", it is modified to be {tableNeme}_ID 
//  - A pointer to a buffer, used for input (if INSERT or UPDATE) or output
//    (if SELECT); for INSERT, the  buffer is left unmodified and the ID of
//    the new record is returned in the second return value (the int) 
//  - Note that the buffer pointer is an interface, implemented by a pointer 
//    to a struct; this may be usable by GO generics
//
// It returns two values:
//  - An error: if it is non-nil, the other return value is invalid
//  - An int that is (if INSERT) the newly-added row ID (else) 0 or 1 to
//    indicate how many records were (i.e. whether a record was) affected.
//
// NOTE: When using whereSpec, if a record is not found, this is indicated
// by the second return value (the int), NOT by the first return value 
// (the error, which is reserved for when the DB rejects the SQL.
//
// NOTE: In an UPDATE, if the whereSpec does not refer to the ID, and the
// ID of the input record does not match the ID of the record found by the
// DB, the function panics. So, for UPDATE, just match on the ID. 
//
// NOTE: Also implement COUNT(*) ?
// . 
func (pSR *SqliteRepo) EngineUnique(dbOp string, tableName string, pWS *DRP.UniquerySpec, RM DRM.RowModel) (error, int) {

     var pTD *DRM.TableDetails
     var pRM DRM.RowModel
     var w = pSR.LogWriter()
     var sSQL string

     // Declare some vars used by multiple ops 
     // Table's column-pointers function 
     var idxdCPF, CPF []any  // with ID; no ID
     // Table's comma-separated column names 
     var idxdCSV, CSV string // with ID; no ID
     // Table's '$'-numbered parameters (Postgres-style)
     var idxdPlcNrs, PlcNrs string // with ID; no ID 

     // Fetch the table's details and fill in the vars 
     pTD = GetTableDetailsByCode(tableName)
     if pTD == nil {
     	// FIXME err msgs 
     	s := "NO TblDtls FOR: " + tableName
     	println(s)
	return errors.New(s), 0
     }
     // For convenience, callers can use "ID", and we fix it 
     if pWS != nil && S.EqualFold("id", pWS.Field) {
     	pWS.Field = pTD.PKname
     }
     pRM = pTD.NewInstance() // output buffer 
     CSV = pTD.ColumnNamesCSV // no ID column 
     CPF = pTD.ColumnPtrsFunc(RM, false) // no ID column 
     idxdCSV = pTD.PKname + ", " + CSV // with ID column 
     idxdCPF = pTD.ColumnPtrsFunc(RM, true) // with ID column
     var i int 
     for i = range len(CPF) {
	 PlcNrs += fmt.Sprintf("$%d, ", i+1) }
     for i = range len(idxdCPF) {
	 idxdPlcNrs += fmt.Sprintf("$%d, ", i+1) }
     PlcNrs     = S.TrimSuffix(PlcNrs, ", ")
     idxdPlcNrs = S.TrimSuffix(idxdPlcNrs, ", ")

     // Log info about the columns 
     // writeFieldDebugInfo(w, pTD)

     // switch dbOp {
     // We only use the first letter of the 
     // DB op, so callers can be creative :-P 
     println("DB OP IS: " + S.ToUpper(dbOp[0:1]))
     switch S.ToUpper(dbOp[0:1]) { 

     	// ======================================================
	case "A", "C", "I", "N":
	// Add, Create, Insert, New 
	// https://www.sqlite.org/lang_insert.html
	// INSERT INTO tblNm (fld1, fld2) VALUES(val1, val2);
	// INSERT INTO tblNm (fld1, fld2) VALUES($1,$2); + any...
     	// ======================================================
	if pWS != nil {
	   return errors.New("EngineUnique: INSERT: unwanted WHERE"), 0 
	}
	// Write table name and all column names (as CSV).
	// Do NOT include the primary key, D.SFT_PRKEY 
	sSQL = "INSERT INTO " +
	        pTD.TableSummary.StorName +
		"(" + CSV + ") " +
		"VALUES(" + PlcNrs + ") " +
		"RETURNING " + pTD.PKname + ";"
	fmt.Fprintf(w, "INSERT.sql: " + sSQL + "\n")
	
	// It is now ready for Exec()
	var theRes sql.Result
	var newID  int64
	var e      error
	// Call Exec(..) on the stmt, with all column ptrs
	theRes, e = pSR.Handle().Exec(sSQL, CPF...)
	if e != nil {
		fmt.Fprintf(w, "INSERT.exec: failed: %s", e)
		return fmt.Errorf("engineunique: insert: exec: %w", e), -1
	}
	newID, e = theRes.LastInsertId()
	if e != nil {
		fmt.Fprintf(w, "INSERT.lastinsertId: failed: %s", e)
		return fmt.Errorf("engineunique: insert: lastinsertId: %w",e),-1
	}
	fmt.Fprintf(w, "INSERT: OK: LastInsertID: %d \n", newID)
	return nil, int(newID)
	
     	// =======================================
        case "F", "G", "L", "R", "S":
        // Fetch, Get, List, Retrieve, Select 
        // https://www.sqlite.org/lang_select.html
	// SELECT fld1, fld2 FROM tblNm WHERE expr
	// https://www.sqlite.org/syntax/expr.html 
     	// =======================================
	if pWS == nil {
	   return errors.New("EngineUnique: SELECT: missing WHERE"), 0
	}
	sSQL =  "SELECT " + idxdCSV +
		" FROM "  + pTD.TableSummary.StorName +
		" WHERE " + pWS.Field + " = " + pWS.Value + ";"
		
	// TODO: QueryRow(..)
	row := pSR.Handle().QueryRow(sSQL)
	// ==========
	// What if there is no row in the result, and .Scan() can't
	// scan a value. What then? The error constant sql.ErrNoRows
	// is returned by QueryRow() when the result is empty.
	// This needs to be handled as a special case in most cases.
	// You should only see this error when you're using QueryRow().
	// If you see this error elsewhere, you're doing something wrong.
	
	var colPtrs []any
	var e error
	// idxdCPF = = pTD.ColumnPtrsFunc(RM, true)
	colPtrs = pRM.ColumnPtrsMethod(true) 
	e = row.Scan(colPtrs...)
	switch e {
	  case sql.ErrNoRows:
	       return nil, 0 // false, nil
	  case nil:
	       return nil, 1 // true, nil 
	  default:
		println("SQL ERROR: (" + e.Error() + ") SQL: " + sSQL)
		return fmt.Errorf("EngineUnique(get) " +
		       "(%s=%s) failed: %w", pWS.Field, pWS.Value, e), 0
	}
	panic("Oops, fallthru in SELECT")
	
     	// ================================================
	case "M", "U": // Modify, Update
	// https://www.sqlite.org/lang_update.html
	// Obnoxious syntax: 
	// UPDATE tblNm SET fld1=val1,fld2=val2 WHERE expr: 
	// (or..) SET fld1=$1, fld2=$2 WHERE expr; + any...
	// https://www.sqlite.org/syntax/expr.html 
     	// ================================================
     	if pWS == nil {
	   return errors.New("EngineUnique: UPDATE: missing WHERE"), 0 
     	   }
// ====================================================
	// For UPDATE (only), we have to generate an SQL
	// string that involves all columns (except the ID). 
	// Write assignment pairts as CSV: f1 = $1, f2 = $2, ...
	// We do NOT include the primary key, D.SFT_PRKEY


	sSQL = "UPDATE " +
	        pTD.TableSummary.StorName + " SET " + 
		"(" + CSV + ") " +
		"VALUES(" + PlcNrs + ") " +
		"RETURNING " + pTD.PKname + ";"
	fmt.Fprintf(w, "INSERT.sql: " + sSQL + "\n")
	
	// It is now ready for Exec()
	var theRes sql.Result
	var newID  int64
	var e      error
	// Call Exec(..) on the stmt, with all column ptrs
	theRes, e = pSR.Handle().Exec(sSQL, CPF...)
	if e != nil {
		fmt.Fprintf(w, "INSERT.exec: failed: %s", e)
		return fmt.Errorf("engineunique: insert: exec: %w", e), -1
	}
	newID, e = theRes.LastInsertId()
	if e != nil {
		fmt.Fprintf(w, "INSERT.lastinsertId: failed: %s", e)
		return fmt.Errorf("engineunique: insert: lastinsertId: %w",e),-1
	}
	fmt.Fprintf(w, "INSERT: OK: LastInsertID: %d \n", newID)
	return nil, int(newID)

// ====================================================

     	// =======================================
	case "D": // Delete, Discard, Drop
	// https://www.sqlite.org/lang_delete.html
	// DELETE FROM tblNm WHERE expr RET'G expr
     	// =======================================
	if pWS == nil {
     	   return errors.New("EngineUnique: SELECT: missing WHERE"), 0 
     	}
     // default:
     }
     return errors.New("engineunique: bad dbOp: " + dbOp), -1
}

func writeFieldDebugInfo(w io.Writer, pTD *DRM.TableDetails) {
     	// TODO: Check the correctness of this! It seemed to overrun with "true"
	CPF := pTD.ColumnPtrsFunc(pTD.BlankInstance, false) // true) // with ID column
	for iCol, cp := range CPF {
	    sn := pTD.ColumnSpecs[iCol].StorName
	    dt := D.SemanticFieldType(pTD.ColumnSpecs[iCol].Datatype)
	    fmt.Fprintf(w, "Column [%d] %s / %s / %T \n", iCol, sn, dt, cp)
	}
}
