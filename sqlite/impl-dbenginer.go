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
//  - SELECT / Retrieve / GET
//  - INSERT / Create / POST 
//  - UPDATE / Update / PUT
//  - DELETE / Delete / DELETE
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
//  - An int that is (if INSERT) the newly-added row ID (else) 0 or 1
//    to indicate whether a record was affected.
//
// NOTE: When using whereSpec, if a record is not found, this is indicated
// by the second return value (the int), NOT by the error, which is reserved
// for when the DB rejects the SQL.
//
// NOTE: In an UPDATE, if the whereSpec does not refer to the ID, and the
// ID of the input record does not match the ID of the record found by the
// DB, the function panics. So, for UPDATE, just match on the ID. 
//
// NOTE: Also implement COUNT(*) ?
// . 
func (pSR *SqliteRepo) EngineUnique(dbOp string, tableName string, whereSpec *DRP.UniquerySpec, RM DRM.RowModel) (error, int) {

     var pTD *DRM.TableDetails
     var w = pSR.LogWriter()
     var sSQL string 
     // Column-pointers function 
     var idxdCPF, CPF []any  // with ID; no ID
     // Comma-separated column names 
     var idxdCSV, CSV string // with ID; no ID
     // '$'-numbered parameters (Postgres-style)
     var idxdPlcNrs, PlcNrs string // with ID; no ID 

     pTD = GetTableDetailsByCode(tableName)
     // TODO: Handle nil return 
     
     CSV = pTD.ColumnNamesCSV // no ID column 
     CPF = pTD.ColumnPtrsFunc(RM, false) // no ID column 
     idxdCSV = pTD.PKname + ", " + CSV // "IDX_" + pTD.StorName + ", " + CSV
     idxdCPF = pTD.ColumnPtrsFunc(RM, true) // with ID column
     var i int 
     for i = range len(CPF) {
	 PlcNrs += fmt.Sprintf("$%d, ", i+1) }
     for i = range len(idxdCPF) {
	 idxdPlcNrs += fmt.Sprintf("$%d, ", i+1) }
     PlcNrs     = S.TrimSuffix(PlcNrs, ", ")
     idxdPlcNrs = S.TrimSuffix(idxdPlcNrs, ", ")

     // Log info about the columns 
     writeFieldDebugInfo(w, pTD)

     // switch dbOp {
     switch S.ToUpper(dbOp[0:0]) { 

     	// ======================================================
	case "A", "C", "I", "N":
	// Add, Create, Insert, New 
	// https://www.sqlite.org/lang_insert.html
	// INSERT INTO tblNm (fld1, fld2) VALUES(val1, val2);
	// INSERT INTO tblNm (fld1, fld2) VALUES($1,$2); + any...
     	// ======================================================
	if whereSpec != nil {
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
	if whereSpec == nil {
	   return errors.New("EngineUnique: SELECT: missing WHERE"), 0
	}
	sSQL =  "SELECT " + idxdCSV +
		" FROM " + pTD.TableSummary.StorName +
		" WHERE " + whereSpec.Field + " = " + whereSpec.Value + ";"
	// TODO: QueryRow(..) 
	
     	// ============================================
        case "M", "U": // Modify, Update
	// https://www.sqlite.org/lang_update.html
	// UPDATE tblNm SET stuff WHERE expr RET'G expr 
	// https://www.sqlite.org/syntax/expr.html 
     	// ============================================
     	if whereSpec == nil {
           return errors.New("EngineUnique: SELECT: missing WHERE"), 0 
     	   }
     	// =======================================
        case "D": // Delete, Discard, Drop
	// https://www.sqlite.org/lang_delete.html
	// DELETE FROM tblNm WHERE expr RET'G expr
     	// =======================================
	if whereSpec == nil {
     	   return errors.New("EngineUnique: SELECT: missing WHERE"), 0 
     	}
     // default:
     }
     return errors.New("engineunique: bad dbOp: " + dbOp), -1
}

func writeFieldDebugInfo(w io.Writer, pTD *DRM.TableDetails) {
	CPF := pTD.ColumnPtrsFunc(pTD.BlankInstance, true) // with ID column
	for iCol, cp := range CPF {
	    sn := pTD.ColumnSpecs[iCol].StorName
	    dt := D.SemanticFieldType(pTD.ColumnSpecs[iCol].Datatype)
	    fmt.Fprintf(w, "Column [%d] %s / %s / %T \n", iCol, sn, dt, cp)
	}
}
