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

// Set value for sDBOP, then check (non)presence of WhereSpec.

// BuildINSERT writes table name + column names + placeholders. 
// Do NOT include the primary key, D.SFT_PRKEY 
func BuildINSERT(pTD *DRM.TableDetails) string { 
	return  "INSERT INTO " + pTD.TableSummary.StorName +
		           "(" + pTD.CSVs.FieldNames   + ") " +
		     "VALUES(" + pTD.CSVs.PlaceNumbers + ") " +
		  "RETURNING " + pTD.PKname            + ";"
}

// WITH WHERE and withOUT WHERE 
func buildSELECT(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }
func buildUPDATE(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }
func buildDELETE(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }

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
//
// TODO: switch on dbOp to call a new mini func that assembles the SQL statement.
// 
// TODO: TD.pCSVs should also use DB.Prepare to gather sql.Stmt's
//
// TODO: Use Result.RowsAffected
// https://pkg.go.dev/database/sql#Result
// RowsAffected returns the number of rows affected 
// by an update, insert, or delete. Not every DB or 
// driver supports it. RowsAffected() (int64, error)
//
// NOTE: When writing the multi-row version of this,
// be sure to call Rows.Cloe()
// . 
func (pSR *SqliteRepo) EngineUnique(dbOp string, tableName string, pFV *DRP.FieldValuePair, pRM DRM.RowModel) (error, int) {

     var pTD   *DRM.TableDetails
     var pCSVs *DRM.ColumnStringsCSV
  // var RM     DRM.RowModel
     var sSQL  string
     var w = pSR.LogWriter()
     var e error 
     // Table's column-pointers funcs 
     var CPF, CPF_wID []any 

     // Fetch the table's details and fill in the vars 
     pTD = GetTableDetailsByCode(tableName)
     if pTD == nil {
     	// FIXME err msgs 
     	s := "NO TblDtls FOR: " + tableName
     	println(s)
	return errors.New(s), 0
     }
     pCSVs = pTD.CSVs
     if pCSVs == nil { panic("nil TableDetails ColumnStrings") }
     
     // For convenience, callers can use "ID", and we fix it 
     if pFV != nil && S.EqualFold("id", pFV.Field) {
     	pFV.Field = pTD.PKname
     }
     if pRM == nil { pRM = pTD.NewInstance() } // output buffer 
     CPF     = pTD.ColumnPtrsFunc(pRM, false) // no ID column 
     CPF_wID = pTD.ColumnPtrsFunc(pRM, true) // with ID column

     // Log info about the columns 
     // writeFieldDebugInfo(w, pTD)

     // switch dbOp {
     // We only use the first letter of the 
     // DB op, so callers can be creative :-P 
     // println("DB OP IS: " + S.ToUpper(dbOp[0:1]))

/*      FieldNames,     FieldNames_wID  string  // "F1, F2, F3" 
        PlaceNumbers,   PlaceNrs_wID    string  // "$1, $2, $3" 
        UpdateNames     string   // "F1 = $1, F2 = $2, F3 = $3"  */

     switch S.ToUpper(dbOp[0:1]) { 

     	// ======================================================
	case "A", "C", "I", "N":
	// Add, Create, Insert, New
	// Use RETURNING to get new ID. 
	// https://www.sqlite.org/lang_insert.html
	// INSERT INTO tblNm (fld1, fld2) VALUES(val1, val2);
	// INSERT INTO tblNm (fld1, fld2) VALUES($1,$2); + any...
	// FIELDS are FieldNames[_wID]. VALUES are PlaceNrs[_wID].
     	// ======================================================
	if pFV != nil {
	   return errors.New("EngineUnique: INSERT: unwanted WHERE"), 0 
	}
	// Do NOT include the primary key, D.SFT_PRKEY 
	sSQL = BuildINSERT(pTD)
	fmt.Fprintf(w, "INSERT.sql: " + sSQL + "\n")
	
	// It is now ready for Exec()
	var theRes sql.Result
	var newID  int64
	// Call Exec(..) on the stmt, with all column ptrs
	theRes, e = pSR.Handle().Exec(sSQL, CPF...)
	if e != nil {
		fmt.Fprintf(w, "engineunique.insert.exec: failed: %s", e)
		return fmt.Errorf("engineunique.insert.exec: %w", e), -1
	}
	newID, e = theRes.LastInsertId()
	if e != nil {
		fmt.Fprintf(w, "engineunique.insert.lastinsertId: failed: %s", e)
		return fmt.Errorf("engineunique.insert: lastinsertId: %w",e),-1
	}
	fmt.Fprintf(w, "INSERT: OK: LastInsertID: %d \n", newID)
	return nil, int(newID)
	
     	// =======================================
        case "F", "G", "L", "R", "S":
        // Fetch, Get, List, Retrieve, Select 
        // https://www.sqlite.org/lang_select.html
	// SELECT fld1, fld2 FROM tblNm WHERE expr
	// https://www.sqlite.org/syntax/expr.html
	// FIELDS are FieldNames_wID. 
     	// =======================================
	if pFV == nil {
	   return errors.New("engineunique.select: missing WHERE"), 0
	}
	// TODO This should use a '$'-placeholder ?? 
	sSQL =  "SELECT " + pCSVs.FieldNames_wID +
		" FROM "  + pTD.TableSummary.StorName +
		" WHERE " + pFV.Field + " = " + pFV.Value + ";"

	// TODO: QueryRow(..)
	row := pSR.Handle().QueryRow(sSQL)
	// ---------------------------------------------------------
	// What if there is no row in the result, and .Scan() can't
	// scan a value. What then? The error constant sql.ErrNoRows
	// is returned by QueryRow() when the result is empty.
	// This needs to be handled as a special case in most cases.
	// You should only see this error if you're using QueryRow().
	// If you see this error elsewhere, you're doin' it  wrong.
	// ---------------------------------------------------------
	e = row.Scan(CPF_wID...)
	switch e {
	  case sql.ErrNoRows:
	       return nil, 0 // false, nil
	  case nil:
	       return nil, 1 // true, nil 
	  default:
		println("SQL ERROR: (" + e.Error() + ") SQL: " + sSQL)
		return fmt.Errorf("engineunique.get: " +
		       "(%s=%s) failed: %w", pFV.Field, pFV.Value, e), 0
	}
	panic("Oops, fallthru in SELECT")
	
     	// ================================================
	case "M", "U": // Modify, Update
	// https://www.sqlite.org/lang_update.html
	// Obnoxious syntax: 
	// UPDATE tblNm SET fld1=val1,fld2=val2 WHERE expr: 
	// (or..) SET fld1=$1, fld2=$2 WHERE expr; + any...
	// https://www.sqlite.org/syntax/expr.html
	// Use UpdateNames. 
     	// ================================================
     	if pFV == nil {
	   return errors.New("engineunique.update: missing WHERE"), 0 
     	   }
	// -----------------------------------------------------
	// For UPDATE (only), we have to generate here+now an 
	// SQL string that involves all columns (except the ID). 
	// Write assignment pairs as CSV: f1 = $1, f2 = $2, ...
	// We do NOT include the primary key, D.SFT_PRKEY, which
	// is used in the WHERE. 
	// -----------------------------------------------------
	sSQL =	"UPDATE " + pTD.TableSummary.StorName +
		" SET " + pCSVs.UpdateNames +
		" WHERE " + pFV.Field + " = " + pFV.Value + ";"
	fmt.Fprintf(w, "UPDATE.sql: " + sSQL + "\n")
	
	// It is now ready for Exec()
	var theRes sql.Result
	// var newID  int64
	// Call Exec(..) on the stmt, with all column ptrs
	theRes, e = pSR.Handle().Exec(sSQL, CPF...)
	if e != nil {
		fmt.Fprintf(w, "UPDATE.exec: failed: %s", e)
		return fmt.Errorf("engineunique.update.exec: %w", e), 0
	}
	var nRA int64
	nRA, e = theRes.RowsAffected()
	if e != nil {
		fmt.Fprintf(w, "engineunique.update.rowsaffected: failed: %s", e)
		return fmt.Errorf("engineunique.update.rowsaffected: %w", e), 0
	}
	/*
	newID, e = theRes.LastInsertId()
	if e != nil {
		fmt.Fprintf(w, "INSERT.lastinsertId: failed: %s", e)
		return fmt.Errorf("engineunique: insert: lastinsertId: %w",e),0
	}
	fmt.Fprintf(w, "INSERT: OK: LastInsertID: %d \n", newID)
	*/
	return nil, int(nRA) // 1 

// ====================================================

     	// =======================================
	case "D": // Delete, Discard, Drop
	// https://www.sqlite.org/lang_delete.html
	// DELETE FROM tblNm WHERE expr RET'G expr
     	// =======================================
	if pFV == nil {
     	   return errors.New("EngineUnique: SELECT: missing WHERE"), 0 
     	}
     // default:
     }
     return errors.New("engineunique: bad dbOp: " + dbOp), -1
}

func writeFieldDebugInfo(w io.Writer, pTD *DRM.TableDetails) {
     	// TODO: Check the correctness of this! It seemed to overrun with "true"
	CPF := pTD.ColumnPtrsFunc(pTD.NewInstance(), false) // true) // with ID column
	for iCol, cp := range CPF {
	    sn := pTD.ColumnSpecs[iCol].StorName
	    dt := D.SemanticFieldType(pTD.ColumnSpecs[iCol].Datatype)
	    fmt.Fprintf(w, "Column [%d] %s / %s / %T \n", iCol, sn, dt, cp)
	}
}
