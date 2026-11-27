package sqlite

// maybe use ~/go/src/github.com/simukti/sqldb-logger

import(
	"fmt"
	"errors"
	"io"
    "uuid"
	"database/sql"
	S "strings"
	"github.com/fbaube/nork"
	D "github.com/fbaube/dsmnd"
	// DRP "github.com/fbaube/datarepo"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// EngineUnique acts on a single DB record, based on the
// value of the UUID v7 primary key; this includes INSERT,
// because we assign the UUID before inserting.
//
//  - Inputs are calling arguments DB-Op and UUIDv7 and
//    maybe the buffer at calling argument *RowModel;
//    UUIDv7 must include the DB table name 
//  - Ouputs are return value error (nil is success) and 
//    maybe  the buffer at calling argument *RowModel 
//
// The basic signature is (error) = func(op,UUIDv7,*RowModel):
// One of four basic actions is performed (listed as SQL/CRUD/HTTP):
//  - INSERT / Create / POST   (record in) 
//  - UPDATE / Update / PUT    (record in) 
//  - SELECT / Retriv / GET    (record out; *Rowmodel has nork.NilUUID if not found)
//  - DELETE / Delete / DELETE (record in - mandatory for now) 
//  - For UPDATE and DELETE,  calling argument UUID must
//    match the UUID in the *Rowmodel, and the UUID must
//    include the DB table name; otherwise the func panics. 
//
// OBSOLETE: It takes four input arguments:
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
// If the call to this func is inside a Transaction,
// we probably need to make it available to this func.
// 
// NOTE: When using whereSpec, if a record is not found, this is indicated
// by the first return value (the int), NOT by the second return value 
// (the error, which is reserved for when the DB rejects the SQL.
//
// NOTE: Also implement COUNT(*) ? Perhaps as DB Op "K" = Kount. 
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
//
// NOTE: When writing the multi-row version of this,
// be sure to call Rows.Cloe()
// . 
// func (pSR *SqliteRepo) EngineUnique(dbOp string, tableName
//       string, anID int, pRMbuf DRM.RowModel) (int, error)  
func (pSR *SqliteRepo) EngineUnique(dbOp string, u *uuid.UUID, pRMbuf DRM.RowModel) error {

     if u == nil || *u == nork.NilUUID || *u != *pRMbuf.UUID() {
        panic("EngineUnique got bad inputs")
     }
     var pTD    *DRM.TableDetails
     var pCSVs  *DRM.ColumnStringsCSV
     var pStmts *DRM.Statements
     var w = pSR.LogWriter()
     var e error
     var dbOp1 string // single character 
     var useQueryRow bool // true only for SELECT 
     var SQL_toUse string // , WHERE_toUse string FIXME FIXME FIXME
     // Table's column-pointers funcs 
     var CPF []any

     var tableName string
     tableName = nork.TableNameInUUID(*u) // FIXME: Is it OK Ascii ? 

     // Fetch the table's details 
     pTD = GetTableDetailsByStorName(tableName)
     if pTD == nil {
     	pTD = GetTableDetailsByDispName(tableName)
     }
     if pTD == nil {
     	// FIXME err msgs 
     	s := "NO TblDtls FOR: " + tableName
     	println(s)
	    errors.New(s) 
    }
     // Fill in convenience vars 
     pCSVs  = pTD.CSVs
     pStmts = pTD.Stmts
     if pCSVs  == nil { panic("nil TableDetails ColumnStrings") }
     if pStmts == nil { panic("nil TableDetails Statements") }
     
     // This CPF cannot be pre-generated during initialization
     // because the ptrs depend on the particular value of pRMbuf.
     CPF = pTD.ColumnPtrsFunc(pRMbuf) 

     // Log info about the columns 
     // writeFieldDebugInfo(w, pTD)

     // Convert dbOp from a long string to a single
     // character (one of "+", "-", "=", "?", "#"), 
     // and choose the correct column ptrs func. 
     // We check only the first letter of the DB op, 
     // so user can (and should?) be creative with
     // the string passed in :-P
     fmt.Fprintln(w, "DB OP IS: " + S.ToUpper(dbOp))
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
	     // USE QueryRow for INSERT! So that LastResultID works ? No.
	     // https://github.com/mattn/go-sqlite3/issues/1140
	     // useQueryRow = true // <==
	     dbOp1 = "+"
        // ? Fetch, Get, List, Retrieve, Select 
        case "F", "G", "L", "R", "S":
	     SQL_toUse = pStmts.SELECTuniqueID
	     useQueryRow = true // <==
	     dbOp1 = "?"
	// = Modify, Update
	case "M", "U":
	     SQL_toUse = pStmts.UPDATEuniqueID
	     dbOp1 = "="
	// - Delete, Discard, Drop
	case "D":
	     SQL_toUse = pStmts.DELETEuniqueID
	     dbOp1 = "-"
	// # Kount
	case "K":
	     dbOp1 = "#"
	default:
	return errors.New("engineunique: bad dbOp: " + dbOp) 
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
	cpf = pTD.ColumnPtrsFunc(pRMbuf)
	cpfid := (cpf [0]).(*uuid.UUID)
	bufID = *cpfid
     }
     // Check re. WHERE ID spec and/or input/output buffer.
     // if INSERT or KOUNT, no WHERE ID.
     if dbOp1 == "+" || dbOp1 == "#" {
     	// Allow the zero value :-D 
	if anID > 0 || bufID > 0 { // rather than: >= 0 
	   return errors.New("EngineUnique: INSERT: unwanted WHERE ID") 
	   }
	ID_toUse = -1 
     // else is ?/=/- SELECT/UPDATE/DELETE 
      } else {
	// Need a search spec: (1) the "anID" int argument, and/or 
	// (2) the "pRMbuf" buffer argument (where check for an ID).
	   if (anID == -1) && (pRMbuf == nil) {
	       return errors.New(dbOpError + 
		  "missing search spec (no WHERE ID spec or buffer with ID") 
	   }
	   if (anID >= 0) && (pRMbuf != nil ) {
	       if anID != bufID {
	       	  if bufID == 0 {
		     ermsg := fmt.Sprintf("WARN: " +
		     	     "Input ID %d != buffer ID 0", anID)
		     println(ermsg)
		     fmt.Fprintln(w, ermsg)
		  } else {
	 	    return fmt.Errorf(dbOpError + "conflicting ID spec: " +
		  	 "arg <%d> != buffer's ID <%d>", anID, bufID)
		  }
	       }
	   }
	   ID_toUse = anID 
 	// If UPDATE, need an  input buffer 
 	// If SELECT, need an output buffer 
	   if (dbOp1 != "-") && (pRMbuf == nil) {
	    	return errors.New(dbOpError + "missing buffer") 
	   }
	   
     }
     
     // ================
     //  TIME to EXECUTE 	
     // ================
     if useQueryRow { // "?" 
     	fmt.Fprintf(w, "QueryRow: %s / %d / %s \n", dbOp1, ID_toUse, SQL_toUse)
	// FIXME: Transaction, if active ?
	row := pSR.Handle().QueryRow(SQL_toUse, ID_toUse)
	// ---------------------------------------------------------
	// What if there is no row in the result, and .Scan() can't
	// scan a value. What then? The error constant sql.ErrNoRows
	// is returned by QueryRow() when the result is empty.
	// This needs to be handled as a special case in most cases.
	// You should only see this error if you're using QueryRow().
	// If you see this error elsewhere, yer doin' it wrong.
	// ---------------------------------------------------------
	e = row.Scan(CPF...) 
	fmt.Fprintf(w, "QueryRow: ret: %v \n", e)
	
	switch e {
	  case sql.ErrNoRows:
           pRMbuf.UUID = nork.NilUUID
	       return nil // no row, no error
	  case nil:
	       return nil // 1 row, no error 
	  default:
		println("SQL ERROR: (" + e.Error() + ") SQL: " + SQL_toUse)
		return 0, fmt.Errorf(dbOpError + 
		       "ID <%d> failed: %w", ID_toUse, e)
	}
	panic("Oops, fallthru in SELECT")
	
      } else {
	var newID int64
     // if dbOp1 == "+" { CPF_toUse = append(CPF_toUse, &newID) } 
     	fmt.Fprintf(w, "Exec: %d / w %d parms / %s \n",
		ID_toUse, len(CPF), SQL_toUse)
	// It is now ready for Exec()
	var theRes sql.Result
	// FIXME: Transaction, if active ?
	theRes, e = pSR.Handle().Exec(SQL_toUse, CPF...)
	if e != nil {
		fmt.Fprintf(w, dbOpError + "exec failed: %s", e)
		return 0, fmt.Errorf(dbOpError + "exec: %w", e) 
	}
	
	if dbOp1 == "+" { // INSERT 
	   // Did NOT use RETURNING to get new ID. 
	   newID, e = theRes.LastInsertId()
	   if e != nil {
		fmt.Fprintf(w, "engineunique.insert.lastinsertId: failed: %s", e)
		return 0, fmt.Errorf("engineunique.insert: lastinsertId: %w", e) 
		} 
	   fmt.Fprintf(w, "INSERT (%s) OK: LastInsertID: %d \n",
	   	tableName, int(newID))
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
	CPF := pTD.ColumnPtrsFunc(pTD.NewInstance()) 
	for iCol, cp := range CPF {
	    sn := pTD.ColumnSpecs[iCol].StorName
	    dt := D.SemanticFieldType(pTD.ColumnSpecs[iCol].Datatype)
	    fmt.Fprintf(w, "Column [%d] %s / %s / %T \n", iCol, sn, dt, cp)
	}
}
