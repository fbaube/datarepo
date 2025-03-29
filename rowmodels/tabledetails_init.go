package rowmodels

import(
	"fmt"
	"errors"
	"strconv"
	S "strings"
	D "github.com/fbaube/dsmnd"
)

// terms: Column Field Key Prop / Value

// csvNumbers is 1-based 
func csvNumbers(pfx string, min int, max int) string {
     var sb S.Builder
     var i int 
     for i = min; i < max; i++ {
     	 _,_ = sb.WriteString(pfx + strconv.Itoa(i) + ", ")
     }
     _,_ = sb.WriteString(pfx + strconv.Itoa(max))
     return sb.String()
}

// csvNames is 0-based. 
func csvNames(CSs []D.ColumnSpec, min int, max int) string {
     var sb S.Builder
     var i int 
     for i = min; i < max; i++ {
     	 _,_ = sb.WriteString(CSs[i].StorName + ", ")
     }
     _,_ = sb.WriteString(CSs[max].StorName)
     return sb.String()
}

// GenerateColumnStringsCSV generates struct [ColumnStringsCSV] 
// for every struct that has been registered using method
// datarepo/RegisterAppTables of interface datarepo/SimpleRepo .
//
// TODO: It should be guarded by a Do.Once()
//
// It does not modify or even access the database. It should 
// be called ASAP after program start. It does not need to be 
// called before a database is opened, but it DOES need to be 
// called any SQL is executed against the database.
//
// Note that a table's primary key, named "{table}_ID", 
// is not listed in the slice of TableDetails.
// .
func GenerateColumnStringsCSV(pTD *TableDetails) error {
     if pTD.CSVs != nil {
     	println("GenerateColumnStringsCSV: DUPE CALL")
	return errors.New("GenerateColumnStringsCSV: dupe initialisation")
	}    
     pTD.CSVs = new(ColumnStringsCSV)
     var CSs []D.ColumnSpec
     CSs   = pTD.ColumnSpecs
     var N = 1+len(pTD.ColumnSpecs) // add in the ID 
     // var cs D.ColumnSpec 

     // PRIMARY KEY (aka "ID") is NOT in ColumnSpecs !!

     // So if (for example) we have a total (incl. ID) of N fields:
     // - The first is the ID, and there are N-1 others
     // - The N  fields are indexed in ColumnSpecs as 0..N-1, i.e. [:N]
     // - Including the ID, the placeholders are numbered $1..$N
     // - Not including it, the placeholders are numbered $1..$N-1

     // For example, if we assume an ID plus four (4)
     // other fields, for a total of << N=5 >> fields
     // (ID plus four others, f[0..3]):
     // type struct { id, f0, f1, f2, f3 } 
     // type ColumnStringsCSV struct {
     // 4 FieldNames_noID        "f0, f1, f2, f3" 
     // 4 PlaceNums_noID         "$1, $2, $3, $4" 
     // 5 FieldNames_wID     "id, f0, f1, f2, f3" 
     // 5 PlaceNums_wID      "$1, $2, $3, $4, $5"
     // 6 PlaceNums_wFV      "$1, $2, $3, $4, $5, $6"
     // 7 PlaceNums_wID_wFV  "$1, $2, $3, $4, $5, $6, $7" // 6,7=K,V 
     // 4 UpdateNames        "f0=$1, f1=$2, f2=$3, f3=$4" // No ID!

     pTD.CSVs.PlaceNums_noID    = csvNumbers("$", 1, N-1) 
     pTD.CSVs.PlaceNums_wID     = csvNumbers("$", 1, N) 
     pTD.CSVs.PlaceNums_wID_wFV = csvNumbers("$", 1, N+2)
     pTD.CSVs.FieldNames_noID   = csvNames(CSs, 0, N-2)
     pTD.CSVs.FieldNames_wID    = pTD.PKname + ", " + csvNames(CSs, 0, N-2)

     // For clarity in composability: No semicolons!
     s := pTD.PKname 
     pTD.CSVs.Where_noVals     = fmt.Sprintf(" WHERE %s = $%d", s, 1)
     pTD.CSVs.Where_wVals_wID  = fmt.Sprintf(" WHERE %s = $%d", s, N+1)
     pTD.CSVs.Where_wVals_noID = fmt.Sprintf(" WHERE %s = $%d", s, N)

     // For later composability: No placeholders for UPDATE's WHERE!
     var sbUpdtNams S.Builder
     var i int
     for i=0; i<N-1; i++ {
     	 sbUpdtNams.WriteString(
		CSs[i].StorName + " = $" +
	 	strconv.Itoa(i+1) + ", ")
	}
     pTD.CSVs.UpdateNames = S.TrimSuffix(sbUpdtNams.String(), ", ")

     println(pTD.TableSummary.StorName +
     	"[" + strconv.Itoa(len(pTD.ColumnSpecs)) + "]:\n" + pTD.CSVs.String())

     return nil
}

// GenerateStatements generates struct [Statements] for
// every struct that has been registered using method
// datarepo/RegisterAppTables of interface datarepo/SimpleRepo
//
// The examples in the inline docu here in this source code
// file assume five (5) fields:
//  - F-numbers (field names) start at "F0"
//  - At element [0] is "F0", the table's primary key
//    (or: the "ID column"), denoted in docu as "F0=ID"
//  - Elements [1..4] are names fields/columns  F1, F2, F3, F4
//  - placeholder-numbers start at "$1", Postgres-style
//
// This func does not modify or even access the DB. It should
// be called ASAP after program start, but not until AFTER
// func [GenerateColumnStringsCSV] is called. It does not
// need to be called before a DB is opened, but it DOES
// need to be called any SQL is executed against the DB.
// .
func GenerateStatements(pTD *TableDetails) error {
     if pTD.Stmts != nil {
     	println("GenerateStatements: DUPE CALL")
	return errors.New("GenerateStatements: dupe initialisation")
	}
     pTD.Stmts = new(Statements)

     // === INSERT ===========================================
     // Add, Create, Insert, New
     // Use INSERT...RETURNING to get the new ID. 
     // https://www.sqlite.org/lang_insert.html
     // INSERT INTO tblNm (F1, F2, F3, F4) VALUES(val1, val2, val3, val4);-
     // INSERT INTO tblNm (F1, F2, F3, F4) VALUES($1,$2,$3,$4); + any...
     // ======================================================
     // table name + column names CSV + placeholders CSV 
     // No F0=ID. No WHERE clause. 
     // FIELDS are FieldNames_noID.
     // VALUES are PlaceNums_noID.
     // ======================================================
     pTD.Stmts.INSERTuniqueID =
	"INSERT INTO " + pTD.TableSummary.StorName +
        	   "(" + pTD.CSVs.FieldNames_noID  + ") " +
             "VALUES(" + pTD.CSVs.PlaceNums_noID   + ") " +
       /* "RETURNING " + pTD.PKname                + */ ";"

     // === SELECT ===========================================
     // Fetch, Get, List, Retrieve, Select
     // https://www.sqlite.org/lang_select.html
     // SELECT F0, F1, F2, F3, F4 FROM tblNm WHERE expr; 
     // https://www.sqlite.org/syntax/expr.html
     // ======================================================
     // table name + column names CSV + WHERE clause
     // With F0=ID. With WHERE clause.
     // FIELDS are FieldNames_wID. 
     // The WHERE clause is tipicly on the ID but need not be. 
     // ======================================================
     pTD.Stmts.SELECTuniqueID =
	"SELECT " + pTD.CSVs.FieldNames_wID +
	" FROM "  + pTD.TableSummary.StorName + pTD.CSVs.Where_noVals + ";"

     // === UPDATE ===========================================
     // Modify, Update
     // https://www.sqlite.org/lang_update.html
     // Obnoxious syntax:
     // UPDATE tblNm SET fld1=val1,fld2=val2 WHERE expr:
     // (or..) SET fld1=$1, fld2=$2 WHERE expr; + any...
     // https://www.sqlite.org/syntax/expr.html
     // Use UpdateNames.
     // -----------------------------------------------------
     // For UPDATE (only), we have to generate here+now an
     // SQL string that involves all columns (except the ID).
     // Write assignment pairs as CSV: f1 = $1, f2 = $2, ...
     // We do NOT include the primary key, D.SFT_PRKEY, which
     // is used in the WHERE.
     // We do NOT include the WHERE clause.
     // -----------------------------------------------------
     pTD.Stmts.UPDATEuniqueID =
     	"UPDATE " + pTD.TableSummary.StorName +
	" SET "   + pTD.CSVs.UpdateNames + pTD.CSVs.Where_wVals_noID + ";"
     
     // === DELETE ============================
     // Delete, Discard, Drop
     // https://www.sqlite.org/lang_delete.html
     // =======================================
     pTD.Stmts.DELETEuniqueID =
     	"DELETE FROM " + pTD.TableSummary.StorName +
	" WHERE " + pTD.PKname + " = $1" + ";"

     return nil
}

/*     
     println("FieldNames    ", pTD.CSVs.FieldNames)
     println("FieldNames wID", pTD.CSVs.FieldNames_wID)
     println("PlaceNmbrs    ", pTD.CSVs.PlaceNumbers)
     println("PlaceNmbrs wID", pTD.CSVs.PlaceNrs_wID)
     println("UpdateNames   ", pTD.CSVs.UpdateNames)
* /
     return nil
}


*/