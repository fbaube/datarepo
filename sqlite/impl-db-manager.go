package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	_ "github.com/mattn/go-sqlite3" // to get init()
	L "github.com/fbaube/mlog"
	"os"
	FP "path/filepath"
	S "strings"
)

// Methods in this file implement
// type DBManager interface {
//	OpenAtPath(string) (string, error)
//	NewAtPath(string) (string, error)
//	OpenExistingAtPath(string) (string, error)

/* OpenFile flags etc.
// Flags to OpenFile wrapping those of the underlying system.
// Not all flags may be implemented on a given system.
const (
    // Specify exactly one of O_RDONLY, O_WRONLY, or O_RDWR
    O_RDONLY int = syscall.O_RDONLY // open the file R/O
    O_WRONLY int = syscall.O_WRONLY // open the file write-only
    O_RDWR   int = syscall.O_RDWR   // open the file R/W
    // The remaining values may be or'ed in to control behavior.
    O_APPEND int = syscall.O_APPEND // append data to the file when writing.
    O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
    O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
    O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
    O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file after open.
)

when file already exists, either truncate it or fail:

openOpts := os.O_RDWR|os.O_CREATE
if OKtoTruncateWhenExists {
    openOpts |= os.O_TRUNC // file will be truncated
} else {
    openOpts |= os.O_EXCL  // file must not exist
}
f, err := os.OpenFile(filePath, openOpts, 0644)
// ... do stuff
*/

/* OS errors
https://pkg.go.dev/os
var (
	// ErrInvalid indicates an invalid argument.
	// Methods on File will return this error when the receiver is nil.
	ErrInvalid = fs.ErrInvalid // "invalid argument"

	ErrPermission = fs.ErrPermission // "permission denied"
	ErrExist      = fs.ErrExist      // "file already exists"
	ErrNotExist   = fs.ErrNotExist   // "file does not exist"
	ErrClosed     = fs.ErrClosed     // "file already closed"

	ErrNoDeadline = errNoDeadline()  // "file type does not support deadline"

	ErrDeadlineExceeded = errDeadlineExceeded() // "i/o timeout"
)

https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
Instead of using os.Create, you should use
os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666) .
That way you'll get an error if the file already exists.
Also this doesn't have a race condition with something
else making the file, unlike a version that checks for
existence beforehand.
https://groups.google.com/g/golang-nuts/c/Ayx-BMNdMFo/m/4rL8FFHr8v4J
Often an os.Exists function is not really needed.
For instance: if you are going to open the file, why check
whether it exists first? Simply call os.IsNotExist(err) after
you try to open the file, and deal with non-existence there.

os.IsExist(err) is good for cases when you expect the
file to not exist yet, but the file actually exists:
os.Symlink, os.Mkdir,
os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
(os.IsExist(err) will trigger when target exists because
O_EXCL means that file should not exist yet)

os.IsNotExist(err) is good for more common cases where you
expect the file to exists, but it actually doesn't exist:
os.Chdir, os.Stat, os.Open, os.OpenFile(without os.O_EXCL),
os.Chmod, os.Chown, os.Close, os.Read, os.ReadAt, os.ReadDir,
os.Readdirnames, os.Seek, os.Truncate
*/

func (p *SqliteDBManager) OpenAtPath(path string) (*SqliteRepo, error) {
	/* N'YET
	if p.DBImplementationName != D.DB_SQLite { panic("not sqlite ?!") }
	*/
	// func IsFileAtPath(aPath string) (bool, *os.FileInfo, error) {
	filexist, fileinfo, filerror := FU.IsFileAtPath(path)
	if filexist && (fileinfo.Size() == 0) {
		// Bye!
		L.L.Warning("DB file is empty (0-len): deleting it.")
		e := os.Remove(path)
		if e != nil {
			return nil, fmt.Errorf("sqlite.OpenAtPath: "+
				"os.Remove(%s): %w", path, e)
		}
	}
	if filexist && (filerror == nil) {
		return p.OpenExistingAtPath(path)
	}
	if (!filexist) && (filerror == nil) {
		return p.NewAtPath(path)
	}
	return nil, fmt.Errorf("sqlite.OpenAtPath: unfathomable: %s", path)
}

// NewAtPath creates a DB at the filepath, opens it, and runs
// standard initialization pragma(s). It does not create any tables
// in it. If a file or dir already exists at the filepath, the func
// returns an error. The filepath can be a relative path, but not "".
//
// The repo type will be "sqlite" (equivalent to "sqlite3").
// .
func (p *SqliteDBManager) NewAtPath(aPath string) (*SqliteRepo, error) {
// func (p *SqliteDBManager) NewAtPath(aPath string) (DR.SimpleRepo, error) {

	if !S.HasSuffix(aPath, ".db") {
		L.L.Warning("newSqliteRepo: missing suffix \".db\": " + aPath)
	}
	if !FP.IsAbs(aPath) {
		// println("BEFOR", aPath)
		aPath = FU.ResolvePath(aPath)
		// println("AFTER", aPath)
	}
	// func IsFileAtPath(aPath string) (bool, *os.FileInfo, error) {
	filexist, fileinfo, filerror := FU.IsFileAtPath(aPath)
	// Handle this case separately, cos it requires user action:
	// something there, but not a regular file ?
	if fileinfo != nil && !filexist {
		return nil, errors.New(
			"sqlite.repo.new: blocked, can't create: " + aPath)
	}
	if filexist || filerror != nil || fileinfo != nil {
		return nil, fmt.Errorf(
			"sqlite.repo: New(%s): problem(s)...", aPath)
	}
	// ALL CLEAR !
	errPfx := fmt.Sprintf("sqlite.repo.new(%s): ", aPath)
	var pDB *sql.DB
	// Try opening it, and then close it immediately.
	file, e := // not os.Create(aPath) but rather
		os.OpenFile(aPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if e != nil {
		return nil, fmt.Errorf(errPfx+"tried os.create, got: %w", e)
	}
	file.Close()

	// It may seem odd that this is necessary, but for
	// some lame retro compatibility, SQLite does not
	// default to enforcing foreign key constraints.
	// NOTE: Now we do this in the connection string.
	// pSR.DoPragma("foreign_keys = ON")

	// pDB, e = sqlx.Open("sqlite3", "file:"+aPath+"?foreign_keys=on")
	if !S.HasPrefix(aPath, "file:") {
		// println("Adding to file URL the prefix \"file:\"")
		aPath = "file:" + aPath
	}
	pDB, e = sql.Open("sqlite3", aPath+"?foreign_keys=on")
	if e != nil {
		return nil, fmt.Errorf(errPfx+"sql.open: %w", e)
	}
	L.L.Info("DB file is open: " + aPath)
	// https://golang.org/pkg/database/sql/#Open
	// Open might simply validate its arguments without creating a DB
	// connection. To verify that the data source name is valid, call Ping.
	if e = pDB.Ping(); e == nil {
		if e = pDB.PingContext(context.Background()); e != nil {
			return nil, fmt.Errorf(
				"reposqlite.new: sql.Ping[Context]: %w", e)
		}
	}
	L.L.Info("DB ping worked")
	// pDB = dbx.NewFromDB(psDB, "sqlite")

	L.L.Info("New DB created at: " + SU.Tildotted(aPath))
	// drivers := sql.Drivers()
	// println("DB driver(s):", fmt.Sprintf("%+v", drivers))
	pSR := new(SqliteRepo)
	pSR.DB = pDB
	pSR.filepath = aPath
	pSR.w = io.Discard
	return pSR, nil
}

// OpenExistingAtPath opens an existing DB at the filepath. It is
// assumed to have any default initializations (i.e. pragmas etc.).
// If no file already exists at the filepath, the func returns an
// error. The filepath can be a relative path, but may not be "".
//
// The repo type will be "sqlite" (equivalent to "sqlite3").
// .
func (p *SqliteDBManager) OpenExistingAtPath(aPath string) (*SqliteRepo, error) {
	if !S.HasSuffix(aPath, ".db") {
	   L.L.Warning("SqliteRepo.openAt: missing \".db\": " + aPath)
	}
	aPath = FU.ResolvePath(aPath)
	// func IsFileAtPath(aPath string) (bool, *os.FileInfo, error) {
	filexist, fileinfo, filerror := FU.IsFileAtPath(aPath)
	if filerror != nil {
		return nil, fmt.Errorf(
			"sqlite.repo.openAt(%s): %w", aPath, filerror)
	}
	var e error
	// file already exists ? HAPPY PATH
	if filexist {
		// ALL CLEAR !
		pSR := new(SqliteRepo)
		pSR.DB, e = sql.Open("sqlite3", aPath+"?foreign_keys=on")
		// pSR.DB = sqlx.MustConnect("sqlite3", aPath)
		pSR.filepath = aPath
		pSR.w = io.Discard
		if e == nil {
			return pSR, nil
		}
		// OOPS, our DB file mysteriously went OTL
		return nil, fmt.Errorf(
			"sqls.connect(%s): file went OTL: %w", aPath, e)
	}
	// something there, but not a regular file ?
	if fileinfo.Mode() != 0 { // Weird syntax complaint from compiler
		return nil, errors.New(
			"sqlite.repo.openAt: blocked, can't open: " + aPath)
	}
	// nothing there
	return nil, errors.New(
		"sqlite.repo.openAt: file not found: " + aPath)
}
