package datarepo

// Backupable methods work with locations, whose type (filepath, dir
// path, or URI/URL) and naming convention (incorporating date & time)
// are determined by the implementation for each DB. Methods exist to
// move DB to, copy DB to, or restore DB from a location. Each method
// returns the location of the new backup or restored-from backup.
//
// It defines methods on a (probably empty) singleton that selects
// (for now: only) SQLite implementations of this interface.
//
// At the CLI?
//   - sqlite3 my_database.sq3 ".backup 'backup_file.sq3'"
//   - sqlite3 m_database.sq3 ".backup m_database.sq3.bak"
//   - sqlite3 my_database .backup > my_database.back
//
// .
type Backupable interface {
	MoveToBackup() (string, error)
	CopyToBackup() (string, error)
	RestoreFromMostRecentBackup() (string, error)
}
