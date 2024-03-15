package sqlite

import (
       "io"
       "fmt"
       "time"
	"database/sql"
	D "github.com/fbaube/dsmnd"
	// DR "github.com/fbaube/datarepo"
)

type SqliteRepo struct {
	*sql.DB
	filepath string
	w io.Writer // DB logging 
	// SBs map[string]DR.StatementBuilder
	// SEs map[string]StatementEngine
}

func (p *SqliteRepo) DBImplementationName() D.DB_type {
	return D.DB_SQLite
}

func (p *SqliteRepo) SetLogWriter(wrtr io.Writer) io.Writer {
     tmpw := p.w
     p.w = wrtr
     fmt.Fprintf(p.w, "# DB logfile opened at %s \n",
     	time.Now().Local().Format(time.RFC3339))
     return tmpw
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
