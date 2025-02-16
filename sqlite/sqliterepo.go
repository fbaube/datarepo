package sqlite

import(
	"io"
	"fmt"
	"database/sql"
	D "github.com/fbaube/dsmnd"
	SU "github.com/fbaube/stringutils"
)

// SqliteRepo implements [datarepo/SimpleRepo].
// 
// 2025.02 sql.Tx is removed from this, because functions
// using it add unnecessary, user-surprising complexity
// that does not appear in Go third party libraries, and
// anyways transactions should be managed by repo users.
type SqliteRepo struct {
	*sql.DB
	filepath string
	w io.Writer // DB logging 
	// SBs map[string]DR.StatementBuilder
	// SEs map[string]StatementEngine
	// *sql.Tx // 2025.02 Removed 
}

func (p *SqliteRepo) DBImplementationName() D.DB_type {
	return D.DB_SQLite
}

func (p *SqliteRepo) SetLogWriter(wrtr io.Writer) io.Writer {
     tmpw := p.w
     p.w = wrtr
     fmt.Fprintf(p.w, "# DB logfile opened at %s \n", SU.Now())
     return tmpw
}

func (p *SqliteRepo) LogWriter() io.Writer {
     return p.w
}

func (p *SqliteRepo) CloseLogWriter() {
     var c io.Closer
     var ok bool 
     c, ok = p.w.(io.Closer)
     if ok {
     	// ignore any error 
     	c.Close()
	}
     p.w = io.Discard
}

/* init() to do type chex
(but don't do DB stuff in init(), cos
 driver might not be initialized yet)
func init() {
	var sr *SqliteRepo
	var x1 RepoAppTables
	sr, _, _ = NewSqliteRepoAtPath("/tmp/tmp")
	x1, ok = (repo.RepoAppTables)(sr)
	// x2, ok := sr.(Backupable)
	// x3, ok := sr.(RepoEntity)
	// x4, ok := sr.(RepoLifecycle)
}
*/

