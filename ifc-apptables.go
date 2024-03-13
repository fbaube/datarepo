package datarepo

import (
	// _ "database/sql"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// RepoAppTables is table-related methods for a specified app's schema.
// The app name is case-insensitive, and used as all lower case, and
// prefixed to table names as "appname_". If the app name is left blank
// (""), a default namespace is used and no prefix is added to table names.
type RepoAppTables interface {
	// SetAppTables specifies the schemata of the specified app's
	// tables, which this interface creates and/or manages. Multiple
	// calls, whether with tables previously specified or not before
	// seen do not conflict; if a table name is repeated but with a
	// different schema, the result is undefined.
	SetAppTables(string, []DRM.TableDetails) error
	// EmptyAllTables deletes (app-level) data from the app's tables
	// but does not delete any tables (i.e. no DROP TABLE are done).
	EmptyAppTables() error
	// CreateTables creates the app's tables per already-supplied
	// schema(ta); if the tables exist, they are emptied of data.
	CreateAppTables() error
}
