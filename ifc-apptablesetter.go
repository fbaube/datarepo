package datarepo

import (
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// AppTableSetter is table-related methods for a specified app's schema.
// 
// An untested feature is that an app name prefix can be added, via 
// the first argument to func [RegisterAppTables]. The app name is
// case-insensitive, and used as all lower case, and prefixed to
// table names as "appname_". If the app name is left blank (""), a
// default namespace is used and no prefix is added to table names.
// .
type AppTableSetter interface {
	// RegisterAppTables processes the schemata of the specified 
	// app's tables, which this interface creates and/or manages.
	//  - Multiple calls to this func do not conflict, whether 
	//    with tables previously specified or not before seen 
	//  - If a table name is repeated but with a different schema,
	//    the result is undefined
	//  - If the tables already exist in the DB, it is not verified
	//    that their structure matches what this schema specifies
	//    (but this might be a future TODO) 
	RegisterAppTables(string, []*DRM.TableDetails) error
	// EmptyAllTables deletes (app-level) data from the app's tables
	// but does not delete any tables (i.e. no DROP TABLE are done).
	EmptyAppTables() error
	// CreateTables creates the app's tables per already-supplied
	// schema(ta); if the tables exist, they are emptied of data. 
	CreateAppTables() error
	
	// GetTableDetailsByString is case-insensitive.
	// Also, we do not define it here, because it is not
	// DB-specific, it is expected for all impementations
	// of SimpleRepo and interface AppTableSeteter.
	// GetTableDetailsByString(string) *DRM.TableDetails
}

