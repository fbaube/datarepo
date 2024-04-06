package datarepo

// Backupable methods are invoked directly by the DB being backed up,
// rather than by some sort of higher-level "Manager". The methods
// work with locations, whose type (filepath, dir path, or URI/URL)
// and naming convention (incorporating date & time) are determined
// by the implementation for each DB. Methods exist to move DB to,
// copy DB to, or restore DB from a location. Each method returns
// the location of the new backup or restored-from backup.
//
// There are equivalent sqlite3 commands at the CLI: 
//   - sqlite3 my_database.db ".backup 'backup_file.db'"
//   - sqlite3 my_database.db ".backup m_database.db.bak"
//   - sqlite3 my_database .backup > my_database.back
// .
type Backupable interface {
	MoveToBackup() (string, error)
	CopyToBackup() (string, error)
	RestoreFromMostRecentBackup() (string, error)
}
