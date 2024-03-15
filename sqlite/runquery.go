package sqlite

import (
	"fmt"
	DRU "github.com/fbaube/datarepo/utils"
)

func (p *SqliteRepo) RunQuery0(*DRU.QuerySpec) (any, error) { // ie. Exec()
	fmt.Fprintf(p.w, "NOT IMPL'D: RunQuery0 \n")
	return nil, nil
}

func (p *SqliteRepo) RunQuery1(*DRU.QuerySpec) (any, error) { // One row, like by_ID
	fmt.Fprintf(p.w, "NOT IMPL'D: RunQuery1 \n")
	return nil, nil
}

func (p *SqliteRepo) RunQueryN(*DRU.QuerySpec) ([]any, error) { // Multiple rows
	fmt.Fprintf(p.w, "NOT IMPL'D: RunQueryN \n")
	return nil, nil
}
