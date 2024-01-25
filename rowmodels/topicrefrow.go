package rowmodels

import (
	"fmt"
	D "github.com/fbaube/dsmnd"
	DRU "github.com/fbaube/datarepo/utils"
)

// TableDetails_TopicrefRow specifies only two foreign keys.
var TableDetails_TopicrefRow = DRU.TableDetails{
        TableSummaryTRF, 
	"idx_topicref", // IDName
	"idx_cnt_map, idx_cnt_tpc", // ColumnNames
	ColumnSpecs_TopicrefRow, 
}

// ColumnSpecs_TopicrefRow contains only foreign keys.
var ColumnSpecs_TopicrefRow = []D.ColumnSpec{
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_cnt_map", "contentity",
		"Referencing map"},
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_cnt_tpc", "contentity",
		"Referenced topic"},
}

// TODO Write col desc's using Desmond !
// TODO Generate ColNames from ColumnSpecs_TopicrefRow

// String implements Stringer. FIXME
func (p *TopicrefRow) String() string {
	return fmt.Sprintf("topicrefrow FIXME")
}
