package rowmodels

// Implement interface RowModeler

func (tro *TopicrefRow) TableDetails() TableDetails {
     return TableDetails_TopicrefRow
}

func (tro *TopicrefRow) ColumnNamesCSV() string {
     return tro.TableDetails().ColumnNamesCSV
}

