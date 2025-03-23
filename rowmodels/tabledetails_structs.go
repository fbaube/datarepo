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
// application's data schema both at the table level and
// field-by-field. A primary key is assumed for every table,
// in the zeroth column, and foreign keys are allowed.
//
// NOTE: The case of multiple foreign indices into the same 
// table (from TopicrefRow into ContentityRow) is 90% resolved. 
//
// As a convention, all DB table and column names should be 
// in all lower case, except perhaps in the names of indices, 
// which have names of the form (??) ID_* or Idx_*. Then SQL 
// keywords can be given in all upper case, per SQL convention. 
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
//  - ColumnNamesCSV (2025.q1: OBS: is now auto-gen'd)
//  - ColumnPtrsFunc
//  - the struct corresponding to them all 
//
// If they are not kept in sync as required, there may be compile
// time errors, but more likely there will be nasty runtime errors,
// and in particular, obscure impenetrable DB errors.
//
// For typical DB operations like INSERT and UPDATE (altho not DELETE),
// the funcs and vars in this file always explicitly pass the names of
// ALL fields/columns to the DB engine, and (for example) the shortcut
// "SELECT *" is never used.
// 
// This means that fields can be defined in any order freely but must 
// be referred to consistently in the same order by funcs and vars.
// EXCEPTION: the table's primary key must be in the very first (i.e.
// zeroth) column,
// 
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

	// ColumnStringsCSV is described in full elsewhere.
	CSVs *ColumnStringsCSV
	// Statements is described in full elsewhere.
	Stmts *Statements
}

// ColumnStringsCSV stores strings useful for composing SQL
// statements. Each string includes all the columns, in order,
// comma-separated. SQL using these strings defaults to setting
// and getting every field in a DB record.
//
// The strings have no trailing commas. Each string (except 
// for UPDATE) has two versions:
//  - a version (suffixed with "_wID") that DOES include the primary 
//    key (always named "{table}_ID"), for (e.g.) output from SELECT
//  - a version (suffixed with "_noID") that does NOT include the 
//    primary key, for (e.g.) input to INSERT (where the ID is new)
//    and input to UPDATE (where the ID is used in a WHERE clause
//    to find the record)
// .
type ColumnStringsCSV struct {

	// FieldNames_* [with|no ID primary key] 
	// list column (i.e. field) names, in order:
	// (w/ ID) "F0=ID, F1, F2, F3, F4" 
	// (no ID)        "F1, F2, F3, F4" 
	FieldNames_wID, FieldNames_noID string
	
	// PlaceNums_* [with|no ID primary key] 
	// list '$'-numbered parameters:
	// (w/ ID) "$1, $2, $3, $4, $5" 
	// (no ID)     "$1, $2, $3, $4"
	// (where) "$1, $2, $3, $4, $5, $6, $7" (w/ ID) 
	PlaceNums_wID, PlaceNums_noID, PlaceNums_wID_wFV string
	
	// FieldUpdates [_noID: NOT with ID primary key] 
	// lists column/field names with "=", and their 
	// new values as '$'-numbered parameters:
	// "F1 = $1, F2 = $2, F3 = $3, F4 = $4"
	UpdateNames string

	// Where_* [with|no ID primary key] lists
	// FIELD NAMES ?? PLACEHOLDERS ?? 
	// for when there is also a WHERE clause:
	Where_noVals, Where_wVals_wID, Where_wVals_noID string 
/*
	pTD.CSVs.Where_noVals     = fmt.Sprintf(" WHERE $%d = $%d", 1, 2)
	pTD.CSVs.Where_wVals_wID  = fmt.Sprintf(" WHERE $%d = $%d", N, N+1)
	pTD.CSVs.Where_wVals_noID = fmt.Sprintf(" WHERE $%d = $%d", N-1, N)
*/
}

// Statements stores several SQL query strings customised for the
// table. Statements vary in whether they include the primary key, 
// and whether they include a WHERE clause.
//
// Statements named "*unique" are for working with single records,
// and are used by method [datarepo.EngineUnique] of interface
// [datarepo.DBEnginer].
// .
type Statements struct {
     	INSERTunique string
	SELECTunique string
	UPDATEunique string
	DELETEunique string 
}

