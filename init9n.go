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

        var repo SimpleRepo
	// type-checking
	// var _ repo.SimpleRepo = (*RS.SqliteRepo)(nil)

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
	
	// Start by checking on the status of the filename.
	// NOTE: This assumes that the DB is SQLite, a single file.
	// Note that a path is used to derive a FILE path.
	var dbFilepath string
	// println("misc.go: BEFOR:", p.Dir)
	// NOTE that if p.Dir is "", ResolvePath won't fix it!
	if p.Dir == "" {
		p.Dir = "."
	}
	// println("misc.go: BEFOR:", p.Dir)
	dbFilepath = FU.ResolvePath(
		p.Dir + FU.PathSep + DEFAULT_FILENAME)
	println("DB: full path:", dbFilepath)
	errPfx := fmt.Errorf("DB.Init9n(%s):", dbFilepath)

	var fileinfo os.FileInfo
	filexist, fileinfo, filerror := FU.IsFileAtPath(dbFilepath)
	if filerror != nil {
		return nil, fmt.Errorf("%s file error: %w", errPfx, filerror)
	}
	var shortPath string 
	shortPath = SU.Tildotted(dbFilepath)
	var openError error
	// If the DB already exists already and 
	// is non-zero then go ahead and open it
	if filexist {
		L.L.Info("DB exists: " + shortPath)
		// If the file is zero-length, it is removed, altho this may
		// remove some special permissions that were configured
		if fileinfo.Size() == 0 {
			L.L.Info("DB is empty: " + shortPath)
			e = os.Remove(dbFilepath)
			if e != nil {
				panic(e)
			}
			filexist = false
		} else {
		        repo, openError = DB_Manager.OpenAtPath(dbFilepath)
		}
	}
	// Now we have (or will have) a repo and we 
	// will use it, so provide a DB logging file.
	dbw, dbe := FU.CreateEmpty("./db.log")
	if dbe != nil {
	   L.L.Warning("Cannot open DB logfile ./db.log: %w", dbe)
	}
	// We are soon ready to maybe return success, but before we can
	// do that, we have to be sure to register apptable details.
	var pSR SimpleRepo
	pSR, ok := repo.(SimpleRepo)
	if !ok {
		panic("init9n.L165")
		return nil, errors.New("processDBargs: is not sqlite")
	}
	
	if p.TableDetailz == nil || len(p.TableDetailz) == 0 {
	   println("DB: missing app table details. Aborting.")
	   return nil, errors.New("Missing app DB table details")
	}
	e = pSR.RegisterAppTables("", p.TableDetailz) // DRM.M5_TableDetails)
	if e != nil {
		return nil, fmt.Errorf("%s can't register tables: %w", errPfx, e)
	}
	
	// If the DB exists and we want to open
	// it as-is, i.e. without zeroing it out,
	// then this is where we return success. 
	if filexist && repo != nil && !p.DoZeroOut {
	   	// Some weirdness ? A non-fatal error ? 
	   	if openError != nil {
		   L.L.Warning("Expect DB problems: " +
		   	"DB init got error: %w", openError)
		}
		L.L.Info("DB file opened OK: " + shortPath)
	   	repo.SetLogWriter(dbw)
		return repo, nil
	}
	if !filexist {
	   	println("DB: Creating anew.")
		L.L.Info("Creating DB: " + shortPath)
		if p.DoZeroOut {
			L.L.Info("Zeroing out the DB is redundant")
		}
		repo, e = DB_Manager.NewAtPath(dbFilepath)
		// Now configure it using pragmas
		ps, pe := repo.DoPragmas(DB_Manager.InitznPragmas())
		if pe != nil {
		   L.L.Warning("Expect DB problems: " +
                        "DB init pragmas got error: %w", pe)
			}
		L.L.Info("Ran init pragmas on DB: " + ps)
	}
	if e != nil {
		return nil, fmt.Errorf("%s DB failure: %w", errPfx, e)
	}
	repo.SetLogWriter(dbw)
	repoAbsPath := repo.Path()
	println("DB: status OK.")
	L.L.Info("DB OK: " + SU.ElideHomeDir(repoAbsPath))

	if !filexist {
		// env.SimpleRepo.ForceExistDBandTables()
		e = pSR.CreateAppTables()

	} else if p.DoZeroOut {
		L.L.Info("Zeroing out DB (init9nb.go)")
		_, e := repo.CopyToBackup()
		if e != nil {
			panic(e)
		}
		pSR.EmptyAppTables()
	}
	return repo, nil
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

