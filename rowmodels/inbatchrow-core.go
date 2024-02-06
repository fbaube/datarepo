package rowmodels

import (
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
)

var TableSummaryINB = D.TableSummary{
	D.SCT_TABLE.DT(), "inbatch", "inb",
	"Input batch of imported files"}

// This file contains four key items that MUST be kept in sync:
//  - ColumnSpecsINB
//  - ColumnNamesCsvINB
//  - ColumnPtrsINB
//  - struct InbatchRow
//
// SEE FILE ./tabledetails.go for more information.

// PKSpecINB should be auto.generated! 
var PKSpecINB = D.ColumnSpec{D.SFT_PRKEY.DT(),
    "idx_inbatch", "Pri.key", "Primary key"} 

// ColumnSpecsINB field order MUST be kept in sync with
// [ColumnNamesCsvINB] and [ColumnPtrsINB] and it specifies:
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

// ColumnNamesCsvINB TODO: this can be left unset and 
// then (easily!) auto-generated from [ColumnSpecsINB].
var ColumnNamesCsvINB = "FilCt, Descr, T_Cre, T_Imp, T_Edt, RelFP, AbsFP" 

// ColumnPtrsFuncINB goes into TableDetails and MUST be kept in sync:
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
     var pk []any
     pk = []any { &inbro.Idx_Inbatch }
     return append(pk, list...)
}

func (inbro *InbatchRow) ColumnPtrsMethod(inclPK bool) []any {
     return ColumnPtrsFuncINB(inbro, inclPK) 
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

