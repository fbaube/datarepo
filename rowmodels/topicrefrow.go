package rowmodels

import (
	"fmt"
	D "github.com/fbaube/dsmnd"
	DRU "github.com/fbaube/datarepo/utils"
)

/*
var TableSummary_ContentityRow = D.TableSummary{
var TableDescriptor_ContentityRow = TableDescriptor{
func (cro *ContentityRow) PtrFields() []any { // barfs on []db.PtrFields
var ColumnSpecs_ContentityRow = []D.ColumnSpec{
type ContentityRow struct {
func (p *ContentityRow) String() string {
*/

// TableSummary_TopicrefRow describes the table.
var TableSummary_TopicrefRow = D.TableSummary{D.SCT_TABLE.DT(),
	"trf", "topicref", "Reference from map to topic"}

// TableDescriptor_TopicrefRow specifies only two foreign keys.
var TableDescriptor_TopicrefRow = DRU.TableDescriptor{
	"topicref",     // Name
	"trf",          // ShortName
	"idx_topicref", // IDName
	//
	"idx_map_contentity, idx_tpc_contentity", // ColumnNames
	// ONLY foreign keys
	// []string{"map_contentity", "tpc_contentity"},
	ColumnSpecs_TopicrefRow, // []D.ColumnSpec
}

// TODO: Still can't sort out the notation for ptr constraints ?!
func (tro *TopicrefRow) PtrFields() []any { // barfs on []db.PtrFields
	return []any{&tro.Idx_Map_Contentity, &tro.Idx_Tpc_Contentity}
}

// ColumnSpecs_TopicrefRow is empty, cos
// the table contains only foreign keys.
var ColumnSpecs_TopicrefRow = []D.ColumnSpec{
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_cnt_map", "contentity",
		"Referencing map"},
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_cnt_tpc", "contentity",
		"Referenced topic"},
}

// TopicrefRow describes a reference from a Map (i.e. TOC) to a Topic.
// Note that "Topic" does NOT necessarily refer to a DITA `topictref`
// element!
//
// The relationship is N-to-N btwn Maps and Topics, so a TopicrefRow
// might not be unique because a topic might be explicitly referenced
// more than once by a map. So for simplicity, let's create only one
// TopicrefRow per topic per map file, and see if it creates problems
// elsewhere later on.
//
// Note also that if we decide to use multi-trees, then perhaps these links
// can count not just as kids for maps, but also as parents for topics.
type TopicrefRow struct {
	Idx_Topicref       int
	Idx_Map_Contentity int
	Idx_Tpc_Contentity int
}

// TODO Write col desc's using Desmond !
// TODO Generate ColNames from ColumnSpecs_TopicrefRow

// String implements Stringer. FIXME
func (p *TopicrefRow) String() string {
	return fmt.Sprintf("topicrefrow FIXME")
}
