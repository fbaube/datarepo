package sqlite

import (
	"errors"
	"fmt"
	"os"
	// "is"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	FP "path/filepath"
)

// MoveToBackup makes a best effort but might fail if
// (for a file) the backup destination is a directory
// or has a permissions problem. The current Repo is
// renamed, so it seems to "disappear" from production./
//
// FIXME: Use best practices to do this: make sure it
// is robust; this will need a check for the DB type.
// FIXME: Before moving, close it if it is open.
func (p *SqliteRepo) MoveToBackup() (string, error) {
	var fromFP string
	if fromFP = p.Path(); fromFP == "" {
		return "", errors.New("sqliterepo.movetobackup: no from-path")
	}
	// if p.Type() != RU.DB_SQLite {
	//	return "", errors.New("sqliterepo.movetobackup: not sqlite")
	// }
	cns := SU.NowAsYMDHM()
	toFP := FU.AppendToFileBaseName(fromFP, "-"+cns)
	toFP, _ = FP.Abs(toFP)
	// func os.Rename(oldpath, newpath string) error
	e := os.Rename(fromFP, toFP)

	if e != nil {
		return "", fmt.Errorf("sqliterepo.movetobackup: "+
			"can't move current DB to <%s>: %w", toFP, e)
	}
	return toFP, nil
}

// CopyToBackup makes a best effort but might fail if
// (for a file) the backup destination is a directory
// or has a permissions problem. The current DB is
// never affected.
//
// FIXME: Use best practices to do this: make sure it
// is robust; this will need a check for the DB type.
// FIXME: If it is open, does it need to be closed -
// and then reopened after the copying ?
func (p *SqliteRepo) CopyToBackup() (string, error) {
	var fromFP string
	if fromFP = p.Path(); fromFP == "" {
		return "", errors.New("sqliterepo.copytobackup: no from-path")
	}
	// if p.Type() != RU.DB_SQLite {
	//	return "", errors.New("sqliterepo.copytobackup: not sqlite")
	// }
	cns := SU.NowAsYMDHM()
	toFP := FU.AppendToFileBaseName(fromFP, "-"+cns)
	toFP, _ = FP.Abs(toFP)
	e := FU.CopyFromTo(fromFP, toFP)
	if e != nil {
		return "", fmt.Errorf("sqliterepo.copytobackup: "+
			"can't copy current DB to <%s>: %w", toFP, e)
	}
	return toFP, nil
}

func (p *SqliteRepo) RestoreFromMostRecentBackup() (string, error) {
	return "", errors.New("RestoreFromMostRecentBackup: not implemented")
}
