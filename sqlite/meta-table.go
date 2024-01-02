package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	D "github.com/fbaube/dsmnd"
	S "strings"
	// DU "github.com/fbaube/dbutils"
	// L "github.com/fbaube/mlog"
)

// DbTblColsInDb returns all column names & types
// for the specified table as found in the DB.
// .
func (pDB SqliteRepo) DbTblColsInDb(tableName string) ([]*D.DbColInDb, error) {
	if tableName == "" {
		return nil, nil
	}
	switch pDB.Type() {
	case D.DB_SQLite:
		// Is OK!
	default:
		return nil, errors.New(
			"simplerepo.sqlite.dbtblcolsindb: not sqlite")
	}
	var e error
	var Rs *sql.Rows
	var CTs []*sql.ColumnType
	var retval []*D.DbColInDb

	Rs, e = pDB.Handle().Query("select * from " + tableName + " limit 1")
	if e != nil {
		return nil, fmt.Errorf("simplerepo.dbtblcolsindb: "+
			"select * : failed on table <%s>: %w", tableName, e)
	}
	CTs, e = Rs.ColumnTypes()
	if e != nil {
		return nil, fmt.Errorf("simplerepo.dbtblcolsindb: "+
			"rs.ColumnTypes failed on table <%s>: %w", tableName, e)
	}
	for _, ct := range CTs {
		dci := new(D.DbColInDb)
		dci.Datatype = D.Datatype(ct.DatabaseTypeName())
		dci.StorName = ct.Name()
		retval = append(retval, dci)
	}
	return retval, nil
}

// DumpTableSchema_sqlite returns the (SQLite) schema
// for the specified table as found in the DB.
// .
func (pDB SqliteRepo) DumpTableSchema_sqlite(tableName string) (string, error) {
	switch pDB.Type() {
	case D.DB_SQLite:
		// Is OK!
	default:
		return "", errors.New(
			"reposqlite.dumptableschema: not sqlite")
	}
	var theCols []*D.DbColInDb
	var sb S.Builder
	var sType string
	var e error

	theCols, e = pDB.DbTblColsInDb(tableName)
	if e != nil {
		return "", fmt.Errorf("simplerepo.dumptableschema.sqlite: "+
			"pDB.DbTblColsInDb<%s> failed: %w", e)
	}
	for i, c := range theCols {
		sType = ""
		if D.BasicDatatype(c.Datatype) != D.BDT_TEXT {
			sType = "(" + string(c.Datatype) + "!)"
		}
		sb.WriteString(fmt.Sprintf(
			"[%d]%s%s / ", i, sType, c.StorName))
	}
	sb.WriteString(fmt.Sprintf("%d fields", len(theCols)))
	return sb.String(), nil
}
