package rowmodels

import (
	DRU "github.com/fbaube/datarepo/utils"
)

// Implement interface RowModeler

func (tro *TopicrefRow) TableDetails() DRU.TableDetails {
     return TableDetails_TopicrefRow
}

func (tro *TopicrefRow) ColumnNamesCSV() string {
     return tro.TableDetails().ColumnNamesCSV
}

// TODO: Still can't sort out the notation for ptr constraints ?!
func (tro *TopicrefRow) ColumnPtrs() []any { // barfs on []db.PtrFields
	return []any{&tro.Idx_Map_Contentity, &tro.Idx_Tpc_Contentity}
}

// TODO: Still can't sort out the notation for ptr constraints ?!
func ColumnPtrsTRF(tro *TopicrefRow) []any { // barfs on []db.PtrFields
	return []any{&tro.Idx_Map_Contentity, &tro.Idx_Tpc_Contentity}
}

