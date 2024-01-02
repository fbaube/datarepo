package sqlite

import (
/*
		"database/sql"
		"errors"
		"fmt"
		"io/ioutil"
		"log"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	DRU "github.com/fbaube/datarepo/utils"
	SU "github.com/fbaube/stringutils"
*/
)

/*

// NewSqlarRecord does content fetching &
// analysis, while "promoting" a [PathProps];
// it works for directories and symlinks too.
// .
// func NewContentityRow(pPP *FU.PathProps) (*RU.ContentityRow, error) {
func NewSqlarRecord(pPP *FU.PathProps, pPA *FU.PathAnalysis) (*RU.ContentityRow, error) {
	if pPP == nil || pPA == nil {
		panic("OOPS")
	}
	var e error
	pNewCR := new(RU.ContentityRow)
	pNewCR.PathAnalysis = pPA

	if !pPP.Exists() {
		L.L.Error(pPP.String())
		return nil, errors.New(
			"input PathProps does not exist: " + pPP.AbsFP.S())
	}
	/*
		if pPP.Size() == 0 {
			// panic("barf")
			return pCR, nil
		}
	* /
	if pPP.IsDir() || pPP.IsSymlink() {
		// COMMENTING THIS OUT IS A FIX
		// pCR.SetError(errors.New("Is directory or symlink"))
		return pNewCR, nil
	}
	if !pPP.IsFile() {
		return pNewCR, errors.New("is not valid file")
	}
	// OK, it's a valid file. But maybe it's empty!
	// If it's not, copy it into a convenient local variable.
	if pPP.Size() == 0 {
		L.L.Progress("Skipping fetch for zero-length content")
	} else {
		// e = pPP.FetchRaw()
		e = pPP.GoGetFileContents()
		if e != nil {
			L.L.Error("newCnty: cannot fetch raw content: %w", e)
			return pNewCR, fmt.Errorf(
				"newCnty: cannot fetch raw content: %w", e)
		}
	}
	if len(pPP.Raw) == 0 {
		L.L.Error("OOPS raw 0, <%s>, funcs_contentityrecord.go L68",
			pPP.AbsFP.S())
		return nil, fmt.Errorf("No raw: skipping: %s", pPP.AbsFP.S())
	}
	// ==============================
	//  Perform content analysis.
	//  pAR will be copied into pCR.
	// ==============================
	var pAR *FU.PathAnalysis
	// pAR, e = FU.DoAnalysis(sCont, FP.Ext(string(pPP.AbsFP)))
	pAR, e = FU.NewPathAnalysis(pPP)
	if e != nil {
		L.L.Error("newCtyRec: doAnalysis failed: " + e.Error())
		return pNewCR, fmt.Errorf("newCtyRec: doAnalysis failed: %w", e)
	}
	// L.L.Warning("==== LEN %d ====", len(pAR.PathProps.Raw))
	// fmt.Printf("==== LEN %d ==== \n", len(pAR.ContentityBasics.Raw))
	L.L.Okay("DoAnalysis: [%.20s]%s", pPP.Raw,
		pAR.ContypingInfo.MultilineString())
	if pAR.MType == "" {
		L.L.Warning("No MType, so trying snift-MIME-type: %s", pAR.MimeTypeAsSnift)
		switch pAR.MimeTypeAsSnift {
		case "text/xml/image/svg+xml":
			println("SVG!!")
			pAR.MType = "xml/cnt/svg"
		}
	}
	if pAR == nil {
		panic("NIL pAR")
	}
	pNewCR.PathAnalysis = pAR
	// SPLIT FILE!
	if !pAR.ContentityBasics.HasNone() {
		L.L.Okay("Key elm triplet: XmlRoot<%s> Meta<%s> Text<%s>",
			pAR.ContentityBasics.XmlRoot.Info(),
			pAR.ContentityBasics.Meta.Info(),
			pAR.ContentityBasics.Text.Info())
	} else if pAR.MarkupType() == SU.MU_type_MKDN {
		// pAR.KeyElms.SetToAllText()
		// L.L.Warning("TODO set MKDN all text, and ranges")
	} else if pAR.MarkupType() == SU.MU_type_BIN {
	} else {
		L.L.Warning("Found no key elms (root,meta,text)")
	}
	// fmt.Printf("D=> NewCR: %s \n", pCR.String())
	return pNewCR, nil
}

// GetContentityAll gets all content in the DB.
func (p SqliteRepo) GetSqlarAll() (pp []*RU.ContentityRow) {
	pp = make([]*RU.ContentityRow, 0, 16)
	rows, err := p.Handle().Queryx("SELECT * FROM CONTENT")
	if err != nil {
		panic("GetContentityAll")
	}
	for rows.Next() {
		p := new(RU.ContentityRow)
		err := rows.StructScan(p)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("    DD:%#v\n", *p)
		pp = append(pp, p)
	}
	return pp
}

// InsertContentityRow adds a content item (i.e. a file) to the DB.
func (p SqliteRepo) InsertSqlarRecord(pC *RU.ContentityRow) (int, error) {
	var rslt sql.Result
	var stmt string
	var e error
	// println("REL:", pC.RelFP)
	// println("ABS:", pC.AbsFP)

	pC.T_Cre = SU.Now() // time.Now().UTC().Format(time.RFC3339)
	pC.T_Imp = SU.Now() // time.Now().UTC().Format(time.RFC3339)
	tx := p.Handle().MustBegin()
	stmt = "INSERT INTO CONTENTITY(" +
		"idx_inbatch, descr, relfp, absfp, " +
		"t_cre, t_imp, t_edt, " +
		// "metaraw, textraw, " +
		"mimetype, mtype, " +
		// roottag, rootatts, " +
		// "xmlcontype, xmldoctype, ditaflavor, ditacontype" +
		"xmlcontype, ditaflavor, ditacontype" +
		") VALUES(" +

		// ":idx_inbatch, :pathprops.relfp, :pathprops.absfp, " +
		":idx_inbatch, :descr, :relfp, :absfp, " +

		// ":times.t_cre, :times.t_imp, :times.t_edt, " +
		":t_cre, :t_imp, :t_edt, " +

		// ":metaraw, :textraw, " +
		// ":mimetype, :mtype, " +
		":mimetype, :mtype, " +
		// ":root.name, :root.atts, " +
		// ":analysisrecord.contentitystructure.root.name, " +
		// ":analysisrecord.contentitystructure.root.atts, " +

		// ":xmlcontype, :xmldoctype, :ditaflavor, :ditacontype);"
		":xmlcontype, :ditaflavor, :ditacontype);"
		// ":doctype, :ditaflavor, :ditacontype);"

	rslt, e = tx.NamedExec(stmt, pC)
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

*/
