package datarepo

import "io"

// SessionLifecycler is session lifecycle operations for databases.
type SessionLifecycler interface {
	// Open is called on an existing repo file, and can be called
	// multiple times in a sessions, so it should not pragma-style
	// initialization; however, options passed in the connection
	// string (such as SQLite's "...?foreign_keys=on") are kosher.
	Open() error
	SetLogWriter(io.Writer) io.Writer
	// TODO:? return it, for use or understanding
	// InitString() string
	// IsOpen pings the DB as a health check.
	IsOpen() bool
	// Verify runs app-level sanity & consistency checks (but things
	// like foreign key integtrity should be delegated to DB setup).
	Verify() error
	Flush() error
	// Close remembers the path (like os.File does).
	Close() error
	CloseLogWriter()
}
