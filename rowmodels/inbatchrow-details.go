package rowmodels

// Implement interface RowModeler

// TableDetailsINB TBS has no foreign keys.
var TableDetailsINB = TableDetails{
        TableSummaryINB,
        "idx_inbatch", // PKname
        ColumnNamesCsv_INB,
        ColumnSpecsINB, // []D.ColumnSpec
	ColumnPtrsFuncINB,
	&InbatchRow{}, 
}

// TableDetails returns the table
// detail info, given any instance. 
func (inb *InbatchRow) TableDetails() TableDetails {
     return TableDetailsINB
}

/*
// ColumnNamesCSV returns them, given any instance. 
func (inb *InbatchRow) ColumnNamesCSV() string {
     return inb.TableDetails().ColumnNamesCSV
}

// ColumnSpecs returns them, given any instance. 
func (inb *InbatchRow) ColumnSpecs() string {
     return inb.TableDetails().ColumnSpecsINB
}
*/

/* STILL FAILS IN go1.21.5
func PtrFieldsOfGen[T *E, E any](inbro T) []any { // barfs on []db.PtrFields
     switch inbro.(type) {
     }
	// return []any{&inbro.Idx_Inbatch, &inbro.FilCt, &inbro.RelFP,
	//	&inbro.AbsFP, &inbro.T_Cre, &inbro.Descr}
	return []any{1,"hi"}
} */

