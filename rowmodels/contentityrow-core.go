package rowmodels

import (
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
	DRU "github.com/fbaube/datarepo/utils"
	CA "github.com/fbaube/contentanalysis"
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

// ColumnNamesCsvCNT TODO: this can be left unset and
// then (easily!) auto-generated from [ColumnSpecsCNT].
var ColumnNamesCsvCNT =
    "IDX_inbatch, RelFP, AbsFP, Descr, T_Cre, T_Imp, T_Edt, " +
    "RawMT, Mimtp, MType, Contt" // ColumnNames

// ColumnPtrsCNT goes into TableDetails and MUST be kept in sync:
//  - field order with [ColumnSpecsCNT] and [ColumnNamesCsvCNT]
//  - field names with [ContentityRow]
func ColumnPtrsCNT(cro *ContentityRow, inclPK bool) []any {
     var list []any
     if cro.PathAnalysis == nil {
     	println("NIL cro.PathAnalysis !!")
	// Dump stack (was trigrd accidentally by debugging stuff)
	// debug.PrintStack()
	cro.PathAnalysis = new(CA.PathAnalysis)
	}
     list = []any {
		// &cro.Idx_Contentity,
		&cro.Idx_Inbatch,
		&cro.PathProps.RelFP, &cro.PathProps.AbsFP,
		&cro.Descr, &cro.T_Cre, &cro.T_Imp, &cro.T_Edt,
		&cro.PathProps.TypedRaw.MarkupType,
		&cro.PathAnalysis.ContypingInfo.MimeType,
		&cro.PathAnalysis.ContypingInfo.MType,
		&cro.PathProps.TypedRaw.Raw, 
		}
	if !inclPK { return list }
	var pk []any
	pk = []any { &cro.Idx_Contentity }
	return append(pk, list...)
}

func (cro *ContentityRow) ColumnPtrs(inclPK bool) []any {
     return ColumnPtrsCNT(cro, inclPK)
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
// The Raw content is in [PathProps.TypedRaw.Raw].
type ContentityRow struct {
	Idx_Contentity int
	Idx_Inbatch    int // NOTE: Rename to FILESET? Could be multiple?
	Descr          string
	// Times is T_Cre, T_Imp, T_Edt string
	DRU.Times
	// PathProps has Raw and is // => EntityProps !!
	// CT.TypedRaw { Raw, SU.MarkupType string };
	// RelFP, ShortFP string;
	// FileMeta { os.FileInfo, exists bool, MU.Errer }
	FU.PathProps
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

