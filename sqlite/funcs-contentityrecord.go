package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	// DRU "github.com/fbaube/datarepo/utils"
	DRM "github.com/fbaube/datarepo/rowmodels"
	SU "github.com/fbaube/stringutils"
	CA "github.com/fbaube/contentanalysis"
)

// NewContentityRow does content fetching &
// analysis, while "promoting" a [PathProps];
// it work for directories and symlinks too.
// .
func NewContentityRow(pPP *FU.PathProps, pPA *CA.PathAnalysis) (*DRM.ContentityRow, error) {
	if pPP == nil || pPA == nil {
		panic("OOPS")
	}
	if pPA.MarkupType() == "UNK" {
		// panic("UNK MarkupType in ExecuteStages")
		return nil, fmt.Errorf(
			"reposqlite.fxCntyrow.NewCR: got MarkupType UNK")
	}
	// var e error
	pNewCR := new(DRM.ContentityRow)
	pNewCR.PathProps = *pPP
	pNewCR.PathAnalysis = pPA

	if !pPP.Exists() {
		L.L.Error(pPP.String())
		return nil, errors.New(
			"input PathProps does not exist: " + pPP.AbsFP.S())
	}
	if pPP.IsDir() || pPP.IsSymlink() {
		// COMMENTING THIS OUT IS A FIX
		// pCR.SetError(errors.New("Is directory or symlink"))
		return pNewCR, nil
	}
	if !pPP.IsFile() {
		return pNewCR, errors.New("is not valid file")
	}
	// =======================
	//  More content analysis
	// =======================
	if pPA.MType == "" {
		L.L.Warning("No MType, so trying snift-MIME-type: %s", pPA.MimeTypeAsSnift)
		switch pPA.MimeTypeAsSnift {
		case "text/xml/image/svg+xml":
			println("SVG!!")
			pPA.MType = "xml/cnt/svg"
		}
	}
	pNewCR.PathAnalysis = pPA
	if pNewCR.MarkupType() == "UNK" {
		panic("UNK MarkupType in NewContentityRow")
	}
	// SPLIT FILE!
	if !pPA.ContentityBasics.HasNone() {
		L.L.Okay("Key elms: XmlRoot<%s> Meta<%s> Text<%s>",
			pPA.ContentityBasics.XmlRoot.Info(),
			pPA.ContentityBasics.Meta.Info(),
			pPA.ContentityBasics.Text.Info())
	} else if pPA.MarkupType() == SU.MU_type_MKDN {
		// pPA.KeyElms.SetToAllText()
		// L.L.Warning("TODO set MKDN all text, and ranges")
	} else if pPA.MarkupType() == SU.MU_type_BIN {
	} else {
		L.L.Warning("Found no key elms (root,meta,text)")
	}
	// fmt.Printf("D=> NewCR: %s \n", pCR.String())
	return pNewCR, nil
}

// GetContentityAll gets all content in the DB.
func (p SqliteRepo) GetContentityAll() (pp []*DRM.ContentityRow, err error) {
	var rows *sql.Rows
	pp = make([]*DRM.ContentityRow, 0, 16)
	q := "SELECT * FROM CONTENT"
	// rows, err := p.Handle().Queryx(q)
	rows, err = p.Handle().Query(q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		panic("GetContentityAll")
	}
	for rows.Next() {
		p := new(DRM.ContentityRow)
		// err := rows.StructScan(p)
		if err = rows.Scan(p.PtrFields()...); err != nil {
			return nil, fmt.Errorf("GetContentityAll: "+
				"row.Scan error: %w \n\t (%s)", err, q)
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetContentityAll: "+
			"rows.Err(): %w \n\t (%s)", err, q)
	}
	return pp, nil
}

// InsertContentityRow adds a content item (i.e. a file) to the DB.
func (p SqliteRepo) InsertContentityRow(pC *DRM.ContentityRow) (int, error) {
	var rslt sql.Result
	var stmt string
	var e error
	// println("REL:", pC.RelFP)
	// println("ABS:", pC.AbsFP)

	pC.T_Cre = SU.Now() // time.Now().UTC().Format(time.RFC3339)
	pC.T_Imp = SU.Now() // time.Now().UTC().Format(time.RFC3339)
	// tx := p.Handle().MustBegin()
	tx, err := p.Handle().Begin()
	if err != nil {
		panic(err)
	}
	stmt = "INSERT INTO CONTENTITY(" +
		"idx_inbatch, descr, relfp, absfp, " +
		"t_cre, t_imp, t_edt, " +
		"mimetype, mtype, " +
		"xmlcontype, ditaflavor, ditacontype" +
		") VALUES(" +

		":idx_inbatch, :descr, :relfp, :absfp, " +
		":t_cre, :t_imp, :t_edt, " +
		":mimetype, :mtype, " +
		":xmlcontype, :ditaflavor, :ditacontype);"

	// rslt, e = tx.NamedExec(stmt, pC)
	fmt.Printf("funcs_contentity.L141: " +
		"skipping NamedExec(INSERT INTO CONTENTITY(values)) \n")
	tx.Commit()
	// println("=== ### ===")
	if e != nil {
		L.L.Error("DB.Add_Contentity: %s", e.Error())
	}
	if e != nil {
		println("========")
		println("DB: NamedExec: ERROR:", e.Error())
		println("========")
		fnam := "./insert-Contentity-failed.sql"
		e = ioutil.WriteFile(fnam, []byte(stmt), 0644)
		if e != nil {
			L.L.Error("Could not write file: " + fnam)
		} else {
			L.L.Dbg("Wrote \"INSERT INTO contentity ... \" to: " + fnam)
		}
		// panic("INSERT CONTENTITY failed")
		return -1, e
	}
	liid, _ := rslt.LastInsertId()
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertFile: ID %d nR %d \n", liid, naff)
	return int(liid), nil
}
