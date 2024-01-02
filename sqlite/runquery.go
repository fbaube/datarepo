package sqlite

import (
	"fmt"
	RU "github.com/fbaube/repoutils"
)

func (p *SqliteRepo) RunQuery0(*RU.QuerySpec) (any, error) { // ie. Exec()
	fmt.Println("NOT IMPL'D: RunQuery0")
	return nil, nil
}

func (p *SqliteRepo) RunQuery1(*RU.QuerySpec) (any, error) { // One row, like by_ID
	fmt.Println("NOT IMPL'D: RunQuery1")
	return nil, nil
}

func (p *SqliteRepo) RunQueryN(*RU.QuerySpec) ([]any, error) { // Multiple rows
	fmt.Println("NOT IMPL'D: RunQueryN")
	return nil, nil
}
