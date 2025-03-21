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
     var N = len(pTD.ColumnSpecs)
     // var cs D.ColumnSpec 

     // Examples, if we assume N=5 fields (f[0]=ID plus four others, f[1..4])
     // type struct { id, f1, f2, f3, f4 } 
     // type ColumnStringsCSV struct {
     // 4 FieldNames_noID    "f1, f2, f3, f4" 
     // 4 PlaceNumbers_noID  "$1, $2, $3, $4" 
     // 5 FieldNames_wID     "id, f1, f2, f3, f4" 
     // 5 PlaceNrs_wID       "$1, $2, $3, $4, $5"
     // 6 PlaceNrs_wFV       "$1, $2, $3, $4, $5, $6"
     // 7 PlaceNrs_wID_wFV   "$1, $2, $3, $4, $5, $6, $7" // 6,7=K,V // SELECT?
     // 4 UpdateNames        "f1=$1, f2=$2, f3=$3, f4=$4" // No ID!

     pTD.CSVs.PlaceNums_noID    = csvNumbers("$", 1, N-1) 
     pTD.CSVs.PlaceNums_wID     = csvNumbers("$", 1, N) 
     pTD.CSVs.PlaceNums_wID_wFV = csvNumbers("$", 1, N+2)
     pTD.CSVs.FieldNames_noID   = csvNames(CSs, 1, N-1)
     pTD.CSVs.FieldNames_wID    = csvNames(CSs, 0, N-1)

     // For clarity in composability: No semicolons! 
     pTD.CSVs.Where_wID  = fmt.Sprintf(" WHERE $%d = $%d", N, N+1)
     pTD.CSVs.Where_noID = fmt.Sprintf(" WHERE $%d = $%d", N-1, N)

     // For later composability: No placeholders for UPDATE's WHERE!
     var sbUpdtNams S.Builder
     var i int
     for i=1; i<N; i++ {
     	 sbUpdtNams.WriteString(
		CSs[i].StorName + " = $" +
	 	strconv.Itoa(i) + ", ")
	}
     pTD.CSVs.UpdateNames = S.TrimSuffix(sbUpdtNams.String(), ", ")

     // So if (for example) we have N fields:
     // - The first is the ID, and there are N-1 others
     // - The N  fields are indexed as 0..N-1, or in slice notation,  [:N]
     // - Non-ID fields are indexed as 1..N-1, or in slice notation, [1:N]
     // - Including the ID, the placeholders are numbered $1..$N
     // - Not including it, the placeholders are numbered $2..$N

     println("FieldNames noID", pTD.CSVs.FieldNames_noID)
     println("FieldNames  wID", pTD.CSVs.FieldNames_wID)
     println("PlaceNmbrs noID", pTD.CSVs.PlaceNums_noID)
     println("PlaceNmbrs  wID", pTD.CSVs.PlaceNums_wID)
     println("UpdateNames    ", pTD.CSVs.UpdateNames)
     println("Where noID", pTD.CSVs.Where_noID)
     println("Where  wID", pTD.CSVs.Where_wID)

     return nil
}

// GeneratePreparedStatements generates struct [Statements] 
// for every struct that has been registered using method
// datarepo/RegisterAppTables of interface datarepo/SimpleRepo .
//
// TODO: It should maybe be guarded by a Do.Once()
//
// NOTE: f-numbers (field names) start at "f0", and f0 is always the ID 
// (primary key). placeholder-numbers start at "$1", Postgres-style. 
//
// It does not modify or even access the DB. It should be called
// ASAP after program start, but AFTER [GenerateColumnStringsCSV]  
// is called. It does not need to be called before a DB is opened, 
// but it DOES need to be called any SQL is executed against the DB.
// .
func GenerateStatements(pTD *TableDetails) error {
     if pTD.Stmts != nil {
     	println("GenerateStatements: DUPE CALL")
	return errors.New("GenerateStatements: dupe initialisation")
	}    
     pTD.Stmts = new(Statements)
     
     // === INSERT ===========================================
     // Add, Create, Insert, New
     // Use RETURNING to get new ID. 
     // https://www.sqlite.org/lang_insert.html
     // INSERT INTO tblNm (fld1, fld2) VALUES(val1, val2);
     // INSERT INTO tblNm (fld1, fld2) VALUES($1,$2); + any...
     // ======================================================
     // table name + column names CSV + placeholders CSV 
     // FIELDS are FieldNames_noID. VALUES are PlaceNums_noID.
     // WithOUT ID (primary key D.SFT_PRKEY). No WHERE clause. 
     // ======================================================
     pTD.Stmts.INSERTunique =
	"INSERT INTO " + pTD.TableSummary.StorName +
        	   "(" + pTD.CSVs.FieldNames_noID  + ") " +
             "VALUES(" + pTD.CSVs.PlaceNums_noID   + ") " +
          "RETURNING " + pTD.PKname                + ";"

     // === SELECT ===========================================
     // Fetch, Get, List, Retrieve, Select
     // https://www.sqlite.org/lang_select.html
     // SELECT fld1, fld2 FROM tblNm WHERE expr
     // https://www.sqlite.org/syntax/expr.html
     // FIELDS are FieldNames_wID. 
     // ======================================================
     // table name + column names CSV + WHERE clause
     // WITH ID (primary key D.SFT_PRKEY). 
     // The WHERE clause is tipicly on the ID but need not be. 
     // ======================================================
     pTD.Stmts.SELECTunique =
	"SELECT " + pTD.CSVs.FieldNames_wID +
	" FROM "  + pTD.TableSummary.StorName + pTD.CSVs.Where_wID + ";"


     return nil
}

/*
// WITH WHERE and withOUT WHERE
func buildSELECT(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }
func buildUPDATE(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }
func buildDELETE(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }
* /

/*     
     println("FieldNames    ", pTD.CSVs.FieldNames)
     println("FieldNames wID", pTD.CSVs.FieldNames_wID)
     println("PlaceNmbrs    ", pTD.CSVs.PlaceNumbers)
     println("PlaceNmbrs wID", pTD.CSVs.PlaceNrs_wID)
     println("UpdateNames   ", pTD.CSVs.UpdateNames)
* /
     return nil
}

===

     	// ================================================
	case "M", "U": // Modify, Update
	// https://www.sqlite.org/lang_update.html
	// Obnoxious syntax: 
	// UPDATE tblNm SET fld1=val1,fld2=val2 WHERE expr: 
	// (or..) SET fld1=$1, fld2=$2 WHERE expr; + any...
	// https://www.sqlite.org/syntax/expr.html
	// Use UpdateNames. 
     	// ================================================
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

     	// =======================================
	case "D": // Delete, Discard, Drop
	// https://www.sqlite.org/lang_delete.html
	// DELETE FROM tblNm WHERE expr RET'G expr
     	// =======================================

*/