package rowmodels

import (
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
)

var TableSummaryINB = D.TableSummary{
	D.SCT_TABLE.DT(), "inbatch", "inb",
	"Input batch of imported files"}

// For implementing interface RowModeler: This file
// contains four key items that MUST be kept in sync:
//  - ColumnSpecsINB
//  - ColumnPtrsINB
//  - InbatchRow
//  - ColumnNamesCsvINB
//
// The order of fields is quite flexible, and so because 
// of how fields are displayed in DB tools, shorter and 
// more-important fields should appear first.
//
// NOTE that:
//  - in principle, both variables and code can 
//    be auto-generated based on [ColumnSpecsINB],
//  - BUT it would be a serious task to do so,
//  - AND might be useless and/or impossible when
//    using generics. 

// ColumnNamesCsvINB can be left unset and then 
// easily auto-generated from [ColumnSpecsINB].
var ColumnNamesCsvINB = "FilCt, Descr, T_Cre, T_Imp, T_Edt, RelFP, AbsFP" 

var PKSpecINB = D.ColumnSpec{D.SFT_PRKEY.DT(),
    "idx_inbatch", "Pri.key", "Primary key"} 

func (inbro *InbatchRow) ColumnPtrs(inclPK bool) []any {
     return ColumnPtrsINB(inbro, inclPK) 
}

// ColumnSpecsINB field order MUST be kept in 
// sync with [ColumnNamesCsvINB] and [ColumnPtrsINB] and it specifies:
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

// ColumnPtrsINB MUST be kept in sync:
//  - field order with [ColumnNamesCsvINB] and [ColumnSpecsINB]
//  - field names with [InbatchRow]
func ColumnPtrsINB(inbro *InbatchRow, inclPK bool) []any { 
     var list []any
     list = []any { &inbro.FilCt, &inbro.Descr,
     	      &inbro.T_Cre, &inbro.T_Imp, &inbro.T_Edt,
	      &inbro.RelFP, &inbro.AbsFP } 
     if !inclPK { return list }
     var pk []any
     pk = []any { &inbro.Idx_Inbatch }
     return append(pk, list...)
}

// InbatchRow field names MUST be kept in sync 
// with [ColumnPtrsINB] and a record describes 
// (in the DB) a single import batch at the CLI.
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

