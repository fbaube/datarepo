package rowmodels

import (
	D "github.com/fbaube/dsmnd"
)

var TableSummaryTRF = D.TableSummary{
	D.SCT_TABLE.DT(), "topicref", "trf",
	"Reference from map to topic"}

// This file contains four key items that MUST be kept in sync:
//  - ColumnSpecsTRF
//  - ColumnNamesCsvTRF
//  - ColumnPtrsTRF
//  - struct TopicrefRow
//
// SEE FILE ./tabledetails.go for more information.

// PKSpecTRF should be auto.generated!
var PKSpecTRF = D.ColumnSpec{D.SFT_PRKEY.DT(),
    "idx_topicref", "Pri.key", "Primary key"} 

// ColumnSpecsTRF field order MUST be kept in sync with
// [ColumnNamesCsvTRF] and [ColumnPtrsTRF] and it specifies:
//   - the primary index into table "contentity"
//     of the map that makes the reference
//   - the primary index into table "contentity"
//     of the topic that is referred to
//   - NOT the primary key, which is handled automatically 
//
// Note that it is unclear ATM so far whether this 
// table also includes maps referring to submaps,
// or topics referring to other topics.
// .
var ColumnSpecsTRF = []D.ColumnSpec{
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_cnt_map", "contentity",
		"Referencing map"},
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_cnt_tpc", "contentity",
		"Referenced topic"},
}

// ColumnNamesCsvTRF TODO: this can be left unset and 
// then (easily!) auto-generated from [ColumnSpecsTRF].
var ColumnNamesCsvTRF = "idx_cnt_map, idx_cnt_tpc" 

// ColumnPtrsFuncTRF MUST be kept in sync:
//  - field order with [ColumnNamesCsvTRF] and [ColumnSpecsTRF]
//  - field names with [TopicrefRow]
// func ColumnPtrsFuncTRF(tro *TopicrefRow, inclPK bool) []any { 
func ColumnPtrsFuncTRF(atro RowModel, inclPK bool) []any {
     var tro *TopicrefRow
     tro = atro.(*TopicrefRow)
     var list []any 
     list = []any { &tro.Idx_Map_Contentity, &tro.Idx_Tpc_Contentity }
     if !inclPK { return list }
     var pk []any
     pk = []any { &tro.Idx_Topicref }
     return append(pk, list...)
}

func (tro *TopicrefRow) ColumnPtrsMethod(inclPK bool) []any {
     return ColumnPtrsFuncTRF(tro, inclPK) 
}

// TopicrefRow describes (in the DB) a reference 
// from a Map (i.e. TOC) to a Topic; field names 
// MUST be kept in sync with [ColumnPtrsTRF]. 
//
// Note that "Topic" does NOT necessarily refer 
// to a DITA `topictref`element!
//
// The relationship is N-to-N btwn Maps and Topics, so a TopicrefRow
// might not be unique because a map might explicitly reference a 
// particular topic more than once. So for simplicity, let's create 
// only one TopicrefRow per map/topic pair, and see if it creates 
// problems elsewhere later on. Maybe a record also needs a "count" field. 
//
// Note also that if we decide to use multi-trees, then perhaps these links
// can count not just as kids for maps, but also as parents for topics.
// .
type TopicrefRow struct {
	Idx_Topicref       int
	Idx_Map_Contentity int
	Idx_Tpc_Contentity int
}

