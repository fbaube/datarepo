package sqlite

// maybe use ~/go/src/github.com/simukti/sqldb-logger

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
//    (NOTE: We might just use an ID passed in in the buffer.) 
//  - A pointer to a buffer, used for input (if INSERT or UPDATE) or output
//    (if SELECT); for INSERT, the  buffer is left unmodified and the ID of
//    the new record is returned in the second return value (the int) 
//  - Note that the buffer pointer is an interface, implemented by a pointer 
//    to a struct; this may be usable by GO generics
//
// It returns two values:
//  - An error: if it is non-nil, the other (i.e. int) return value is invalid
//  - An int that is (if INSERT) the newly-added row ID (else) 0 or 1 to
//    indicate how many records were (i.e. whether a record was) affected.
//
// NOTE: When using whereSpec, if a record is not found, this is indicated
// by the second return value (the int), NOT by the first return value 
// (the error, which is reserved for when the DB rejects the SQL.
//
// NOTE: In an UPDATE, if the whereSpec does not refer to the ID, and the
// ID of the input record does not match the ID of the record found by the
// DB, the function panics. So, for UPDATE, it might be wise only to pass 
// in the ID for the WHERE. 
//
// NOTE: Also implement COUNT(*) ?
//
// TODO: TD.pCSVs could also use DB.Prepare to gather sql.Stmt's, but then we
// have to pass in a Connection, which screws up modularity pretty severely.
//
// TODO: switch on dbOp to call a new mini func that assembles the SQL statement.
// 
// TODO: Use Result.RowsAffected
// https://pkg.go.dev/database/sql#Result
// RowsAffected returns the number of rows affected by an update, insert,
// or delete. Not every DB or driver supports it. [altho mattn/SQLite does.]
// RowsAffected() (int64, error)
//
// NOTE: When writing the multi-row version of this,
// be sure to call Rows.Cloe()
// . 
func (pSR *SqliteRepo) EngineUnique(dbOp string, tableName string, pFV *DRP.FieldValuePair, pRMbuf DRM.RowModel) (error, int) {

     var pTD    *DRM.TableDetails
     var pCSVs  *DRM.ColumnStringsCSV
     var pStmts *DRM.Statements
     var w = pSR.LogWriter()
     var e error
     var dbOp1 string // single character 
     var useQueryRow bool // true only for SELECT 
     var SQL_toUse string // , WHERE_toUse string FIXME FIXME FIXME
     // Table's column-pointers funcs 
     var CPF_noID, CPF_wID, CPF_toUse []any 

     // Fetch the table's details 
     pTD = GetTableDetailsByCode(tableName)
     if pTD == nil {
     	// FIXME err msgs 
     	s := "NO TblDtls FOR: " + tableName
     	println(s)
	return errors.New(s), 0
     }
     // Fill in vars 
     pCSVs  = pTD.CSVs
     pStmts = pTD.Stmts
     if pCSVs  == nil { panic("nil TableDetails ColumnStrings") }
     if pStmts == nil { panic("nil TableDetails Statements") }
     
     // For convenience, WHERE can use "ID"
     // (without table name), and we fix it 
     if pFV != nil && S.EqualFold("id", pFV.Field) {
     	println("FIXED WHERE:", pFV.Field, "=>", pTD.PKname)
     	pFV.Field = pTD.PKname
     }
     CPF_noID = pTD.ColumnPtrsFunc(pRMbuf, false)   // no ID column 
     CPF_wID  = pTD.ColumnPtrsFunc(pRMbuf, true)   // with ID column

     // Log info about the columns 
     // writeFieldDebugInfo(w, pTD)

     // Convert dbOp from a long string to a single
     // character (one of "+", "-", "=", "?"), and 
     // choose the correct column ptrs func. 
     // We check only the first letter of the DB op, 
     // so be creative with the string passed in :-P
     // println("DB OP IS: " + S.ToUpper(dbOp[0:1]))
     switch S.ToUpper(dbOp[0:1]) {
     
	// + Add, Create, Insert, New
	case "A", "C", "I", "N":
	     SQL_toUse = pStmts.INSERTunique
	     CPF_toUse = CPF_noID 
	     dbOp1 = "+"
        // ? Fetch, Get, List, Retrieve, Select 
        case "F", "G", "L", "R", "S":
	     SQL_toUse = pStmts.SELECTunique
	     CPF_toUse = CPF_wID
	     useQueryRow = true // <==
	     dbOp1 = "?"
	// = Modify, Update
	case "M", "U":

	     dbOp1 = "="
	// - Delete, Discard, Drop
	case "D":

	     dbOp1 = "-"
	default:
	return errors.New("engineunique: bad dbOp: " + dbOp), 0
     }

     dbOpString := dbOp + "(" + dbOp1 + ")"
     dbOpError := "engineunique: " + dbOpString + ": "

     // Check re. WHERE spec and/or input/output buffer.
     // if INSERT, no WHERE.
     if dbOp1 == "+" {
	if pFV != nil {
	   return errors.New("EngineUnique: INSERT: unwanted WHERE"), 0 
	   }
     // else is SELECT/UPDATE/DELETE 
      } else {
	// Need a search spec: either (1) a WHERE spec, or 
	// (2) an input buffer (where we can find an ID).
	// TODO: If a WHERE spec, it must be on a UNIQUE column. 
	   if (pFV == nil) && (pRMbuf == nil ) {
	       return errors.New(dbOpError + 
		  "missing search spec (no WHERE spec or buffer with ID"), 0 
	   }
 	// If UPDATE, need an  input buffer 
 	// If SELECT, need an output buffer 
	   if (dbOp1 != "-") && (pRMbuf == nil) {
	    	return errors.New(dbOpError + "missing buffer"), 0 
	   }
	// Now set WHERE_toUse !!! TODO TODO TODO 
     }

     // ================
     //  TIME to EXECUTE 	
     // ================
      fmt.Fprintf(w, "SQL: " + SQL_toUse + "\n")	

     if useQueryRow { // "?" 
	row := pSR.Handle().QueryRow(SQL_toUse)
	// ---------------------------------------------------------
	// What if there is no row in the result, and .Scan() can't
	// scan a value. What then? The error constant sql.ErrNoRows
	// is returned by QueryRow() when the result is empty.
	// This needs to be handled as a special case in most cases.
	// You should only see this error if you're using QueryRow().
	// If you see this error elsewhere, yer doin' it wrong.
	// ---------------------------------------------------------
	e = row.Scan(CPF_toUse...) // _noID // BUT WHAT ABOUT no-WHERE ???
	switch e {
	  case sql.ErrNoRows:
	       return nil, 0 // false, nil
	  case nil:
	       return nil, 1 // true, nil 
	  default:
		println("SQL ERROR: (" + e.Error() + ") SQL: " + SQL_toUse)
		return fmt.Errorf(dbOpError + 
		       "(%s=%s) failed: %w", pFV.Field, pFV.Value, e), 0
	}
	panic("Oops, fallthru in SELECT")
	
     } else { 
	// It is now ready for Exec()
	var theRes sql.Result
	var newID  int64
	theRes, e = pSR.Handle().Exec(SQL_toUse, CPF_toUse...)
	if e != nil {
		fmt.Fprintf(w, dbOpError + "exec failed: %s", e)
		return fmt.Errorf(dbOpError + "exec: %w", e), -1
	}
	
	if dbOp == "+" { // INSERT 
	// Used RETURNING to get new ID. 
	// Call Exec(..) on the stmt, with all column ptrs (except ID) 
	   newID, e = theRes.LastInsertId()
	   if e != nil {
		fmt.Fprintf(w, "engineunique.insert.lastinsertId: failed: %s", e)
		return fmt.Errorf("engineunique.insert: lastinsertId: %w",e),-1
		}
	   fmt.Fprintf(w, "INSERT: OK: LastInsertID: %d \n", newID)
	   return nil, int(newID)
	}
	// UPDATE, DELETE 
	var nRA int64
	nRA, e = theRes.RowsAffected()
	if e != nil {
		fmt.Fprintf(w, "engineunique.update.rowsaffected: failed: %s", e)
		return fmt.Errorf("engineunique.update.rowsaffected: %w", e), 0
		}
	return nil, int(nRA) // nr of rows affected: 0 or 1 
     }
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
