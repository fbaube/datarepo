package sqlite

import (
	"fmt"
	DRU "github.com/fbaube/datarepo/utils"
)

func (p *SqliteRepo) RunQuery0(*DRU.QuerySpec) (any, error) { // ie. Exec()
	fmt.Println("NOT IMPL'D: RunQuery0")
	return nil, nil
}

func (p *SqliteRepo) RunQuery1(*DRU.QuerySpec) (any, error) { // One row, like by_ID
	fmt.Println("NOT IMPL'D: RunQuery1")
	return nil, nil
}

func (p *SqliteRepo) RunQueryN(*DRU.QuerySpec) ([]any, error) { // Multiple rows
	fmt.Println("NOT IMPL'D: RunQueryN")
	return nil, nil
}
