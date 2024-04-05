package sqlite

import (
       "io"
       "fmt"
       "errors"
	"database/sql"
	D "github.com/fbaube/dsmnd"
	SU "github.com/fbaube/stringutils"
)

// SqliteRepo implements [datarepo/SimpleRepo]
/*
InitznPragmas
IsInTx
Tx
Tx for write Conn only
*/
type SqliteRepo struct {
	*sql.DB
	filepath string
	w io.Writer // DB logging 
	// SBs map[string]DR.StatementBuilder
	// SEs map[string]StatementEngine
	// TODO: Tx should be protected by a Mutex
	*sql.Tx 
}

func (p *SqliteRepo) DBImplementationName() D.DB_type {
	return D.DB_SQLite
}

func (p *SqliteRepo) GetTx() *sql.Tx {
	return p.Tx
}

func (p *SqliteRepo) IsInTx() bool {
	return p.Tx != nil 
}

func (p *SqliteRepo) Commit() error {
     if p.Tx == nil {
     	return errors.New("SqliteRepo.Commit: no active transaction")
	}
     e := p.Tx.Commit()
     p.Tx = nil 
     if e != nil {
     	return fmt.Errorf("SqliteRepo.Commit failed: %w", e)
	}
     return nil
}

func (p *SqliteRepo) Rollback() error {
     if p.Tx == nil {
     	return errors.New("SqliteRepo.Rollback: no active transaction")
	}
     e := p.Tx.Rollback()
     p.Tx = nil 
     if e != nil {
     	return fmt.Errorf("SqliteRepo.Rollback failed: %w", e)
	}
     return nil
}

func (p *SqliteRepo) SetLogWriter(wrtr io.Writer) io.Writer {
     tmpw := p.w
     p.w = wrtr
     fmt.Fprintf(p.w, "# DB logfile opened at %s \n", SU.Now())
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

func (p *SqliteRepo) Exec(query string, args ...any) (sql.Result, error) {
     if p.Tx != nil {
     	return p.Tx.Exec(query, args...)
     } else {
        return p.DB.Exec(query, args...)
     }
}

func (p *SqliteRepo) Query(query string, args ...any) (*sql.Rows, error) {
     if p.Tx != nil {
     	return p.Tx.Query(query, args...)
     } else {
        return p.DB.Query(query, args...)
     }
}

func (p *SqliteRepo) QueryRow(query string, args ...any) *sql.Row {
     if p.Tx != nil {
     	return p.Tx.QueryRow(query, args...)
     } else {
        return p.DB.QueryRow(query, args...)
     }
}
