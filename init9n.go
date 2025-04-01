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
	   return nil, errors.New("bad DB type: " + string(p.DB_type))
	}
	if p.DB_type == "" {
	   println("DB: Type is missing: using SQLite.")
	}

	// ===================================================
	// Start by checking on the status of theDB  filename.
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
        var repo SimpleRepo
	// type-checking
	// var _ repo.SimpleRepo = (*RS.SqliteRepo)(nil)
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
		L.L.Warning("DB exists but 0-len: " + dbShortPath)
		e = os.Remove(dbFilepath)
		if e != nil {
			panic(e)
		}
		filexist = false
	}
	// ==========================================
	// Now we have (or will have) a repo and we 
	// will use it, so provide a DB logging file.
	// ==========================================
	dbw, dbe := FU.CreateEmpty("./db.log")
	if dbe != nil {
	   L.L.Warning("Cannot open DB logfile ./db.log: %w", dbe)
	}
	// ===============================
	// Now let's get down to business.
	// ===============================
	var crearror error 
	if filexist {
		repo, crearror = DB_Manager.OpenAtPath(dbFilepath)
	} else {
	   	println("DB: Creating anew.")
		L.L.Info("Creating DB: " + dbShortPath)
		if p.DoZeroOut {
			L.L.Info("Zeroing out the DB is redundant")
		}
		repo, crearror = DB_Manager.NewAtPath(dbFilepath)
		// Now configure it using pragmas
		ps, pe := repo.DoPragmas(DB_Manager.InitznPragmas())
		if pe != nil {
		   L.L.Warning("Expect DB problems: " +
                        "DB init pragmas got error: %w", pe)
			}
		L.L.Info("Ran init pragmas on DB: " + ps)
	}
	if crearror != nil {
		return nil, fmt.Errorf(errPfx + "DB failure: %w", crearror)
	}
	// ======================================
	// Soon we can maybe return success, but
	// before that, register apptable details.
	// ======================================
	var pSR SimpleRepo
	pSR, ok := repo.(SimpleRepo)
	if !ok {
		// panic("init9n.L155")
		return nil, errors.New("db.init9n: is not sqlite simplerepo")
	}
	if p.TableDetailz == nil || len(p.TableDetailz) == 0 {
	   println("DB: missing app table details. Aborting.")
	   return nil, errors.New("Missing app DB table details")
	}
	e = pSR.RegisterAppTables("", p.TableDetailz) 
	if e != nil {
		return nil, fmt.Errorf(errPfx + "can't register tables: %w", e)
	}
	repo.SetLogWriter(dbw)
	// ======================================	
	// If the DB exists and we want to open
	// it as-is, i.e. without zeroing it out,
	// then this is where we return success.
	// ======================================
	if filexist && repo != nil && !p.DoZeroOut {
		L.L.Info("DB file opened OK: " + dbShortPath)
		return repo, nil
	}
	println("DB: status OK.")
	L.L.Info("DB OK: " + repo.Path())
	// ===================================
	// Otherwise, there is a bit more work
	// to do before we return success.
	// ===================================
	if !filexist {
		e = pSR.CreateAppTables()
	} else if p.DoZeroOut {
		L.L.Info("Zeroing out DB (init9nb.go)")
		_, e := repo.CopyToBackup()
		if e != nil {
			panic(e)
		}
		e = pSR.EmptyAppTables()
	}
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

