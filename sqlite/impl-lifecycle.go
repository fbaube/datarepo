package sqlite

// SQLiteConn implements driver.Conn

// SQLiteContext behave sqlite3_context
// https://www.sqlite.org/c3ref/context.html
// The context in which an SQL function executes
// is stored in an sqlite3_context object. A ptr
// to an sqlite3_context object is always the
// first parameter to app-defined SQL functions.

// https://www.sqlite.org/backup.html
// Function << sqlite3_backup_init() >> is called to create
// an sqlite3_backup object to copy data from database pDb
// to the backup database file identified by zFilename.
// Function << sqlite3_backup_step() >> is called with a
// parameter of 5 to copy 5 pages of database pDb to the
// backup database (file zFilename).
// If there are still more pages to copy from database pDb,
// then the function sleeps for 250 milliseconds (using the
// << sqlite3_sleep() >> utility) and then returns to step 2.
// Function << sqlite3_backup_finish() >> is called to clean
// up resources allocated by << sqlite3_backup_init() >>

import (
	"database/sql"
	"fmt"
	"strings"
)

// RepoLifecycle is lifecycle operations for databases.
//
// Expected constructors for any struct "pkg.XxxRepo" that implements this:
//  - Default call, optimized for a DB that is assumed to existing already,
//    but it can fall back (without error) on NewRepoAtPath(..):
//  - OpenRepoAtPath(string) (pkg.XxxRepo, string, error)
//  - Call optimized for a path with no DB (yet), and if something is there,
//    it will return an error:
//  - NewRepoAtPath(string)  (pkg.XxxRepo, string, error) // does not Open() tho
//
// type RepoLifecycle interface {

// BONUS FUNC
func (p *SqliteRepo) DoPragma(s string) {
	if !strings.HasPrefix(s, "PRAGMA") {
		s = "PRAGMA " + s
	}
	_, e := p.DB.Exec(s)
	if e != nil {
		panic("SQLite PRAGMA failed: " + e.Error())
	}
}

// Open also, if needed, does pragma-style initialization.
func (p *SqliteRepo) Open() error {
	// println("Open")
	return nil
}

// NewAtPath also, if needed, does pragma-style initialization.
func (p *SqliteRepo) NewAtPath() error {
	// println("Open")
	return nil
}

// IsOpen also, if possible, pings the DB as a health check.
func (p *SqliteRepo) IsOpen() bool {
	if p.DB == nil {
		return false
	}
	e := p.DB.Ping()
	return (e == nil)
}

// Flush is basically a no-op.
func (p *SqliteRepo) Flush() error {
	return nil
}

// Close remembers the path.
func (p *SqliteRepo) Close() error {
	// println("Close")
	// Conn.Close()
	e := p.DB.Close()
	if e != nil {
		println("db.close failed:", e.Error())
		return fmt.Errorf("sqliterepo.close failed: %w", e)
	}
	return e
}

// Verify runs app-level sanity & consistency checks (things
// like foreign key validity should be delegated to DB setup).
func (p *SqliteRepo) Verify() error {
	var stmt *sql.Stmt
	var e error
	stmt, e = p.Handle().Prepare("PRAGMA integrity_check;")
	if e == nil {
		return e
	}
	_, e = stmt.Exec() // rslt,e := ...
	if e == nil {
		return e
	}
	stmt, e = p.Handle().Prepare("PRAGMA foreign_key_check;")
	if e == nil {
		return e
	}
	_, e = stmt.Exec() // rslt,e := ...
	if e == nil {
		return e
	}
	return nil
	// liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("DD:mustExecStmt: ID %d nR %d \n", liid, naff)
}
