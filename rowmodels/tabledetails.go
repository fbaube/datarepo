package rowmodels

import(
	// "fmt"
	"errors"
	"strconv"
	S "strings"
	D "github.com/fbaube/dsmnd"
)

// Column Field Key Prop / Value 

// ColumnStringsCSV stores strings useful for composing SQL
// statements. Each string includes all the columns, in order,
// comma-separated. SQL using these strings defaults to setting
// and getting every field in a DB record.
//
// The strings have no trailing commas. Each string (except 
// for UPDATE) has a "full" version (suffixed with "_wID")
// that includes the primary key (always named "{table}_ID") 
// for output from SELECT, and (importantly!) a version 
// withOUT the "{table}_ID" primary key, for input to INSERT 
// (where the ID is new) and input to UPDATE (where the ID
// finds the record).
// .
type ColumnStringsCSV struct {
	// FieldNames   [+withID primarykey] is a list of 
	// column (i.e. field) names, in order: "F1, F2, F3" 
	FieldNames,     FieldNames_wID  string
	// PlaceNumbers [+withID primarykey] is a list of 
	// '$'-numbered parameters (like Postgres): "$1, $2, $3",
	// and with two extra when using "WHERE Field = Value" 
	PlaceNumbers,   PlaceNrs_wID, PlaceNrs_wID_wFV  string
	// FieldUpdates [NOT with primarykey] is a list of
	// column/field names with "=", and the values as
	// '$'-numbered parameters: "F1 = $1, F2 = $2, F3 = $3" 
	UpdateNames     string
}

// Statements stores several SQL query strings customised for the
// table. Statements vary in whether they include the primary key, 
// and whether they include a WHERE clause.
//
// Statements named "*unique" are for working with single records,
// and are used by method [datarepo.EngineUnique] of interface
// [datarepo.DBEnginer].
// .
type Statements struct {
     	INSERTunique string
	SELECTunique string
	UPDATEunique string
	DELETEunique string 
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
     var colSpex []D.ColumnSpec
     colSpex   = pTD.ColumnSpecs

     // Examples, if we assume 5 fields (F1=ID plus four others)
     // type struct { id, f2, f3, f4, f5 } 
     // type ColumnStringsCSV struct {
     // 4 FieldNames             "f2, f3, f4, f5" 
     // 4 PlaceNumbers           "$1, $2, $3, $4" 
     // 5 FieldNames_wID     "id, f2, f3, f4, f5" 
     // 5 PlaceNrs_wID       "$1, $2, $3, $4, $5"
     // 6 PlaceNrs_wFV       "$1, $2, $3, $4, $5, $6"
     // 7 PlaceNrs_wID_wFV   "$1, $2, $3, $4, $5, $6, $7" // 6,7=K,V // SELECT?
     // 4 UpdateNames        "f2=$1, f3=$2, f4=$3, f5=$4" // No ID! 
     
     var sbFN, sbPN, sbUN S.Builder
     var columnName, placeNumber string
     var iCol, nCols int
     nCols = len(pTD.ColumnSpecs)
     // Now make lots of small, simple loops
     for iCol = 0; iCol < nCols; iCol++ { 
     }

     for i, pCS := range colSpex {
	columnName  = pCS.StorName
	placeNumber = strconv.Itoa(i+1)
	if i == 0 && !S.HasSuffix(columnName, "ID") {
	}
     	// fmt.Printf("%s[%d]=%s: %s \n",
	// 	pTD.StorName, i, columnName, pCS.String())
	sbFN.WriteString(columnName  + ", ")
	sbPN.WriteString("$" + placeNumber + ", ")
	sbUN.WriteString(columnName  + " = $" + placeNumber + ", ")
     }
     sFN := S.TrimSuffix(sbFN.String(), ", ")
     sPN := sbPN.String()
     sUN := S.TrimSuffix(sbUN.String(), ", ")
     pTD.CSVs.FieldNames   = sFN
     pTD.CSVs.PlaceNumbers = sPN
     pTD.CSVs.UpdateNames  = sUN
     pTD.CSVs.FieldNames_wID = pTD.PKname + ", " + sFN
     pTD.CSVs.PlaceNrs_wID   = sPN + "$" + strconv.Itoa(len(colSpex)+1)
     pTD.CSVs.PlaceNumbers   = S.TrimSuffix(pTD.CSVs.PlaceNumbers, ", ")
/*     
     println("FieldNames    ", pTD.CSVs.FieldNames)
     println("FieldNames wID", pTD.CSVs.FieldNames_wID)
     println("PlaceNmbrs    ", pTD.CSVs.PlaceNumbers)
     println("PlaceNmbrs wID", pTD.CSVs.PlaceNrs_wID)
     println("UpdateNames   ", pTD.CSVs.UpdateNames)
*/
     return nil
}

// GeneratePreparedStatements generates struct [PreparedStatements] 
// for every struct that has been registered using method
// datarepo/RegisterAppTables of interface datarepo/SimpleRepo .
//
// TODO: It should be guarded by a Do.Once()
//
// It does not modify or even access the databaseDB. It should be called
// ASAP after program start, but AFTER [GenerateColumnStringsCSV] is 
// called. It does not need to be called before a DB is opened, but
// it DOES need to be called any SQL is executed against the DB.
// .
func GenerateStatements(pTD *TableDetails) error {
     if pTD.Stmts != nil {
     	println("GenerateStatements: DUPE CALL")
	return errors.New("GenerateStatements: dupe initialisation")
	}    
     pTD.Stmts = new(Statements)
     // var colSpex []D.ColumnSpec
     // colSpex   = pTD.ColumnSpecs
     
     // type ColumnStringsCSV struct {
     //      FieldNames,   FieldNames_wID  string  "F1, F2, F3" 
     //      PlaceNumbers, PlaceNrs_wID    string  "$1, $2, $3" 
     //      UpdateNames                   string  "F1=$1,F2=$2,F3=$3"
     
     // === INSERT ===========================================
     // Add, Create, Insert, New
     // Use RETURNING to get new ID. 
     // https://www.sqlite.org/lang_insert.html
     // INSERT INTO tblNm (fld1, fld2) VALUES(val1, val2);
     // INSERT INTO tblNm (fld1, fld2) VALUES($1,$2); + any...
     // FIELDS are FieldNames[_wID]. VALUES are PlaceNrs[_wID].
     // ======================================================
     // table name + column names CSV + placeholders CSV 
     // WithOUT ID (primary key D.SFT_PRKEY). No WHERE clause. 
     // ======================================================
     pTD.Stmts.INSERTunique =
	"INSERT INTO " + pTD.TableSummary.StorName +
        	   "(" + pTD.CSVs.FieldNames   + ") " +
             "VALUES(" + pTD.CSVs.PlaceNumbers + ") " +
          "RETURNING " + pTD.PKname            + ";"

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
     // TODO Should WHERE use '$'-placeholders ?? 
     pTD.Stmts.SELECTunique =
	"SELECT " + pTD.CSVs.FieldNames_wID +
	" FROM "  + pTD.TableSummary.StorName +
     //	" WHERE " + pFV.Field + " = " + pFV.Value + ";"
	" WHERE $FIXME1 = $FIXME2;" // FIXME FIXME FIXME FIXME FIXME 

     return nil
}

/*
// WITH WHERE and withOUT WHERE
func buildSELECT(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }
func buildUPDATE(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }
func buildDELETE(pTD *DRM.TableDetails, pFV DRP.FieldValuePair) { }
* /

     var sbFN, sbPN, sbUN S.Builder
     var columnName, placeNumber string
     for i, pCS := range colSpex {
	columnName  = pCS.StorName
	placeNumber = strconv.Itoa(i+1)
     	// fmt.Printf("%s[%d]=%s: %s \n",
	// 	pTD.StorName, i, columnName, pCS.String())
	sbFN.WriteString(columnName  + ", ")
	sbPN.WriteString("$" + placeNumber + ", ")
	sbUN.WriteString(columnName  + " = $" + placeNumber + ", ")
     }
     sFN := S.TrimSuffix(sbFN.String(), ", ")
     sPN := sbPN.String()
     sUN := S.TrimSuffix(sbUN.String(), ", ")
     pTD.CSVs.FieldNames   = sFN
     pTD.CSVs.PlaceNumbers = sPN
     pTD.CSVs.UpdateNames  = sUN
     pTD.CSVs.FieldNames_wID = pTD.PKname + ", " + sFN
     pTD.CSVs.PlaceNrs_wID   = sPN + "$" + strconv.Itoa(len(colSpex)+1)
     pTD.CSVs.PlaceNumbers   = S.TrimSuffix(pTD.CSVs.PlaceNumbers, ", ")
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