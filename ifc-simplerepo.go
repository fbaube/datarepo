package datarepo

import (
       "io"
	_ "github.com/mattn/go-sqlite3"
)

// SimpleRepo is an interface that combines several 
// other interfaces, and can be fully specified by
//  1. an implementation, currently limited to
//     [dsmnd.DB_SQLite] ("sqlite"), plus 
//  2. a filepath or a URL, for opening or creating, 
//     which may be either relative or absolute 
//
// Each field in the struct is tipicly a ptr, and 
// tipicly they all point to the same single object.
//
// A SimpleRepo is expected to implement DBBackups.
// .
type SimpleRepo interface {
     // Entity is the Repo's basics: type, path, etc.
	Entity
     // Backupable is copy, move, restoreFrom 
	Backupable
     // SessionLifecycler is open, close, etc.
	SessionLifecycler
     // StatementBuilder uses [TableDescriptor] and [QuerySpec] 
	StatementBuilder

     // OBS? Caller is basic DB access ops but FIXME: it uses generics.
     // So maybe have to use typecasting to access these methods.
     // Caller[DRM.InbatchRow]
     	
     // AppTableSetter is table mgmt for a specific Repo-using app.
	AppTableSetter
     // DBEnginer is basic DB ops. 
	DBEnginer
	LogWriter() io.Writer
}

