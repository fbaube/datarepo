package rowmodels

import (
	"fmt"
)

// TableDetailsTRF specifies only two foreign keys.
var TableDetailsTRF = TableDetails{
        TableSummaryTRF, 
	"idx_topicref", // IDName
	"idx_cnt_map, idx_cnt_tpc", // ColumnNames
	ColumnSpecsTRF, 
}

// TODO Write col desc's using Desmond !
// TODO Generate ColNames from ColumnSpecsTRF

// String implements Stringer. FIXME
func (p *TopicrefRow) String() string {
	return fmt.Sprintf("topicrefrow FIXME")
}
