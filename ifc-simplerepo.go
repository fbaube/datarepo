package datarepo

import (
	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/fbaube/sqlite3"
)

// SimpleRepo is an interface that combines several 
// other interfaces, and can be fully described by
//  1. (NOTE: OBS?) a [DBImplementation], 
//     currently limited to "sqlite", plus 
//  2. a filepath or a URL, which may 
//     be either relative or absolute 
//
// Each field in the struct is tipicly a ptr, and 
// tipicly they all point to the same single object.
//
// A SimpleRepo is expected to implement DBBackups.
// .
type SimpleRepo interface {
     // Implementation is func [DBImplementationName]
     // is currently limited to "sqlite"
     // Implementation
     // Entity is type, path, etc.
	Entity
     // Backupable is copy, move, restoreFrom 
	Backupable
     // SessionLifecycler is open, close, etc.
	SessionLifecycler
     // StatementBuilder uses [TableDescriptor] and [QuerySpec] 
	StatementBuilder
     // QueryRunner is for generics and has funcs that return 0,1,N rows
	QueryRunner

	RepoAppTables
	// DBManager
}

