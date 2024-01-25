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

