package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	L "github.com/fbaube/mlog"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// GetAll_Inbatch gets all input batches in the system.
func (p SqliteRepo) GetAll_Inbatch() (pp []*DRM.InbatchRow, err error) {
	var rowsx *sql.Rows
	var e error
	q := "SELECT * FROM INBATCH"
	// rowsx, e = p.Handle().Queryx(q)
	rowsx, err = p.Handle().Query(q)
	if e != nil {
		L.L.Error("DB.GetAll_Inbatch: %w", e)
		return nil, fmt.Errorf("funcs_inbatch.L19")
	}
	pp = make([]*RM.InbatchRow, 0, 16)
	for rowsx.Next() {
		p := new(RM.InbatchRow)
		// ====================
		// e = rowsx.StructScan(p)
		panic("funcs_inbatch L27")
		// ====================
		if e != nil {
			L.L.Error("DB.GetAll_Inbatch.StructScan: %w", e)
		}
		L.L.Dbg("Got Inbatch: %+v", *p)
		pp = append(pp, p)
	}
	return pp, nil
}

// Add_Inbatch adds an input batch to the DB and returns its primary index.
func (p SqliteRepo) Add_Inbatch(pIB *DRM.InbatchRow) (int, error) {
	var rslt sql.Result
	var stmt string
	var e error

	if pIB.FilCt == 0 {
		pIB.FilCt = 1
	} // HACK

	pIB.T_Cre = time.Now().UTC().Format(time.RFC3339)
	// tx := p.Handle().MustBegin()
	tx, err := p.Handle().Begin()
	if err != nil {
		panic(err)
	}
	stmt = "INSERT INTO INBATCH(" +
		"descr, filct, t_cre, relfp, absfp" +
		") VALUES(" +
		":descr, :filct, :t_cre, :relfp, :absfp)" // " RETURNING i_INB", p)
	// rslt, e = tx.NamedExec(stmt, pIB)
	fmt.Printf("funcs_inbatch.L59: "+
		"skipping NamedExec(INSERT INTO INBATCH(values)) <%s>\n", stmt)
	tx.Commit()
	// println("=== ### ===")
	if e != nil {
		L.L.Error("DB.Add_Inbatch: %w", e)
	}
	/* SQL stuff
		Query(...) (*sql.Rows, error) - unchanged
		QueryRow(...) *sql.Row - unchanged
		Extensions:
		Queryx(...) (*sqlx.Rows, error) - Query, but return an sqlx.Rows
		QueryRowx(...) *sqlx.Row -- QueryRow, but return an sqlx.Row
		New semantics:
		Get(dest interface{}, ...) error // to fetch one scannable
		Select(dest interface{}, ...) error // to fetch multi scannables
		Scannable means: simple datum not struct OR struct w no exported fields OR
	implements sql.Scanner f
		"SELECT * FROM INBATCH"
	*/

	// func StructScan(rows rowsi, dest interface{}) error
	// StructScan all rows from an sql.Rows or an sqlx.Rows into the dest slice.
	// StructScan will scan in the entire rows result; to get fewer, use Queryx
	// and see sqlx.Rows.StructScan. If rows is sqlx.Rows, it will use its mapper,
	// otherwise it will use the default.
	// ============

	/* potential code
	var egInb = Inbatch{}
	var rowsx *sqlx.Rows
	rowsx, e = p.DB.Queryx("SELECT * FROM INBATCH")
	TestInbatch(rowsx, &egInb)
	*/

	// WORK HERE

	L.L.Warning("TODO: RETURNING (inbatch ID)")
	liid, err := rslt.LastInsertId()
	if err != nil {
		panic(err)
	}
	// naff, _ := rslt.RowsAffected()
	// fmt.Printf("    DD:InsertInbatch: ID=%d (nR=%d) \n", liid, naff)
	return int(liid), nil
}
