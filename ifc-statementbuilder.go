package datarepo

import DRU "github.com/fbaube/datarepo/utils"

type StatementBuilder interface {
	BuildQueryStmt(*DRU.QuerySpec) (string, error)
	BuildCreateTableStmt(*DRU.TableDetails) (string, error)
}
