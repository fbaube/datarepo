package utils

import (
	_ "database/sql" // to get init()
	// D "github.com/fbaube/dsmnd"
	// L "github.com/fbaube/mlog"
	// S "strings"
)

/*

func (tc *TableDescriptor) StmtCreateAppTable() string {

	var CTS string // the Create Table SQL string
	var hasFKs bool
	hasFKs = (tc.ForenKeys != nil && len(tc.ForenKeys) > 0)

	// === CREATE TABLE
	CTS = "CREATE TABLE " + tc.Name + "(\n"
	// == PRIMARY KEY
	CTS += "idx_" + tc.Name + " integer not null primary key autoincrement, "
	CTS += "-- NOTE: integer, not int. \n"
	if hasFKs {
		// === FOREIGN KEYS
		// []string{"map_contentity", "tpc_contentity"},
		for _, tbl := range tc.ForenKeys {
			if S.Contains(tbl, "_") {
				i := S.LastIndex(tbl, "_")
				minTbl := tbl[i+1:]
				L.L.Info("DB compound index: " + tbl + " indexes " + minTbl)
				CTS += "idx_" + tbl + " integer not null references " + minTbl + ", \n"
			} else {
				// idx_inb integer not null references INB,
				// "not null" might be problematic during development.
				CTS += "idx_" + tbl + " integer not null references " + tbl + ", \n"
			}
		}
	}
	for _, fld := range tc.Columns {
		switch fld.Fundatype {
		case D.INTG:
			// e.g.: filect int not null check (filect >= 0) default 0
			// also: `Col1 INTEGER CHECK (typeof(Col1) == 'integer')`
			//
			CTS += fld.StorName + " int not null"
			// CTS += fld.Code + " int not null check (typeof(" + fld.Code + ") == 'int')"
			/* add checks
			switch tc.intRanges[i] {
			case 1:
				// check (filect >= 0)
				CTS += " check (" + fld + " > 0), \n"
			case 0:
				CTS += " check (" + fld + " >= 0), \n"
			default: // case -1:
				CTS += ", \n"
			}
			* /
			CTS += ", \n"
		case D.TEXT:
			CTS += fld.StorName + " text not null check " +
				"(typeof(" + fld.StorName + ") == 'text'), \n"
		default:
			panic("Unhandled: " + fld.Fundatype)
		}
	}
	if hasFKs {
		// FOREIGN KEY(idx_inb) REFERENCES INB(idx_inb)
		for _, tbl := range tc.ForenKeys {
			// idx_inb integer not null references INB,
			// TMP := "foreign key(idx_" + tbl + ") references " + tbl + "(idx_" + tbl + "), \n"
			// println("TMP:", TMP)
			CTS += "foreign key(idx_" + tbl + ") references " + tbl + "(idx_" + tbl + "), \n"
		}
	}

	CTS = S.TrimSuffix(CTS, "\n")
	CTS = S.TrimSuffix(CTS, " ")
	CTS = S.TrimSuffix(CTS, ",")
	CTS += "\n);"

	return CTS

	/* write file w SQL to create table
	fnam := "./create-table-" + tc.Name + ".sql"
	e := ioutil.WriteFile(fnam, []byte(CTS), 0644)
	if e != nil {
		L.L.Error("Could not write file: " + fnam)
	} else {
		L.L.Dbg("Wrote \"CREATE TABLE " + tc.Name + " ... \" to: " + fnam)
	}
	* /
}
*/
