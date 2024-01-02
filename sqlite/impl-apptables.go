package sqlite

import (
	_ "database/sql" // to get init()
	"errors"
	"fmt"
	L "github.com/fbaube/mlog"
	DRU "github.com/fbaube/datarepo/utils"
	"io/ioutil"
	S "strings"
)

// theMap holds all table schemata passed to SetAppTables(..).
// Map key is "appname_tablename", where appname is forced to
// lower case and tablename is taken from the TableDescriptor.
// Map key is simply "tablename" if appname was "".
//
// Map value is the associated instance of TableDescriptor.
// .
var theMap map[string]DRU.TableDescriptor

func init() {
	theMap = make(map[string]DRU.TableDescriptor)
}

// RepoAppTables is table-related methods for a specified app's schema.
// The app name is case-insensitive, and used as all lower case, and
// prefixed to table names as "appname_". If the app name is left blank
// (""), a default namespace is used and no prefix is added to table names.
//
// type RepoAppTables interface {

// SetAppTables specifies the schemata of the specified app's
// tables, which this interface creates and/or manages. Multiple
// calls, whether with tables previously specified or not before
// seen do not conflict; if a table name is repeated but with a
// different schema, the result is undefined.
// .
func (p *SqliteRepo) SetAppTables(appName string, cfg []DRU.TableDescriptor) error {
	L.L.Info("SetAppTables: got %d table definitions", len(cfg))
	for _, c := range cfg {
		theMap[ /*pfx+*/ S.ToLower(c.Name)] = c
		L.L.Info("Reg'd the config for app table: " +
			S.ToLower(c.Name))
	}
	return nil

}

// EmptyAllTables deletes (app-level) data from the app's tables
// but does not delete any tables (i.e. no DROP TABLE are done).
// The DB should be open when it is called (so that the connection
// object exists). The DB should have a path, but mainly just for
// error messages; the requirement could be removed.
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
		CTS := "DELETE FROM " + S.ToLower(c.Name) + ";"
		// p.Handle().MustExec(CTS)
		_, err := p.Exec(CTS)
		if err != nil {
			strerr := err.Error()
			if S.HasPrefix(strerr, "no such table:") {
				L.L.Info("No such table: " + c.Name)
			} else {
				L.L.Error("reposqlite.emptyAllTbls: " + strerr)
				return fmt.Errorf(
					"sqliterepo.emptyAppTbls(%s) "+
						"failed: %w", p.Path(), e)
			}
		} else {
			L.L.Info("Deleted all from table: " + S.ToLower(c.Name))
		}
	}
	L.L.Warning("SQLAR not emptied, utils/repo/sqlite/impl_apptables.go L83")
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
	println("CreateAppTables")
	// func (pDB SqliteRepo) CreateTable_sqlite(ts U.TableDescriptor) error {
	// FIXME Check table name prefix (e.g. "mmmc_") ?
	for _, td := range theMap {
		e := p.createAppTable(&td)
		if e != nil {
			// FIXME
			return e
		}
	}
	L.L.Warning("SQLAR is To-Do, utils/repo/sqlite/impl_apptables.go L114")
	return nil
}

func (p *SqliteRepo) createAppTable(td *DRU.TableDescriptor) error {

	var CTS string // the Create Table SQL string
	var e error

	CTS, e = p.BuildCreateTableStmt(td)
	if e != nil {
		return fmt.Errorf("Cannot crate app table: %s: %w", td.Name, e)
	}
	fnam := "./create-table-" + td.Name + ".sql"
	e = ioutil.WriteFile(fnam, []byte(CTS), 0644)
	if e != nil {
		L.L.Error("Could not write file: " + fnam)
	} else {
		L.L.Dbg("Wrote \"CREATE TABLE " +
			td.Name + " ... \" to: " + fnam)
	}
	// p.Handle().MustExec(CTS)
	_, err := p.Exec(CTS)
	if err != nil {
		// panic(err)
		fmt.Println("p.Exec(CreTbl) barfed: %s", err.Error())
		return err
	}
	ss, e := p.DumpTableSchema_sqlite(td.Name)
	if e != nil {
		return fmt.Errorf("simplerepo.createtable.sqlite: "+
			"DumpTableSchema<%s> failed: %w", e)
	}
	L.L.Dbg(td.Name + " SCHEMA: " + ss)
	// println("TODO: Insert record with IDX 0 and string descr's")
	//    and ("TODO: Dump all table records (i.e. just one)")
	return nil
}
