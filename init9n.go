package datarepo

import (
	"errors"
	"fmt"
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	// DRS "github.com/fbaube/datarepo/sqlite"
	DRM "github.com/fbaube/datarepo/rowmodels"
	_ "github.com/mattn/go-sqlite3"
	SU "github.com/fbaube/stringutils"
	"os"
	_ "database/sql"
)

// Init9nArgs is for database management. Note that if field [UseDB]
// is false, the contents and usage of this struct are undefined.
type Init9nArgs struct {
     // D.DB_type is so far only D.DB_SQLite = "sqlite"
     D.DB_type
     // BaseFilename defaults to "m5.db"
     BaseFilename string 
     Dir string
     // LogWriter is initialized to [io.Discard]. If it is an open file, 
     // it is passed using func [datarepo/SimpleRepo.SetLogWriter].
     // It could actually be just an io.WriterCloser. 
     // LogWriter *os.File
     
     // DoImport requires DB access, so it is present
     // here, but it is not otherwise processed here. 
     DoImport bool
     // DoZeroOut says initialize the DB with the 
     // app's tables but with no data in them 
     DoZeroOut bool
     // DoBackup says before DoingZeroOut on an
     // existing DB, first copy it to a backup
     // copy using a hard-coded naming scheme 
     DoBackup bool
     // TableDetailz are app tables' details 
     TableDetailz []*DRM.TableDetails
}

var DEFAULT_FILENAME = "m5.db"

// ProcessInit9nArgs processes DN initialization arguments. 
// It can process either a new DB OR an existing DB.
//
// If the returned error is non-nil, it need not be fatal.
// Check whether the SimpleRepo return value is nil. 
// 
// TODO: It should not use a complex logger (e.g. [mlog]),
// because we want to avoid that kind of dependency in 
// a standalone library. 
// .
func (p *Init9nArgs) ProcessInit9nArgs() (SimpleRepo, error) {

	var mustAccessTheDB bool
	var e error
	// This will always be true, until it is decided
	// under what conditions p.Dir might be "". Code
	// below presumes that "" means "." 
	mustAccessTheDB = p.DoImport || p.DoZeroOut || p.Dir != ""
	if !mustAccessTheDB {
		return nil, nil 
	}
	if p.DB_type != "" && p.DB_type != D.DB_SQLite {
	   return nil, errors.New("DB: bad type: " + string(p.DB_type))
	}
	if p.DB_type == "" {
	   println("DB: type is missing: using SQLite.")
	   p.DB_type = D.DB_SQLite
	}
	// ===================================================
	// Start by checking on the status of the DB filename.
	// NOTE: This assumes the DB is SQLite, a single file.
	// ===================================================
	// Derive a FILE path.
	// ===================
	var dbFilepath string
	// NOTE that if p.Dir is "", ResolvePath won't fix it!
	if p.Dir == "" {
		p.Dir = "."
	}
	dbFilepath = FU.ResolvePath(
		p.Dir + FU.PathSep + DEFAULT_FILENAME)
	println("DB: full path:", dbFilepath)
	// ================================
	// Declare vars and check for file.
	// ================================
	errPfx := "db.init9n(" + dbFilepath + "): "
	var fileinfo os.FileInfo
	filexist, fileinfo, filerror := FU.IsFileAtPath(dbFilepath)
	if filerror != nil {
		return nil, fmt.Errorf(errPfx + "file error: %w", filerror)
	}
	var dbShortPath string 
	dbShortPath = SU.Tildotted(dbFilepath)
	// =========================================
	// If it exists already but is zero-len then
	// remove it, even tho this may remove any
	// special permissions that were configured.
	// =========================================
	if filexist && fileinfo.Size() == 0 { 
		L.L.Warning(errPfx + "exists but is 0-len: removing")
		e = os.Remove(dbFilepath)
		if e != nil {
			panic(e)
		}
		filexist = false
	}
	// ===============================
	// Now let's get down to business.
	// ===============================
	// Prepare a DB logging file.
	// ==========================
	dbw, dbe := FU.CreateEmpty("./db.log")
	if dbe != nil {
	   L.L.Warning("Cannot open DB logfile ./db.log: %w", dbe)
	}
	// =================
	// Declare repo var.
	// =================
        var repo SimpleRepo // basic var 
	// =====================================
	// Register apptable details.
	// (Not the same as creating apptables.)
	// =====================================
	if p.TableDetailz == nil || len(p.TableDetailz) == 0 {
	   println("DB: missing app table details. Aborting.")
	   return nil, errors.New("missing app DB table details")
	}
	// =====================================
	// Now it's time to work on the DB file. 
	// =====================================
	// Happy path: it exists,
	// so do stuff and return.
	// =======================
	if filexist {
		repo, e = DB_Manager.OpenAtPath(dbFilepath)
		if e != nil {
			return nil, fmt.Errorf(errPfx + "can't open: %w", e)
		}
		repo.SetLogWriter(dbw)
		e = repo.RegisterAppTables("", p.TableDetailz) 
		if e != nil {
		   return nil, fmt.Errorf(errPfx + "can't register tables: %w", e)
		}
		// println(fmt.Sprintf("repo <%T> \n", repo))
		// =====================================
		// If it exists and we want to open it
		// as-is, i.e. without zeroing it out,
		// then this is where we return success.
		// =====================================
		if !p.DoZeroOut {
		   L.L.Info("DB file opened OK: " + dbShortPath)
		   return repo, nil
		}
		// p.DoZeroOut 
		L.L.Info("(backing up and) zeroing out DB")
		_, e := repo.CopyToBackup()
		if e != nil {
			panic(e)
		}
		e = repo.EmptyAppTables()
		L.L.Info("DB file emptied OK: " + dbShortPath)
		return repo, e 
	}
	// =========
	// !filexist 
	// =========
   	println("DB: Creating anew.")
	L.L.Info("Creating DB: " + dbShortPath)
	if p.DoZeroOut {
		L.L.Info("new DB: zeroing it out is redundant")
	}
	// "NewAtPath creates a DB at the filepath, opens it, 
	// and runs standard initialization pragma(s). It does
	// not create any tables in it. If a file or dir already
	// exists at the filepath, the func returns an error."
	repo, e = DB_Manager.NewAtPath(dbFilepath)
	if e != nil {
		return nil, fmt.Errorf(errPfx + "new() failure: %w", e)
	}
	repo.SetLogWriter(dbw)
	e = repo.RegisterAppTables("", p.TableDetailz) 
	if e != nil {
		return nil, fmt.Errorf(errPfx + "can't register tables: %w", e)
	}
	// println(fmt.Sprintf("repo <%T> \n", repo))
	// ===========================
	// Configure it using pragmas.
	// ===========================
	ps, pe := repo.DoPragmas(DB_Manager.InitznPragmas())
	if  pe != nil {
	    L.L.Warning("Expect DB problems: " +
                       "DB init pragmas got error: %w", pe)
	    }
	L.L.Info("init pragmas on DB returned: " + ps)
	// ======================================	
	println("DB: new: OK.")
	L.L.Info("DB new() OK: " + repo.Path())
	e = repo.CreateAppTables()
	return repo, e
}

/*
// inputExts more than covers the file types associated with the LwDITA spec.
// Of course, when we check for them we do so case-insensitively.
// Ones we don't handle don't go here, e.g. ".adoc" Asciidoc.
var inputExts = []string{
	".dita", ".map", ".ditamap", ".xml",
	".md", ".markdown", ".mdown", ".mkdn", 
	".html", ".htm", ".xhtml", ".png", ".gif", ".jpg"}

// AllGLinks gathers all GLinks in the current run's input set.
var AllGLinks mcfile.GLinks
*/

