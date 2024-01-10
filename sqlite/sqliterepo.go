package sqlite

import (
	"database/sql"
	D "github.com/fbaube/dsmnd"
	// DR "github.com/fbaube/datarepo"
)

type SqliteRepo struct {
	*sql.DB
	filepath string
	// SBs map[string]DR.StatementBuilder
	// SEs map[string]StatementEngine
}

func (p *SqliteRepo) DBImplementationName() D.DB_type {
	return D.DB_SQLite
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
