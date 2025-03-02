package datarepo

import (
       "io"
	_ "github.com/mattn/go-sqlite3"
	// DRM "github.com/fbaube/datarepo/rowmodels"
)

// SimpleRepo is an interface that combines several 
// other interfaces, and can be fully described by
//  1. an implementation, currently limited 
//     to [dsmnd.DB_SQLite] ("sqlite")
//  2. a filepath or a URL, for opening or creating, 
//     which may be either relative or absolute 
//
// Each field in the struct is tipicly a ptr, and 
// tipicly they all point to the same single object.
//
// A SimpleRepo is expected to implement DBBackups.
// .
type SimpleRepo interface {
     // Entity is Repo type, path, etc.
	Entity
     // Backupable is copy, move, restoreFrom 
	Backupable
     // SessionLifecycler is open, close, etc.
	SessionLifecycler
     // StatementBuilder uses [TableDescriptor] and [QuerySpec] 
	StatementBuilder

     // Caller is basic DB access operations but FIXME: uses generics.
     // So maybe have to use typecasting to access these methods.
     // Caller[DRM.InbatchRow]
     	
     // AppTableSetter is table mgmt for a specific Repo-using app.
	AppTableSetter
	
	DBEnginer

	LogWriter() io.Writer
}

