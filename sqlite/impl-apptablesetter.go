package sqlite

import (
	_ "database/sql" // to get init()
	"errors"
	"fmt"
	L "github.com/fbaube/mlog"
	DRM "github.com/fbaube/datarepo/rowmodels"
	"io/ioutil"
	S "strings"
)

// mapByDispName & mapByStorName holds all table schemata passed
// to RegisterAppTables(..). Keys are forced to lower case. 
// Map key is "appname_tablename", where appname is forced to
// lower case and tablename is taken from the TableDetails.
// Map key is simply "tablename" if appname was "".
//
// Map value is the associated instance of TableDetails.
// .
var mapByDispName, mapByStorName map[string]*DRM.TableDetails

func init() {
	mapByDispName = make(map[string]*DRM.TableDetails)
	mapByStorName = make(map[string]*DRM.TableDetails)
}

// Interface AppTableSetter is table-related methods for a specified
// app's schema. Specifying a non-empty string is untested !
//
// The app name is case-insensitive, and used as all lower case, and
// prefixed to table names as "appname_". If the app name is left
// blank (""), a default namespace is used and no prefix is added
// to table names.

// RegisterAppTables processes the schemata of the specified 
// app's tables, which this interface creates and/or manages. 
// This includes filling many fields in struct [TableDetails].
// 
// It registers them in a singleton map in package datarepo/sqlite.
// However this func is left as a method of a SimpleRepo so that 
// future corrections are easier. 
//
// Multiple calls to this do not conflict, whether with tables
// previously specified or not before; if a table name is repeated 
// but with a different schema, the result is undefined.
//
//  - StorName: the name of the table IN THE DB - a "long name"
//  - DispName: a short name (3 ltrs!) for use in building up names 
// .
func (p *SqliteRepo) RegisterAppTables(appName string, pTDs []*DRM.TableDetails) error {
	L.L.Info("RegisterAppTables: got %d table definitions", len(pTDs))
	var td *DRM.TableDetails
	var lcDN, lcSN string
	var i int
	for i, td = range pTDs {
	       	lcDN = S.ToLower(td.DispName) 
		lcSN = S.ToLower(td.StorName)
		println("REG TBL DTLS:", lcDN, lcSN)
		mapByDispName[lcDN] = td
		mapByStorName[lcSN] = td
		L.L.Info("Reg'd config for app table [%d]: %s/%s",
			i, lcDN, lcSN)
		// Do schema-related initialisations
		_ = DRM.GenerateColumnStringsCSV(td)
		_ = DRM.GenerateStatements(td)
	}
	fmt.Printf("MAPS lens stor<%d> disp<%d> \n",
		len(mapByStorName), len(mapByDispName))
	return nil
} 

// EmptyAllTables deletes (app-level) data from the app's tables
// but does not delete any tables (i.e. no DROP TABLE are done).
// The DB should be open when it is called (so that the connection
// object exists). The DB should have a path, but mainly just for
// error messages; the requirement could be removed.
//
// NOTE: If a table does not exist, it has to be created.
//
// NOTE: Errors appeared here because each table was in map twice,
// under both StorNam and DispName, so that was rolled back. 
// .
func (p *SqliteRepo) EmptyAppTables() error {
	if p.Path() == "" {
		return errors.New("sqliterepo.emptyapptables: no path")
	}
	if !p.IsOpen() { // p.Handle() == nil {
		return fmt.Errorf("sqliterepo.emptyapptables(%s): "+
			"not open", p.Path())
	}
	var e error
	for _, td := range mapByStorName {
		CTS := "DELETE FROM " + S.ToLower(td.StorName) + ";"
		// p.Handle().MustExec(CTS)
		_, err := p.Exec(CTS)
		if err == nil {
			L.L.Info("Deleted all from table: " +
				S.ToLower(td.StorName))
			continue
		}
		strerr := err.Error()
		if S.HasPrefix(strerr, "no such table:") {
			L.L.Info("No such table: " + td.StorName)
		// OOPS! Create it!
		e2 := p.createAppTable(td)
		if e2 != nil {
		   return fmt.Errorf("EmptyAppTbls.CreTbl failed: %w", e2)
		// panic(e2.Error())
		   }
		} else {
			L.L.Error("reposqlite.emptyAllTbls: " + strerr)
			return fmt.Errorf(
				"sqliterepo.emptyAppTbls(%s) "+
					"failed: %w", p.Path(), e)
		}
	}
	L.L.Info("emptyapptables: not SQLAR, sqlite/impl-apptablesetter L114")
	if e != nil {
		return fmt.Errorf(
			"sqliterepo.emptyapptables(%s) failed: %w", p.Path(), e)
	}
	return nil
}

// CreateAppTables creates the app's tables per already-supplied
// schema(ta); if the tables exist, they are emptied of data.
// It uses our simplified SQLite DB model wherein
//   - Every column is either string ("TEXT") or int ("INTEGER"),
//   - Every column is NOT NULL (because NULL is evil),
//   - Every column has type checking (TBS), and
//   - Every table has a primary index field, and
//   - Every index (both primary and foreign) includes the full name of the
//     table, which simplifies column creation and table cross-referencing
//     (and in particular, JOINs).
//
// .
func (p *SqliteRepo) CreateAppTables() error {
	// println("CreateAppTables")
	// func (pDB SqliteRepo) CreateTable_sqlite(ts U.TableDetails) error {
	// FIXME Check table name prefix (e.g. "mmmc_") ?
	for _, td := range mapByStorName {
		e := p.createAppTable(td)
		if e != nil {
			// FIXME
			return e
		}
	}
	// L.L.Warning("SQLAR is To-Do, datarepo/sqlite/impl_apptables.go L123")
	return nil
}

func (p *SqliteRepo) createAppTable(td *DRM.TableDetails) error {

	var CTS string // the Create Table SQL string
	var e error

	CTS, e = p.NewCreateTableStmt(td)
	if e != nil {
		return fmt.Errorf(
		     "Cannot create app table: %s: %w", td.StorName, e)
	}
	fnam := "./create-table-" + td.StorName + ".sql"
	e = ioutil.WriteFile(fnam, []byte(CTS), 0644)
	if e != nil {
		L.L.Error("Could not write file: " + fnam)
	} else {
		L.L.Debug("Wrote \"CREATE TABLE " +
			td.StorName + " ... \" to: " + fnam)
	}
	// p.Handle().MustExec(CTS)
	_, err := p.Exec(CTS)
	if err != nil {
		// panic(err)
		return fmt.Errorf("CreAppTbl(%s) failed: %w", td.StorName, err)
	}
	fmt.Printf("CreAppTbl(%s) OK! \n", td.StorName)
	/*
	ss, e := p.DumpTableSchema_sqlite(td.StorName)
	if e != nil {
		return fmt.Errorf("simplerepo.createtable.sqlite: "+
			"DumpTableSchema<%s> failed: %w", e)
	}
	L.L.Debug(td.StorName + " SCHEMA: " + ss)
	*/
	// println("TODO: (maybe) Insert record with IDX 0 and string descr's")
	//    and ("TODO: (then) Dump all table records (i.e. just one)")
	return nil
}

func GetTableDetailsByStorName(s string) *DRM.TableDetails {
     s = S.ToLower(s) 
  // println("GetTableDetailsByString: " + s)
     ret, ok := mapByStorName[s]
     if !ok { // ret == nil {
	fmt.Printf("MAP len stor<%d> \n", len(mapByStorName))
     	for k, v := range mapByStorName {
	    fmt.Printf("\nMAP StorName <%s>: %s %s \n", s, k, v.DispName)
	    }
	}
     return ret 
}

func GetTableDetailsByDispName(s string) *DRM.TableDetails {
     s = S.ToLower(s) 
  // println("GetTableDetailsByString: " + s)
     ret, ok := mapByDispName[s]
     if !ok { // ret == nil {
	fmt.Printf("MAP len disp<%d> \n", len(mapByDispName))
     	for k, v := range mapByDispName {
	    fmt.Printf("\nMAP DispName <%s>: %s %s \n", s, k, v.StorName) 
	    }
	}
     return ret 
}

