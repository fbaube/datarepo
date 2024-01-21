package rowmodels

import (
	DRU "github.com/fbaube/datarepo/utils"
	// "github.com/fbaube/nurepo/db"
)

// Implement interface RowModeler

func (inb *InbatchRow) TableDetails() DRU.TableDetails {
     return TableDetailsINB
}

func (inb *InbatchRow) ColumnNamesCSV() string {
     return inb.TableDetails().ColumnNamesCSV
}

/* STILL FAILS IN go1.21.5
func PtrFieldsOfGen[T *E, E any](inbro T) []any { // barfs on []db.PtrFields
     switch inbro.(type) {
     }
	// return []any{&inbro.Idx_Inbatch, &inbro.FilCt, &inbro.RelFP,
	//	&inbro.AbsFP, &inbro.T_Cre, &inbro.Descr}
	return []any{1,"hi"}
} */

