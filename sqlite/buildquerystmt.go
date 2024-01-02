package sqlite

import RU "github.com/fbaube/repoutils"

func (p *SqliteRepo) BuildQueryStmt(qs *RU.QuerySpec) (string, error) {
	switch qs.DbOp {
	case RU.OpCreateTable:

	}
	return "", nil
}
