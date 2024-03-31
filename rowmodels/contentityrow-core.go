package rowmodels

import (
       "slices"
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	DRU "github.com/fbaube/datarepo/utils"
	CA "github.com/fbaube/contentanalysis"
	L "github.com/fbaube/mlog"
	// "github.com/fbaube/nurepo/db"
	// "runtime/debug"
)

// TableSummaryCNT summarizes the table.
var TableSummaryCNT = D.TableSummary{
    D.SCT_TABLE.DT(), "contentity", "cnt", "Content entity"}

// This file contains four key items that MUST be kept in sync:
//  - ColumnSpecsINB
//  - ColumnNamesCsvINB
//  - ColumnPtrsINB
//  - struct InbatchRow
//
// SEE FILE ./tabledetails.go for more information.

// PKSpecCNT should be auto.generated!
var PKSpecCNT = D.ColumnSpec{D.SFT_PRKEY.DT(),
    "idx_contentity", "Pri.key", "Primary key"}

// ColumnSpecsCNT field order MUST be kept in sync with
// [ColumnNamesCsvCNT] and [ColumnPtrsCNT] and it specifies:
//   - a primary key (actually, it does NOT - a primary
//     key is assumed, and handled elsewhere)
//   - a foreign key "inbatch"
//   - two path fields (rel & abs)
//   - three time fields (creation, import, last-edit)
//   - a description
//   - three content-type fields (raw markup type, MIME-type, MType); 
//     NOTE: these are persisted in the DB because
//   - - they are useful in searching thru content
//   - - they can be expensive to calculate at import time
//   - - they can be overridden by choices made by users
//   - the content itself
//   - (not for now!) XML content type and XML DOCTYPE
//   - (not for now!) two LwDITA fields (flavor
//     [xdita,hdita!,mdita]), LwDITA content type)
//
// .
var ColumnSpecsCNT = []D.ColumnSpec{
	D.ColumnSpec{D.SFT_FRKEY.DT(), "idx_inbatch", "inbatch",
		"Input batch of imported content"},
	D.DD_RelFP,
	D.DD_AbsFP,
	D.ColumnSpec{D.SFT_FTEXT.DT(), "descr", "Description",
		"Content entity description"},
	D.DD_T_Cre,
	D.DD_T_Imp,
	D.DD_T_Edt,
	D.ColumnSpec{D.SFT_TOKEN.DT(), "rawmt", "Markup type", "Raw markup type"},
	D.ColumnSpec{D.SFT_STRNG.DT(), "mimtp", "MIME type", "MIME type"},
	D.ColumnSpec{D.SFT_STRNG.DT(), "mtype", "MType", "MType"},
	D.ColumnSpec{D.SFT_FTEXT.DT(), "contt", "Content", "Entity raw content"},
	// D.ColSpec{D.SFT_TOKEN.DT(), "xmlcontype", "XML contype", "XML content type"},
	// D.ColSpec{D.SFT_TOKEN.DT(), "xmldoctype", "XML Doctype", "XML Doctype"},
	// D.ColSpec{D.SFT_TOKEN.DT(), "ditaflavor", "LwDITA flavor", "LwDITA flavor"},
	// D.ColSpec{D.SFT_TOKEN.DT(), "ditacontype", "LwDITA contype", "LwDITA cnt type"},
}

// ColumnNamesCsv_CNT TODO: this can be left unset and
// then (easily!) auto-generated from [ColumnSpecsCNT].
var ColumnNamesCsv_CNT =
    "IDX_inbatch, RelFP, AbsFP, Descr, T_Cre, T_Imp, T_Edt, " +
    "RawMT, Mimtp, MType, Contt"

func ColumnNamesCsvCNT(inclPK bool) string {
     if !inclPK { return ColumnNamesCsv_CNT }
     return "IDX_contentity, " + ColumnNamesCsv_CNT
     }


/*
func (RM RowModel) ColumnPtrsMthdCNT(inclPK bool) []any {
     return ColumnPtrsFuncCNT(RM, inclPK)
     }
*/

// ColumnPtrsFuncCNT goes into a field in struct TableDetails
// and MUST be kept in sync:
//  - field order with [ColumnSpecsCNT] and [ColumnNamesCsvCNT]
//  - field names with [ContentityRow]
// func ColumnPtrsFuncCNT(cro *ContentityRow, inclPK bool) []any {
func ColumnPtrsFuncCNT(acro RowModel, inclPK bool) []any {
     if acro == nil {
     	panic("ColumnPtrsFuncCNT nil RowModel")
	}
     var cro *ContentityRow
     cro = acro.(*ContentityRow)
     if cro == nil {
     	L.L.Warning("ColumnPtrsFuncCNT nil *ContentityRow")
	cro = new(ContentityRow)
	}
     var list []any
     if cro.PathAnalysis == nil {
     	if cro.MarkupType == SU.MU_type_DIRLIKE {
	   L.L.Info("Making dummy PathAnalysis for DIR in DB")
	} else {
     	  L.L.Warning("CtyRowCore.ColPtrsFunc: NIL cro.PathAnalysis")
	  // Dump stack (but now comment out cos was
	  // trigrd accidentally by debugging stuff)
	  // debug.PrintStack()
	}
	cro.PathAnalysis = new(CA.PathAnalysis)
     }
     list = []any {
		&cro.Idx_Inbatch,
		&cro.FSItem.FPs.RelFP, &cro.FSItem.FPs.AbsFP,
		&cro.Descr, &cro.T_Cre, &cro.T_Imp, &cro.T_Edt,
		&cro.FSItem.TypedRaw.MarkupType,
		&cro.PathAnalysis.ContypingInfo.MimeType,
		&cro.PathAnalysis.ContypingInfo.MType,
		&cro.FSItem.TypedRaw.Raw, 
		}
	if !inclPK { return list }
	// names = slices.Insert(names, 1, "Bill", "Billie")
	list = slices.Insert(list, 0, any(&cro.Idx_Contentity))
	return list
}

func (cro *ContentityRow) ColumnPtrsMethod(inclPK bool) []any {
     return ColumnPtrsFuncCNT(cro, inclPK)
}

func (cro *ContentityRow) ColumnNamesCsv(inclPK bool) string {
     return ColumnNamesCsvCNT(inclPK)
}

/*

func (cro *ContentityRow) TableDetails() TableDetails {
     return TableDetailsCNT
}

// TableDetailsCNT specifies 11 DB columns,
// incl primary key (assumed) and one foreign key, "inbatch".
var TableDetailsCNT = TableDetails{
        TableSummaryCNT, 
	"idx_contentity", // IDName
	"IDX_inbatch, RelFP, AbsFP, Descr, T_Cre, T_Imp, T_Edt, " +
		"RawMT, Mimtp, MType, Contt", // ColumnNames
	// One foreign key: "inbatch"
	ColumnSpecsCNT, // []D.ColumnSpecs
}

*/

// ContentityRow describes (in the DB) the entity's content
// plus its "dead properties" - basically, properties that
// are set by the user, rather than calculated as needed.
// Raw content is in [FSItem.TypedRaw.Raw] and directory
// typing is in [FSItem.TypedRaw.MarkupType] and they
// are MUTUALLY EXCLUSIVE.
type ContentityRow struct {
	Idx_Contentity int
	Idx_Inbatch    int // NOTE: Rename to FILESET? Could be multiple?
	Descr          string
	// Times is T_Cre, T_Imp, T_Edt string
	DRU.Times
	// FSItem has Raw [the byte content) AND path and 
	// name information AND whether it is a directory
	// (indicated by SU.MarkupType == MU_type_DIRLIKE).
	// NOTE that directory-like and byte-content are
	// mutually exclusive !!
	// CT.TypedRaw { Raw, SU.MarkupType string };
	// RelFP, ShortFP string;
	// FileMeta { os.FileInfo, exists bool, MU.Errer }
	FU.FSItem
	// PathAnalysis is a ptr, so that we get a
	// NPE if it is not initialized properly;
	// or if analysis failed, if (for example)
	// the content is too short.
	// FU.PathAnalysis is
	// XU.ContypingInfo { FileExt, MimeType, =>
	//   MimeTypeAsSnift, MType string }
	// ContentityBasics { XmlRoot, Text, Meta CT.Span; // => TopLevel !!
	//     MetaFormat string; MetaProps SU.PropSet }
	// XmlContype string
	// *XU.ParsedPreamble
	// *XU.ParsedDoctype
	// DitaFlavor  string
	// DitaContype string
	*CA.PathAnalysis // NEED DETAIL
	// Contt string

	// For these next two fields, instead put the refs & defs
	//   into another table that FKEY's into this table.
	// ExtlLinkRefs // links that point outside this File
	// ExtlLinkDefs // link targets in-file that are visible outside this File
	// Linker = an outgoing link
	// Linkee = the target of an outgoing link
	// Linkable = a symbol that CAN be a Linkee
}

