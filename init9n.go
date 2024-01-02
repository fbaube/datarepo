package datarepo

import (
	"errors"
	"fmt"
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	DRS "github.com/fbaube/datarepo/sqlite"
	RM "github.com/fbaube/rowmodels"
	R "github.com/fbaube/repo"
	// _ "github.com/fbaube/sqlite3"
	_ "github.com/mattn/go-sqlite3"
	SU "github.com/fbaube/stringutils"
	"os"
	_ "database/sql"
	DRU "github.com/fbaube/datarepo/utils"
)

type Init9nArgs struct {
     // D.DB_type is so far only D.DB_SQLite = "sqlite"
     D.DB_type
     // BaseFilename defaults to "mmmc.db"
     BaseFilename string 
     Dir string
     // DoImport is a flag that access to a DB is required,
     // and could also be set for other such operations 
     DoImport bool
     // DoZeroOut says initialize the DB with the 
     // app's tables but with no data in them 
     DoZeroOut bool
     // TableDetails are app table details 
     TableDetails []DRU.TableDescriptor
}

var DEFAULT_FILENAME = "mmmc.db"

// ProcessInit9nArgs processes DN initialization arguments. 
// It can process either a new DB OR an existing DB.
// 
// TODO: It should not use a logger (i.e. [mlog]),
// because we want to avoid that kind of dependency
// in a standalone library. 
// .
func (p *Init9nArgs) ProcessInit9nArgs() (SimpleRepo, error) {

        var repo SimpleRepo
	// type-checking
	// var _ repo.SimpleRepo = (*RS.SqliteRepo)(nil)

	var mustAccessTheDB bool
	var e error
	mustAccessTheDB = p.DoImport || p.DoZeroOut || p.Dir != ""
	if !mustAccessTheDB {
		return nil, nil 
	}
	if p.DB_type != D.DB_SQLite {
	   return nil, errors.New("bad DB type: " + string(p.DB_type))
	}

	// Start by checking on the status of the filename.
	// This all assumes that the DB is SQLite, a single file.
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
	L.L.Info("DB resolved path: " + dbFilepath)
	errPfx := fmt.Errorf("processDBargs(%s):", dbFilepath)
	// func IsFileAtPath(aPath string) (bool, *os.FileInfo, error) {

	var fileinfo os.FileInfo
	filexist, fileinfo, filerror := FU.IsFileAtPath(dbFilepath)
	if filerror != nil {
		panic("L71")
		return nil, fmt.Errorf("%s file error: %w", errPfx, filerror)
	}
	s := SU.ElideHomeDir(dbFilepath)
	if filexist {
		L.L.Info("DB exists: " + s)
		if fileinfo.Size() == 0 {
			L.L.Info("DB is empty: " + s)
			e = os.Remove(dbFilepath)
			if e != nil {
				panic(e)
			}
			filexist = false
		} else {
		        repo, e = DRS.OpenRepoAtPath(dbFilepath)
			// If the DB exists and we want to open
			// it as-is, i.e. without zeroing it out,
			// then this is where we return success:
			if e == nil && !p.DoZeroOut {
				L.L.Info("DB opened: " + s)
				return repo, nil
			}
		}
	}
	if !filexist {
		L.L.Info("Creating DB: " + s)
		if p.DoZeroOut {
			L.L.Info("Zeroing out the DB is redundant")
		}
		repo, e = DRS.NewRepoAtPath(dbFilepath)
	}
	if e != nil {
		return nil, fmt.Errorf("%s DB failure: %w", errPfx, e)
	}
	repoAbsPath := repo.Path()
	L.L.Info("DB OK: " + SU.ElideHomeDir(repoAbsPath))

	pSQR, ok := repo.(*DRS.SqliteRepo)
	if !ok {
		panic("L100")
		return nil, errors.New("processDBargs: is not sqlite")
	}
	e = pSQR.SetAppTables("", RM.MmmcTableDescriptors)
	/* type RepoAppTables interface {
		// SetAppTables specifies schemata
		SetAppTables(string, []U.TableConfig) error
		// EmptyAllTables deletes (app-level) data
		EmptyAppTables() error
		// CreateTables creates/empties the app's tables
		CreateAppTables() error
	} */
	if !filexist {
		// env.SimpleRepo.ForceExistDBandTables()
		e = pSQR.CreateAppTables()

	} else if p.DoZeroOut {
		L.L.Progress("Zeroing out DB")
		_, e := repo.CopyToBackup()
		if e != nil {
			panic(e)
		}
		pSQR.EmptyAppTables()
	}
	return repo, nil
}

/*
// inputExts more than covers the file types associated with the LwDITA spec.
// Of course, when we check for them we do so case-insensitively.
var inputExts = []string{
	".dita", ".map", ".ditamap", ".xml",
	".md", ".markdown", ".mdown", ".mkdn",
	".html", ".htm", ".xhtml", ".png", ".gif", ".jpg"}

// AllGLinks gathers all GLinks in the current run's input set.
var AllGLinks mcfile.GLinks
*/

