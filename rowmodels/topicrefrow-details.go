package rowmodels

// TableDetailsTRF specifies only two foreign keys.
var TableDetailsTRF = TableDetails{
        TableSummaryTRF, 
	"idx_topicref", // PKname
	ColumnNamesCsvTRF, // "idx_cnt_map, idx_cnt_tpc", // ColumnNames
	ColumnSpecsTRF, // []D.ColumnSpec
	ColumnPtrsTRF, 
}

// TableDetails returns the table
// detail info, given any instance.
func (tro *TopicrefRow) TableDetails() TableDetails {
     return TableDetailsTRF
}

