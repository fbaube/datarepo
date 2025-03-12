package rowmodels

import(
	// "fmt"
	"errors"
	"strconv"
	S "strings"
	D "github.com/fbaube/dsmnd"
)

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
	// '$'-numbered parameters (like Postgres): "$1, $2, $3" 
	PlaceNumbers,   PlaceNrs_wID    string
	// FieldUpdates [NOT with primarykey] is a list of
	// column/field names with "=", and the values as
	// '$'-numbered parameters: "F1 = $1, F2 = $2, F3 = $3" 
	UpdateNames     string
}

// PreparedStatements stores several instances of [sql.Stmt]
// prepared statements customised for the table. Statements
// vary in whether they include the primary key, and whether
// they include a WHERE clause.
//
// Statements named "*unique" are for working with single 
// records, and are used by method [datarepo.EngineUnique]
// of interface [datarepo.DBEnginer].
// .
type PreparedStatements struct {
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
     // type ColumnStringsCSV struct {
     //   FieldNames,     FieldNames_wID  string  "F1, F2, F3" 
     //   PlaceNumbers,   PlaceNrs_wID    string  "$1, $2, $3" 
     //   UpdateNames,    UpdateNames_wID string  "F1=$1,F2=$2,F3=$3" 
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
*/
     return nil
}

// GeneratePreparedStatements generates struct [PreparedStatements] 
// for every struct that has been registered using method
// datarepo/RegisterAppTables of interface datarepo/SimpleRepo .
//
// TODO: It should be guarded by a Do.Once()
//
// It does not modify or even access the database. It
// should be called ASAP after program start, but AFTER
// GenerateColumnStringsCSV is called. It does not need to 
// be called before a database is opened, but it DOES need 
// to be called any SQL is executed against the database.
// .
func GeneratePreparedStatements(pTD *TableDetails) error {
// FIXME
     if pTD.CSVs != nil {
     	println("GeneratePreparedStatements: DUPE CALL")
	return errors.New("GeneratePreparedStatements: dupe initialisation")
	}    
     pTD.CSVs = new(ColumnStringsCSV)
     var colSpex []D.ColumnSpec
     colSpex   = pTD.ColumnSpecs
     // type ColumnStringsCSV struct {
     //   FieldNames,     FieldNames_wID  string  "F1, F2, F3" 
     //   PlaceNumbers,   PlaceNrs_wID    string  "$1, $2, $3" 
     //   UpdateNames,    UpdateNames_wID string  "F1=$1,F2=$2,F3=$3" 
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
*/
     return nil
}

