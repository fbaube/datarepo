package rowmodels

import (
       "slices"
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
)

var TableSummaryINB = D.TableSummary{
	D.SCT_TABLE.DT(), "inbatch", "inb",
	"Input batch of imported files"}

// This file contains four key items that
// MUST be kept perfectly in sync:
//  - ColumnSpecsINB
//  - ColumnNamesCsvINB
//  - ColumnPtrsINB
//  - struct InbatchRow
//
// SEE FILE ./tabledetails.go for more information.

// PKSpecINB specifies the table's primary key.
// TODO: It should be auto.generated! 
var PKSpecINB = D.ColumnSpec{D.SFT_PRKEY.DT(),
    "idx_inbatch", "Pri.key", "Primary key"} 

// ColumnSpecsINB field order MUST be kept in sync with
// [ColumnNamesCsvINB] and [ColumnPtrsINB]. It specifies:
//   - file count
//   - two path fields (rel & abs) (placed at the end
//     because they tend to be looong)
//   - three time fields (creation, import, last-edit)
//     (the meaning of creation is TBD) 
//   - description
//   - NOT the primary key, which is handled automatically 
// .
var ColumnSpecsINB = []D.ColumnSpec{
	D.ColumnSpec{D.SFT_COUNT.DT(), "filct",
		"Nr. of files", "Number of files"}, 
	D.ColumnSpec{D.SFT_FTEXT.DT(), "descr",
		"Batch descr.", "Inbatch description"}, 
	D.DD_T_Cre, 
	D.DD_T_Imp, 
	D.DD_T_Edt, 
	D.DD_RelFP,
	D.DD_AbsFP,
}

// ColumnNamesCsv_INB is TODO: It should be auto-generated!
var ColumnNamesCsv_INB = "FilCt, Descr, T_Cre, T_Imp, T_Edt, RelFP, AbsFP" 

// ColumnNamesCsvINB is TODO: this can be left unset and 
// then (easily!) auto-generated from [ColumnSpecsINB].
func ColumnNamesCsvINB(inclPK bool) string {
     if !inclPK { return ColumnNamesCsv_INB }
     return "IDX_inbatch, " + ColumnNamesCsv_INB
     }

// ColumnPtrsFuncINB supplies a field in TableDetails.
// NOTE: It MUST be kept in sync:
//  - field order with [ColumnSpecsINB] and [ColumnNamesCsvINB] 
//  - field names with [InbatchRow]
// func ColumnPtrsFuncINB(inbro *InbatchRow, inclPK bool) []any { 
func ColumnPtrsFuncINB(ainbro RowModel, inclPK bool) []any {
     var inbro *InbatchRow
     inbro = ainbro.(*InbatchRow)
     var list []any
     list = []any { &inbro.FilCt, &inbro.Descr,
     	      &inbro.T_Cre, &inbro.T_Imp, &inbro.T_Edt,
	      &inbro.RelFP, &inbro.AbsFP } 
     if !inclPK { return list }
     // names = slices.Insert(names, 1, "Bill", "Billie")
     list = slices.Insert(list, 0, any(&inbro.Idx_Inbatch))
     return list
}

// ColumnPtrsMethod is NOTE: Maybe do the
// "switch (Rowmodeler).RowmodelImplName" trick here. 
func (inbro *InbatchRow) ColumnPtrsMethod(inclPK bool) []any {
     return ColumnPtrsFuncINB(inbro, inclPK) 
}

func (cro *InbatchRow) ColumnNamesCsv(inclPK bool) string {
     return ColumnNamesCsvINB(inclPK)
}

// InbatchRow describes (in the DB) a single import batch
// (probably at the CLI field names); field names MUST be 
// kept in sync with [ColumnPtrsINB].
//  - NOTE: Maybe rename Inbatch* to FileSet*
//    (and INB to FLS) ?
//  - TODO: Maybe represent this (or, each file)
//    with a dsmnd.NSPath: Batch.nr+Path
type InbatchRow struct {
	Idx_Inbatch int
	FilCt int
	Descr string
	RelFP string
	AbsFP FU.AbsFilePath
	T_Cre string
	T_Imp string
	T_Edt string
}

