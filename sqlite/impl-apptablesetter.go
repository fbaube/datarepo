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

// theMap holds all table schemata passed to RegisterAppTables(..).
// Map key is "appname_tablename", where appname is forced to
// lower case and tablename is taken from the TableDetails.
// Map key is simply "tablename" if appname was "".
//
// Map value is the associated instance of TableDetails.
// .
var theMap map[string]*DRM.TableDetails

func init() {
	theMap = make(map[string]*DRM.TableDetails)
}

// Interface AppTableSetter is table-related methods for a specified
// app's schema. The app name is case-insensitive, and used as all
// lower case, and preixed to table names as "appname_". If the app
// name is left blank (""), a default namespace is used and no prefix
// is added to table names.

// RegisterAppTables processes the schemata of the specified app's
// tables, which this interface creates and/or manages. Multiple
// calls, whether with tables previously specified or not before
// seen do not conflict; if a table name is repeated but with a
// different schema, the result is undefined.
// .
func (p *SqliteRepo) RegisterAppTables(appName string, cfg []*DRM.TableDetails) error {
	L.L.Info("RegisterAppTables: got %d table definitions", len(cfg))
	var c *DRM.TableDetails
	for _, c = range cfg {
	       	sindex := S.ToLower(c.DispName) // c.StorName)
		// println("REG TBL DTLS: " + sindex)
		theMap[sindex] = c
		L.L.Info("Reg'd the config for app table: " +
			S.ToLower(c.StorName))
		// Do schema-related initialisations
		_ = DRM.GenerateColumnStringsCSV(c)
	}
	return nil
}

// EmptyAllTables deletes (app-level) data from the app's tables
// but does not delete any tables (i.e. no DROP TABLE are done).
// The DB should be open when it is called (so that the connection
// object exists). The DB should have a path, but mainly just for
// error messages; the requirement could be removed.
//
// NOTE: If a table does not exist, it has to be created.
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
	for _, c := range theMap {
		CTS := "DELETE FROM " + S.ToLower(c.StorName) + ";"
		// p.Handle().MustExec(CTS)
		_, err := p.Exec(CTS)
		if err != nil {
			strerr := err.Error()
			if S.HasPrefix(strerr, "no such table:") {
				L.L.Info("No such table: " + c.StorName)
			// OOPS! Create it!
			e2 := p.createAppTable(c)
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
		} else {
			L.L.Info("Deleted all from table: " + S.ToLower(c.StorName))
		}
	}
	L.L.Info("SQLAR not emptied, utils/repo/sqlite/impl_apptables.go L83")
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
	for _, td := range theMap {
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
		return fmt.Errorf("CreAppTbl failed: %w", err)
	}
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

func GetTableDetailsByCode(s string) *DRM.TableDetails {
     s = S.ToLower(s) 
  // println("GetTableDetailsByCode: " + s)
  // return theMap[S.ToLower(s)]
     ret := theMap[s]
     if ret == nil {
     	for k, v := range theMap {
	    fmt.Printf("MAP: %+v %+v \n", k, *v) }
	    }
     return ret 
}