package sqlite

import (
	"fmt"
	D "github.com/fbaube/dsmnd"
	DRM "github.com/fbaube/datarepo/rowmodels"
	S "strings"
	// "time"
)

/* REF
https://www.sqlite.org/lang_createtable.html
Create [ Temp | Temporary ] Table
    [ If Not Exists ] [schemaname.]tablename
1) as select-statement
2) ( columndef,+ tableconstraint,* )
[ table-options:Strict ] ;
*/

/*
A digression on FKEY's ...

Here is the typical use case for a foreign key:

idx_inbatch integer not null references inbatch, -- int64 in Go
foreign key(idx_inbatch) references inbatch(idx_inbatch)

However... a table can have multiple foreign key references into
another table. Here is the use case:

idx_cnt_map integer not null references contentity,
idx_cnt_tpc integer not null references contentity,
foreign key(idx_cnt_map) references contentity(idx_contentity),
foreign key(idx_cnt_tpc) references contentity(idx_contentity),

We can detect this situation, and correct for it,
only by counting underscores in the StorName.
*/

// NewCreateTableStmt expects [TableDetails.IDName] to be set. 
// . 
func (pSR *SqliteRepo) NewCreateTableStmt(pTD *DRM.TableDetails) (string, error) {
	{
		var sbDbg S.Builder
		sbDbg.WriteString(fmt.Sprintf(
			"=== %s.GenCreTblStmt.ColSpex ===\n", pTD.StorName))
		for _, pCS := range pTD.ColumnSpecs {
		    cnm := pCS.StorName // column name
		    bdt := pCS.Datatype
		    sbDbg.WriteString(fmt.Sprintf("%s:%s, ", cnm, bdt))
		}
		fmt.Fprintf(pSR.w, sbDbg.String() + "\n")
	}
        // sqlStmt is for building the actual SQL statement
	var sqlStmt S.Builder
	sqlStmt.WriteString(fmt.Sprintf("CREATE TABLE %s(\n", pTD.StorName))

	// ======================================
	//  PRIMARY KEY IS ASSUMED - DO IT FIRST
	//     idx_mytable integer not null
	//     primary key autoincrement,
	// ======================================
	if pTD.PKname == "" {
	   panic("pTD(" + pTD.StorName + ").PKname is not set")
	}
	// stmt is too long, so use a newline
	sqlStmt.WriteString(pTD.PKname + 
		" integer not null \n    primary key autoincrement, " +
		"-- NOTE: integer, not int. \n")
	// =====================================
	// NOW ALL THE FIELDS EXPLICITLY DEFINED
	// =====================================
	for _, pCS := range pTD.ColumnSpecs {
		colName := pCS.StorName // column name in DB
		// println("Creating column:", pCS.String())
		
		// Get the BDT (which determines the 
		// storage keyword for SQLite) and the 
		// SFT (for resolving some of the cases)
		SFT := D.SemanticFieldType(pCS.Datatype)
		BDT := SFT.BasicDatatype()

		switch BDT {
		/* Unimplemented:
		BDT_NIL  =  BasicDatatype("nil") 
		BDT_FLOT // SQLITE_FLOAT   2
		BDT_BLOB // SQLITE_BLOB    4
		BDT_NULL // SQLITE_NULL    5
		BDT_LIST // List (simple one-dimensional lists) 
		BDT_CLXN // Collection (more-complicated data strux)
		BDT_OTHR // reserved: expansion 
		BDT_NONE // reserved 
		*/
		
		// Text fields are simple
		case D.BDT_TEXT: // maps to SQLITE_TEXT 3
			sqlStmt.WriteString(colName + " text not null,\n")
			
		// Int(eger) fields are simple (they are not keys) 
		case D.BDT_INTG: // maps to SQLITE_INTEGER 1 
			// TODO: int not null check (filect >= 0) default 0
			sqlStmt.WriteString(colName + " int not null,\n")
			
		// Date/time fields are assumed to be ISO-8601,
		// using SQLITE TEXT 
		case D.BDT_DYTM: // maps to SQLYT_DATETIME 6
		  // println("column DATE-TIME:",
		  // 	D.Datum(SFT.Descriptor()).String())
		  // Use ISO-8601 strings ("YYYY-MM-DD HH:MM:SS. SSS")
		     sqlStmt.WriteString(colName + " text not null, -- ISO-8601 \n")

		// Keys get real complicated real fast 
		case D.BDT_KEYY: // PRIMARY/FOREIGN/OTHER KEY (SQLite "INTEGER")
		switch SFT { 
		  case D.SFT_PRKEY: 
			panic("DUPE PRIMARY KEY: " +
				"explicit decl'n is dupe of implicit")
		  case D.SFT_FRKEY: 
			//> e.g. D.ColumnSpec{D.FKEY, "idx_inbatch", 
			//>     "inbatch", "Input batch of imported content"},
			// The referencing fields's name is (e.g.) idx_inbatch
			refgField := colName
			// The referenced table's name is (e.g.) inbatch
			refdTable := pCS.DispName

			// Count up underscores (per comment above)
			ss := S.Split(refgField, "_")
			// fmt.Printf("FKey name split: %#v \n", ss)
			
			switch len(ss) {

			case 2: // normal case, for example:
				//> idx_inbatch integer not null ref's inbatch,
				//> foreign key(idx_inbatch) references
				//>     inbatch(idx_inbatch)
				// BETTER:
				// in table "contentity":
				// idx_contentity integer not null
				//     primary key autoincrement,
				// idx_inbatch integer not null
				//     references inbatch(idx_inbatch),
				sqlStmt.WriteString(fmt.Sprintf(
					"%s integer not null references %s(%s),\n",
					refgField, refdTable, refgField))
			case 3: // multiple indices into same table, e.g.
				//> idx_cnt_map integer not null ref's cont'y,
				//> frn key(idx_cnt_map) refs cont'y(idx_cont'y),
				//> idx_cnt_tpc integer not null ref's cont'y,
				//> frn key(idx_cnt_tpc) refs cont'y(idx_cont'y),
				// We have to deduce "idx_contentity", which
				// we can do confidently after passing checks.
				var refdField string
				if S.EqualFold(ss[0], "idx") &&
					S.EqualFold(ss[1][0:1], refdTable[0:1]) {
					refdField = "idx_" + refdTable
					sqlStmt.WriteString(fmt.Sprintf(
					  "%s integer not null references %s(%s),\n",
					   refgField, refdTable, refdField))
				} else {
					return "", fmt.Errorf(
						"Malformed a_b_c FKEY: %s,%s,%s",
						refgField, refdTable, refdField)
				}
			default:
				return "", fmt.Errorf("Malformed FKEY: "+
					"%s,%s,%s", refgField, refdTable)
			}
		  }
		default:
			fmt.Fprintf(pSR.w, "pTD.StorName CreTbl OOPS: " +
				"col<%s> bdt<%s> sft<%s> \n", colName, BDT, SFT) 
			panic(pCS.Datatype)
		}
	}
	// trim off final ",\n"
	ss := sqlStmt.String()
	// and add STRICT 
	stmt3 := ss[0:len(ss)-2] + "\n) STRICT;"
	fmt.Fprintf(pSR.w, "=== %s.GenCreTblStmt.SQL ===\n%s \n",
		pTD.StorName, stmt3) 
	return stmt3, nil
}

/*
CREATE TABLE contentity( -- STRICT!
idx_inbatch integer not null references inbatch,
relfp text not null check (typeof(relfp) == 'text'),
absfp text not null check (typeof(absfp) == 'text'),
t_cre text not null check (typeof(t_cre) == 'text'),
t_imp text not null check (typeof(t_imp) == 'text'),
t_edt text not null check (typeof(t_edt) == 'text'),
descr text not null check (typeof(descr) == 'text'),
mimetype text not null check (typeof(mimetype) == 'text'),
mtype     text not null check (typeof(mtype)     == 'text'),
xmlcontype text not null check (typeof(xmlcontype) == 'text'),
ditaflavor text not null check (typeof(ditaflavor) == 'text'),
ditacontype text not null check (typeof(ditacontype) == 'text'),
foreign key(idx_inbatch) references inbatch(idx_inbatch)
);
*/
