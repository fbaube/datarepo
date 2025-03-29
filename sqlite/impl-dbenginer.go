package sqlite

// maybe use ~/go/src/github.com/simukti/sqldb-logger

import(
	"fmt"
	"errors"
	"io"
	"database/sql"
	S "strings"
	D "github.com/fbaube/dsmnd"
	// DRP "github.com/fbaube/datarepo"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// EngineUnique acts on a single DB record, based on the
// value of the primary key (except of course for INSERT).
//
// (It was thought that instead of using only the ID, there
// could be a WHERE clause that uses any column specified to
// be UNIQUE, but actually that approach makes no sense.)
//
// The basic signature is (int,error) = func(op,table,int,buffer):
// One of four basic actions is performed (listed as SQL/CRUD/HTTP):
//  - INSERT / Create / POST (record in) (returns new-ID)
//  - UPDATE / Update / PUT  (record in) (returns 0/1)
//  - SELECT / Retriv / GET      (ID in) (returns 0/1 + record) 
//  - DELETE / Delete / DELETE   (ID in) (returns 0/1)
//  - let nAR = nr of records affected 
//  - (newID,e)  = insert(0,inbuffer)  // optimize for use in batches 
//  - (nAR,e) = update(anID,inbuffer)  // anID can be -1, else match buffer's
//  - (nAR,e) = select(anID,outbuffer) // anID >= 0
//  - (nAR,e) = delete(anID,nil)       // anID >= 0; buffer OK, then ID's match
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
// by the first return value (the int), NOT by the second return value 
// (the error, which is reserved for when the DB rejects the SQL.
//
// NOTE: If anID >= 0 and buffer is non-nil, and anID does not match the ID
// of the record in the buffer, the function panics. 
//
// NOTE: Also implement COUNT(*) ? Perhaps as "K" = Kount. 
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
func (pSR *SqliteRepo) EngineUnique(dbOp string, tableName string, anID int, pRMbuf DRM.RowModel) (int, error) {

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
     pTD = GetTableDetailsByString(tableName)
     if pTD == nil {
     	// FIXME err msgs 
     	s := "NO TblDtls FOR: " + tableName
     	println(s)
	return 0, errors.New(s) 
     }
     // Fill in convenience vars 
     pCSVs  = pTD.CSVs
     pStmts = pTD.Stmts
     if pCSVs  == nil { panic("nil TableDetails ColumnStrings") }
     if pStmts == nil { panic("nil TableDetails Statements") }
     
     // These CPF vars cannot be pre-generated during initialization
     // because the ptrs depend on the particular value of pRMbuf.
     CPF_noID = pTD.ColumnPtrsFunc(pRMbuf, false)   // no ID column 
     CPF_wID  = pTD.ColumnPtrsFunc(pRMbuf, true)   // with ID column

     // Log info about the columns 
     // writeFieldDebugInfo(w, pTD)

     // Convert dbOp from a long string to a single
     // character (one of "+", "-", "=", "?", "#"), 
     // and choose the correct column ptrs func. 
     // We check only the first letter of the DB op, 
     // so user can (and should?) be creative with
     // the string passed in :-P
     fmt.Fprintln(w, "DB OP IS: " + S.ToUpper(dbOp[0:1]))
     switch S.ToUpper(dbOp[0:1]) {

// Recapping: 
// The basic signature is (int,error) = func(op,table,int,buffer):
// One of four basic actions is performed (listed as SQL/CRUD/HTTP):
//  - INSERT / Create / POST (record in) (returns new-ID)
//  - UPDATE / Update / PUT  (record in) (returns 0/1)
//  - SELECT / Retriv / GET      (ID in) (returns 0/1 + record) 
//  - DELETE / Delete / DELETE   (ID in) (returns 0/1)
//  - let nAR = nr of records affected 
//  - (newID,e)  = insert(0,inbuffer)  // optimize for use in batches 
//  - (nAR,e) = update(anID,inbuffer)  // anID can be -1, else match buffer's
//  - (nAR,e) = select(anID,outbuffer) // anID >= 0
//  - (nAR,e) = delete(anID,nil)       // anID >= 0; if buffer, ID's must match

	// + Add, Create, Insert, New
	case "A", "C", "I", "N":
	     SQL_toUse = pStmts.INSERTuniqueID
	     // We'll get the new ID OK, so don't try
	     // to rely on hypothetical writeback 
	     CPF_toUse = CPF_noID 
	     dbOp1 = "+"
        // ? Fetch, Get, List, Retrieve, Select 
        case "F", "G", "L", "R", "S":
	     SQL_toUse = pStmts.SELECTuniqueID
	     // WHERE uses ID, so don't NEED to pass
	     // via CPF, but let's use the writeback 
	     CPF_toUse = CPF_wID
	     useQueryRow = true // <==
	     dbOp1 = "?"
	// = Modify, Update
	case "M", "U":
	     SQL_toUse = pStmts.UPDATEuniqueID
	     // WHERE uses ID, and no writeback, so don't pass via CPF 
	     CPF_toUse = CPF_noID
	     dbOp1 = "="
	// - Delete, Discard, Drop
	case "D":
	     SQL_toUse = pStmts.DELETEuniqueID
	     dbOp1 = "-"
	// # Kount
	case "K":
	     dbOp1 = "#"
	default:
	return 0, errors.New("engineunique: bad dbOp: " + dbOp) 
     }

     dbOpString := dbOp + "(" + dbOp1 + ")"
     dbOpError := "engineunique: " + dbOpString + ": "
     fmt.Fprintln(w, dbOpString, ",", SQL_toUse)

     // anID is the argument ID passed in.
     // Now extract the ID if a buffer was passed in. 
     var bufID = -1
     var ID_toUse int 
     if pRMbuf != nil {
     	var cpf []any
	cpf = pTD.ColumnPtrsFunc(pRMbuf, true)
	cpfid := (cpf [0]).(*int)
	bufID = *cpfid
     }
     // Check re. WHERE ID spec and/or input/output buffer.
     // if INSERT or KOUNT, no WHERE ID.
     if dbOp1 == "+" || dbOp1 == "#" {
     	// Allow the zero value :-D 
	if anID > 0 || bufID > 0 { // rather than: >= 0 
	   return 0, errors.New("EngineUnique: INSERT: unwanted WHERE ID") 
	   }
     // else is ?/=/- SELECT/UPDATE/DELETE 
      } else {
	// Need a search spec: (1) the "anID" int argument, and/or 
	// (2) the "pRMbuf" buffer argument (where check for an ID).
	   if (anID == -1) && (pRMbuf == nil) {
	       return 0, errors.New(dbOpError + 
		  "missing search spec (no WHERE ID spec or buffer with ID") 
	   }
	   if (anID >= 0) && (pRMbuf != nil ) {
	       if anID != bufID {
	       	  return 0, fmt.Errorf(dbOpError + "conflicting ID spec: " +
		  	 "arg <%d> != buffer's ID <%d>", anID, bufID) 
	       }
	   }
	   ID_toUse = anID 
 	// If UPDATE, need an  input buffer 
 	// If SELECT, need an output buffer 
	   if (dbOp1 != "-") && (pRMbuf == nil) {
	    	return 0, errors.New(dbOpError + "missing buffer") 
	   }
	   
     }
     // Handle WHERE clause
   

     // ================
     //  TIME to EXECUTE 	
     // ================
     if useQueryRow { // "?"
     	fmt.Fprintf(w, "QueryRow: %d / %s \n", ID_toUse, SQL_toUse)
	row := pSR.Handle().QueryRow(SQL_toUse, ID_toUse)
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
	       return 0, nil // no error, no nRows 
	  case nil:
	       return 1, nil // no error, 1 row 
	  default:
		println("SQL ERROR: (" + e.Error() + ") SQL: " + SQL_toUse)
		return 0, fmt.Errorf(dbOpError + 
		       "ID <%d> failed: %w", ID_toUse, e)
	}
	panic("Oops, fallthru in SELECT")
	
     } else { 
	// It is now ready for Exec()
	var theRes sql.Result
	var newID  int64
	theRes, e = pSR.Handle().Exec(SQL_toUse, CPF_toUse...)
	if e != nil {
		fmt.Fprintf(w, dbOpError + "exec failed: %s", e)
		return 0, fmt.Errorf(dbOpError + "exec: %w", e) 
	}
	
	if dbOp == "+" { // INSERT 
	// Used RETURNING to get new ID. 
	// Call Exec(..) on the stmt, with all column ptrs (except ID) 
	   newID, e = theRes.LastInsertId()
	   if e != nil {
		fmt.Fprintf(w, "engineunique.insert.lastinsertId: failed: %s", e)
		return 0, fmt.Errorf("engineunique.insert: lastinsertId: %w",e) 
		}
	   fmt.Fprintf(w, "INSERT: OK: LastInsertID: %d \n", newID)
	   return int(newID), nil 
	}
	// UPDATE, DELETE 
	var nRA int64
	nRA, e = theRes.RowsAffected()
	if e != nil {
		fmt.Fprintf(w, "engineunique.update.rowsaffected: failed: %s", e)
		return 0, fmt.Errorf("engineunique.update.rowsaffected: %w", e) 
		}
	return int(nRA), nil // nr of rows affected: 0 or 1 
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
