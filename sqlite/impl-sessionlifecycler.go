package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	// L "github.com/fbaube/mlog"
)

// SessionLifecycler is session lifecycle operations for databases.

/*
var initializationPragmas =
`PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 1000000000;
PRAGMA foreign_keys = true;
PRAGMA temp_store = memory;`

var readonlyPragma = "PRAGMA query_only = TRUE;"
*/

func (p *SqliteRepo) DoPragma(s string) {
	if !strings.HasPrefix(s, "PRAGMA") {
		s = "PRAGMA " + s
	}
	_, e := p.DB.Exec(s)
	if e != nil {
		panic("SQLite PRAGMA failed: " + e.Error())
	}
}

// PRAGMA schema.integrity_check;

// PRAGMA schema.journal_mode = DELETE (dflt) | WAL
// WAL journaling mode is persistent; after being set,
// it stays in effect across multiple DB connections
// and after closing and reopening the DB.

// "the pragma value is associated with the connection object."

// "When you are using PRAGMA user_version to read the value,
// this statement behaves exactly like a query, i.e., SELECT
// user_version FROM somewhere. So just use Query()."

// DoPragmas can handle multiline, but does expect the caller
// to provide all required instances of the PRAGMA keyword.
// The PRAGMA menu: https://www.sqlite.org/pragma.html
//
// NOTE: "No error messages are generated if an unknown pragma
// is issued. Unknown pragmas are simply ignored. This means if
// there is a typo in a pragma statement, the library does not
// inform the user of the fact."
//
// A pragma may have an optional schema-name before the pragma name.
//
// PRAGMAs that return results and that have no side-effects can be
// accessed from ordinary SELECT statements as table-valued functions. 
//
// Arguments:
//  - A pragma can take either zero or one argument.
//  - The argument is may be either in parentheses or it may
//    ae separated from the pragma name by an equal sign; the
//    two syntaxes yield identical results.
//  - In many pragmas, the argument is a boolean, and can be one of:
//    1 yes true on / 0 no false off
//  - Keyword arguments can optionally appear in quotes. (Example:
//    'yes' [FALSE].)
//  - Some pragmas takes a string literal as their argument. When
//    pragma takes a keyword argument, it will usually also take 
//    a numeric equivalent as well. For example, "0" and "no" mean
//    the same thing, as does "1" and "yes".
//  - When querying the value of a setting, many pragmas return the
//    number rather than the keyword.
// .
func (p *SqliteRepo) DoPragmas(s string) (string, error) {
	_, e := p.DB.Exec(s)
	if e != nil {
		return "", fmt.Errorf("SQLite PRAGMA failed: %w", e)
	}
	return "", nil
}

// Open should not do pragma-style initialization.
func (p *SqliteRepo) Open() error {
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

// Flush forces WAL sync. Notes:
// https://turso.tech/blog/something-you-probably-want-to-know-about-if-youre-using-sqlite-in-golang-72547ad625f1
//
// A "Rows" object represents a read operation to an SQLite DB.
// Unless rows is closed, explicitly by calling rows.Close() or
// implicitly by reading all the data from it, SQLite treats the
// operation as ongoing. It turns out that ongoing reads prevent
// checkpoint operation from advancing. As a result, one leaked
// Rows object can prevent the DB from ever transferring changes
// from the WAL file to the main DB file and truncating the WAL
// file. This means the WAL file grows indefinitely. At least
// until the process is restarted, which can take a long time
// for a server application.
// .
func (p *SqliteRepo) Flush() error {
     _, err := p.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
     if err != nil {
     	return fmt.Errorf("SQLiteRepo.Flush: PRAGMA wal_checkpoint: %w", err)
 	}
     return nil 
}

// Close remembers the path.
func (p *SqliteRepo) Close() error {
	// println("Close")
	// Conn.Close()
	e := p.DB.Close()
	if e != nil {
		// println("db.close failed:", e.Error())
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
