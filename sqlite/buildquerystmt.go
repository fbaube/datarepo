package sqlite

import DRU "github.com/fbaube/datarepo/utils"

func (p *SqliteRepo) BuildQueryStmt(qs *DRU.QuerySpec) (string, error) {
	switch qs.DbOp {
	case DRU.OpCreateTable:

	}
	return "", nil
}
