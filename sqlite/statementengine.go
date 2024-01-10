package sqlite

// StatementEngine implements interface [StatementBuilder].
// It is DB-specific (in this case, SQLite) and also specific
// to the tables that it is configured with (cnt, inb, trf).
//
// Each table's configuration (via [TableDetails]) contains
// these key data items:
//  - [ColumnSpecs], of type sliceOf [ColumnSpec]), which
//    specifies the [SemanticFieldType], which determines
//    the SQLite storage type (INT, INTEGER, BOOL, BLOB, ...)
//  - ColumnNameCSVs, of type string, which lists all of the
//    table's fields (except primary ID), which is used to
//    build SQL statements that do not use the "*" wildcard 

