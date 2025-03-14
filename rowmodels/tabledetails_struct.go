package rowmodels

import(
	D "github.com/fbaube/dsmnd"
)

// columnPtrsFunc is used in struct [TableDetails.ColumnPtrsFunc],
// below. Note that ColumnPtrsFunc is a func, and/but while there
// is also a method on interface [RowModel] with the signature: 
// (p RowModel) ColumnPtrsMethod(bool) []any

type columnPtrsFunc func(RowModel, bool) []any 
// type columnPtrsMthd func(bool) []any 

type NewInstanceFunc func() RowModel

// TableDetails is key to the entire data repository scheme:
// it contains metadata required to manage corresppondences
// between DB columns and struct fields. It specifies an 
// application's data schema requirements both at the table
// level and field-by-field. A primary key is assumed for
// every table, and foreign keys are allowed.
//
// NOTE: The case of multiple foreign indices into the same 
// table (from TopicrefRow into ContentityRow) is 90% resolved. 
//
// As a convention, all DB table and column names should be 
// in all lower case, except perhaps in the names of indices, 
// which have names of the form Idx_*. Then SQL keywords can
// be given in all upper case, as per SQL convention. 
// (Note that enforcing all this might be a bit patchy.) 
//
// Notes on particular fields:
//   - The field [ColumnSpecs] is a slice of [dsmnd.ColumnSpec],
//     which has four text fields:
//     [dsmnd.BasicDatatype], StorNam, DispName, Description.
//   - TODO: (Maybe): the field [ColumnSpecs] could be nil or
//     len 0. If so then it should "probably" be autogenerated
//     (incompletely. tho) by reflection from the contents of
//     a same-named table already existing in the DB.
//
// IMPORTANT NOTE ABOUT THIS STRUCT
//
// This struct contains four key items that MUST be kept 
// in sync in terms of DB columns defined, struct fields 
// defined, and their order of appearance in DB statements:
//  - ColumnSpecs
//  - ColumnNamesCSV
//  - ColumnPtrsFunc
//  - the struct corresponding to them all 
//
// If they are not kept in sync as required, there may be compile
// time errors, but more likely there will be nasty runtime errors,
// and in particular, obscure impenetrable DB errors.
//
// (In principle this can all be done with reflection and code
// generation, starting from the column specs, but it would be
// a big implementation effort.) 
//
// For typical DB operations like INSERT and UPDATE (altho not DELETE),
// the funcs and vars in this file always explicitly pass the names of
// ALL fields/columns to the DB engine, and (for example) the shortcut
// "SELECT *" is never used.
// 
// This means that fields can be defined in any order freely but must 
// be referred to consistently in the same order by funcs and vars.
// Note tho that DB tools will have trouble displaying (on screens)
// the values of ALL fields, and for this reason, fields whose values
// are shorter and more important should appear (i.e. be defined) first,
// so that CREATE TABLE lists them first, ahead of much-longer and/or 
// less-important fields.
// 
// Notes on date-time fields:
//   - SQLite can use string, int, or real. But, date-time fields
//     based on TableDetails and [dsmnd.DbColumnSpec]) use strings
//     (SQLite "TEXT"), which are expected to be ISO-8601 (or RFC
//     3339) (and probably UTC). Text is the first option listed in
//     https://www.sqlite.org/datatype3.html#date_and_time_datatype:
//   - TEXT: "YYYY-MM-DD HH:MM:SS.SSS"
//   - REAL as Julian day numbers: the day count since 24 Nov 4714 BC
//   - INTEGER as Unix time: the seconds count since 1970-01-01 00:00:00 UTC
//   - NOTE: For TEXT "YYYY-MM-DD HH:MM:SS.SSS", this might often end up
//     in ISO format, which has a "T" instead of the blank " ". So for
//     better readability, and to avoid line breaks, we have a utility
//     that replaces either a blank (" ") or an ISO "T" with a "_".
//
// TODO: Rename struct fields to start with lower case,
// and provide exported accessors ("getters"). 
//
// TODO: Some day all this nitpicky stuff about keeping fields
// in sync could all be done with code generation. 
// .
type TableDetails struct {

     	// [dsmnd.TableSummary] is a [dsmnd.Datum] 
	// and has four fields, used thusly:
	//  - [dsmnd.BasicDatatype]: [D.SCT_TABLE] 
	//  - StorName: the name of the table in the DB,
	//    e.g. "inbatch", "contentity", "topicref" 
	//  - DispName: a short (three-letter) version 
	//    of the name to embed in the names of other
	//    variables, e.g. inb, cnt, trf 
	//  - Description: as is appropriate 
     	D.TableSummary
	
	// PKname is the (auto-generatable as "idx_foo"!) name 
	// of the index (i.e. primary key) field, which we use 
	// in the same format BOTH as primary key in own-table 
	// AND as foreign key in other tables. Enabling natural
	// joins, without using "AS"! 
	PKname string
 
	// ColumnNamesCSV is all column names (except primary key),
	// ready-to-use in SQL, in a specific (auto-generatable!)
	// order, comma-separated. We omit the primary key so that
	// we can use it for SQL INSERT statements too. REPLACED!! 
	// ColumnNamesCSV string
	
	// ColumnSpecs is a list of [dsmnd.D.ColumnSpec] that omits
	// the primary key (which can be brought in when needed).
	ColumnSpecs []D.ColumnSpec
	
	// ColumnPtrsFunc return a slice of ptrs to every field 
	// in the passed-by-ptr struct (that implements interface
	// [Rowmodel]), and sometimes includes the primary key.
	// The slice is used for database Scan(..) funcs.
	// 
	// Note that ColumnPtrsFunc is a func, and/but there 
	// is also this method in interface [RowModel]:
	// (p RowModel) ColumnPtrsMethod(bool) []any
	// .
	ColumnPtrsFunc columnPtrsFunc

	// BlankInstance might be needed at some point
	// BlankInstance RowModel
	// NewInstance returns a ptr to just that, 
	// and is useful for resolving generics. 
	NewInstance NewInstanceFunc
	
	// We used to have ForenKeys defined by name only, but
	// this was insufficient information, because we need
	// the table name AND (sometimes) a unique field name. 
	// ForenKeys   []string

	// FuncNew func() RowModel

	CSVs *ColumnStringsCSV
	PrepdStmts *PreparedStatements
}

