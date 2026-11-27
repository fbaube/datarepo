package sqlite

// maybe use ~/go/src/github.com/simukti/sqldb-logger

import(
	"fmt"
	"errors"
	"io"
	"database/sql"
	S "strings"
	"github.com/fbaube/oid"
	D "github.com/fbaube/dsmnd"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// EngineUnique acts on one DB record. All three arguments
// are required, and altho the pointers must be non-nil,
// requirements on their contents depend on which dbOp.
//
// Input arguments:
//  - dbOp is a single letter, and some dbOps have more
//    than one available, as a convenience
//  - anOID is assumed to be every table's primary key
//    and is basically a UUID v7
//  - anOID must always include the target DB table name 
//  - anOID is even required for INSERT, because we
//    assign an OID to a record *before* inserting it 
//  - anOID is even required for KOUNT, because it 
//    includes the target DB table name
//  - pRowBuf is even required for DELETE because its
//    OID must match the OID argument (as an error check)
//  - pRowBuf is even required for KOUNT because the 
//    integer count is written into its last two bytes 
//  - Except for KOUNT, pRowBuf.OID must match anOID 
//
// Outputs: 
//  - error is nil if SQL execution returns no error
//  - If a SELECT fails to find the item, pRowBuf.OID 
//    is set to oid.Nil 
//  - For KOUNT, the integer value is returned in the last
//    two bytes at pRowBuf (this is admittedly hacky) 
//
// DB operations and Use of pRowBuf:
//  - One of four basic actions is performed (listed here as
//    SQL/CRUD/HTTP), and they vary in how pRowBuf is used
//  - INSERT / Create / POST : (record-in) 
//  - UPDATE / Update / PUT  : (record-in) 
//  - SELECT / Retriv / GET  : (record-out) 
//    (if the OID is not found, RowModel.OID <= oid.Nil) 
//  - DELETE / Delete / DELETE (record-in is for OID error checking) 
//  - KOUNT (record-out receives the integer count)
//
// === FURTHER DOCU HERE IS OLD/DODGY ===
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
// If the call to this func is inside a Transaction,
// we probably need to make it available to this func.
// 
// NOTE: When using whereSpec, if a record is not found, this is indicated
// by the first return value (the int), NOT by the second return value 
// (the error, which is reserved for when the DB rejects the SQL.
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
// OLD: func (pSR *SqliteRepo) EngineUnique(dbOp string, tableName
//            string, anID int, pRowBuf DRM.RowModel) (int, error)
// . 
func (pSR *SqliteRepo) EngineUnique(dbOp string, anOID *oid.OID, pRowBuf DRM.RowModel) error {

     var dbOpIsCount = ("K" == S.ToUpper(dbOp[0:1]))
     
     // For all ops, including SELECT and KOUNT,
     // we need pointers anOID and pRowBuf to be non-nil. 
     if anOID == nil { 
	    return errors.New("EngineUnique: nil OID") 
	    }
     if anOID.IsNil() { 
	    return errors.New("EngineUnique: Nil OID") 
	    }
     if pRowBuf == nil { 
	    return errors.New("EngineUnique: nil row buffer") 
	    }
     if *anOID != *pRowBuf.OID() && !dbOpIsCount {
        return errors.New("EngineUnique: OID != row buffer OID")
        }
     var tableName string
     tableName = anOID.TableName()
     if !anOID.HasValidTableName() {
         return errors.New("EngineUnique: bad DB table name <" + tableName + ">")
         }
         
     var pTD    *DRM.TableDetails
     var pCSVs  *DRM.ColumnStringsCSV
     var pStmts *DRM.Statements
     var w = pSR.LogWriter()
     var e error
     var dbOp1 string     // single character 
     var SQL_toUse string // , WHERE_toUse string FIXME FIXME FIXME
     // Table's column-pointers funcs 
     var CPF []any

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
     // because the ptrs depend on the particular value of pRowBuf.
     CPF = pTD.ColumnPtrsFunc(pRowBuf) 

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
//  - INSERT / Create / POST : (record-in) 
//  - UPDATE / Update / PUT  : (record-in) 
//  - SELECT / Retriv / GET  : (record-out) 
//    (if the OID is not found, RowModel.OID <= oid.Nil) 
//  - DELETE / Delete / DELETE (record-in is for OID error checking) 
//  - KOUNT (record-out receives the integer count)

	// "+" Add, Create, Insert, New
	case "A", "C", "I", "N":
	     SQL_toUse = pStmts.INSERTuniqueID
	     // USE QueryRow for INSERT! So that LastResultID works ? No.
	     // https://github.com/mattn/go-sqlite3/issues/1140
	     // useQueryRow = true // <==
	     dbOp1 = "+"
    // "?" Fetch, Get, List, Retrieve, Select 
    case "F", "G", "L", "R", "S":
	     SQL_toUse = pStmts.SELECTuniqueID
	     dbOp1 = "?"
	// "=" Modify, Update
	case "M", "U":
	     SQL_toUse = pStmts.UPDATEuniqueID
	     dbOp1 = "="
	// "-" Delete, Discard, Drop
	case "D":
	     SQL_toUse = pStmts.DELETEuniqueID
	     dbOp1 = "-"
	// "#" Kount
	case "K":
	     dbOp1 = "#"
         SQL_toUse = pStmts.COUNTuniqueID
         // FIXME TBS 
	default:
         return errors.New("EngineUnique: bad dbOp <" + dbOp + ">") 
    }

     dbOpString := dbOp + "(" + dbOp1 + ")"
     dbOpError  := "EngineUnique: " + dbOpString + ": "
     fmt.Fprintln(w, dbOpString, ",", SQL_toUse)

//   var cpf []any
//	 cpf = pTD.ColumnPtrsFunc(pRowBuf)
     
     // ================
     //  TIME to EXECUTE 	
     // ================
     if dbOp1 == "?" { // SELECT // might aalso be: if !dbOpIsCount { 
         fmt.Fprintf(w, "QueryRow: %s / %d / %s \n",
             dbOp1, anOID, SQL_toUse)
	     // FIXME: Transaction, if active ?
	     row := pSR.Handle().QueryRow(SQL_toUse, anOID)
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
           pRowBuf.SetOID(oid.Nil)
	       return nil // no row, no error
	  case nil:
	       return nil // 1 row, no error 
	  default:
		println("SQL ERROR: (" + e.Error() + ") SQL: " + SQL_toUse)
		return fmt.Errorf(dbOpError + 
		       "ID <%d> failed: %w", anOID, e)
	}
	panic("Oops, fallthru in SELECT")
	
      } else {
	var newID int64
     // if dbOp1 == "+" { CPF_toUse = append(CPF_toUse, &newID) } 
     	fmt.Fprintf(w, "Exec: %d / w %d parms / %s \n",
		anOID, len(CPF), SQL_toUse)
	// It is now ready for Exec()
	var theRes sql.Result
	// FIXME: Transaction, if active ?
	theRes, e = pSR.Handle().Exec(SQL_toUse, CPF...)
	if e != nil {
		fmt.Fprintf(w, dbOpError + "exec failed: %s", e)
		return fmt.Errorf(dbOpError + "exec: %w", e) 
	}
	
	if dbOp1 == "+" { // INSERT 
	   // Did NOT use RETURNING to get new ID. 
	   newID, e = theRes.LastInsertId()
	   if e != nil {
		fmt.Fprintf(w, "engineunique.insert.lastinsertId: failed: %s", e)
		return fmt.Errorf("engineunique.insert: lastinsertId: %w", e) 
		} 
	   fmt.Fprintf(w, "INSERT (%s) OK: LastInsertID: %d \n",
	   	tableName, int(newID))
	   return nil 
	}
	// UPDATE, DELETE
    /*
	var nRA int64
	nRA, e = theRes.RowsAffected()
	if e != nil {
		fmt.Fprintf(w, "engineunique.update.rowsaffected: failed: %s", e)
		return fmt.Errorf("engineunique.update.rowsaffected: %w", e) 
		}
        */
	return nil // nr of rows affected: 0 or 1 
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
