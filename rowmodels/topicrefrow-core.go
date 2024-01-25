package rowmodels

import (
	D "github.com/fbaube/dsmnd"
)

var TableSummaryTRF = D.TableSummary{
	D.SCT_TABLE.DT(), "topicref", "trf",
	"Reference from map to topic"}

// For implementing interface RowModeler: This file
// contains four key items that MUST be kept in sync:
//  - ColumnSpecsTRF
//  - ColumnPtrsTRF
//  - TopicrefRow
//  - ColumnNamesCsvTRF
//
// The order of fields is quite flexible, and so because 
// of how fields are displayed in DB tools, shorter and 
// more-important fields should appear first.
//
// NOTE that:
//  - in principle, both variables and code can 
//    be auto-generated based on [ColumnSpecsTRF],
//  - BUT it would be a serious task to do so,
//  - AND might be useless and/or impossible when
//    using generics. 

// ColumnNamesCsvTRF can be left unset and then 
// easily auto-generated from [ColumnSpecsTRF].
var ColumnNamesCsvTRF = "idx_cnt_map, idx_cnt_tpc" // ColumnNames

var PKSpecTRF = D.ColumnSpec{D.SFT_PRKEY.DT(),
    "idx_topicref", "Pri.key", "Primary key"} 

func (tro *TopicrefRow) ColumnPtrs(inclPK bool) []any {
     return ColumnPtrsTRF(tro, inclPK) 
}

// ColumnSpecsTRF field order MUST be kept in 
// sync with [ColumnNamesCsvTRF] and [ColumnPtrsTRF] and it specifies:
//   - file count
//   - two path fields (rel & abs) (placed at the end
//     because they tend to be looong)
//   - three time fields (creation, import, last-edit)
//     (the meaning of creation is TBD) 
//   - description
//   - NOT the primary key, which is handled automatically 
// .
var ColumnSpecsTRF = []D.ColumnSpec{
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_cnt_map", "contentity",
		"Referencing map"},
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_cnt_tpc", "contentity",
		"Referenced topic"},
}

// ColumnPtrsTRF MUST be kept in sync:
//  - field order with [ColumnNamesCsvTRF] and [ColumnSpecsTRF]
//  - field names with [TopicrefRow]
func ColumnPtrsTRF(tro *TopicrefRow, inclPK bool) []any { 
     var list []any 
     list = []any { &tro.Idx_Map_Contentity, &tro.Idx_Tpc_Contentity }
     if !inclPK { return list }
     var pk []any
     pk = []any { &tro.Idx_Topicref }
     return append(pk, list...)
}

// TopicrefRow describes a reference from a Map (i.e. TOC) to a Topic.
// Note that "Topic" does NOT necessarily refer to a DITA `topictref`
// element!
//
// TopicrefRow field names MUST be kept in sync 
// with [ColumnPtrsTRF] and a record describes 
// (in the DB) a single import batch at the CLI.
//
// The relationship is N-to-N btwn Maps and Topics, so a TopicrefRow
// might not be unique because a map might explicitly reference a 
// particular topic more than once. So for simplicity, let's create 
// only one TopicrefRow per map/topic pair, and see if it creates 
// problems elsewhere later on. Maybe a record needs a "count" field. 
//
// Note also that if we decide to use multi-trees, then perhaps these links
// can count not just as kids for maps, but also as parents for topics.
// .
type TopicrefRow struct {
	Idx_Topicref       int
	Idx_Map_Contentity int
	Idx_Tpc_Contentity int
}

