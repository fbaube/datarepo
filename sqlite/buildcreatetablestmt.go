package sqlite

import (
	"fmt"
	D "github.com/fbaube/dsmnd"
	DRU "github.com/fbaube/datarepo/utils"
	S "strings"
	// "time"
)

/* REF
https://www.sqlite.org/lang_createtable.html
Create [ Temp | Temporary ] Table [ If Not Exists ] [schemaname.]tablename
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

// BuildCreateTableStmt expects [TableDetails.IDName] to be set. 
// . 
func (pSR *SqliteRepo) BuildCreateTableStmt(pTD *DRU.TableDetails) (string, error) {
     	// sb is for building the actual SQL statement.
	// sbDng is for a debugging message. 
	var sb, sbDbg S.Builder

	sb.WriteString(fmt.Sprintf("CREATE TABLE %s(\n", pTD.StorName))
	
	sbDbg.WriteString("reposqlite.GenCreTblStmt: ")
	for _, pCS := range pTD.ColumnSpecs {
		cnm := pCS.StorName // column name
		bdt := pCS.Datatype
		sbDbg.WriteString(fmt.Sprintf("%s:%s, ", cnm, bdt))
	}
	fmt.Printf(sbDbg.String() + "\n")

	// PRIMARY KEY IS ASSUMED - DO IT FIRST
	// idx_mytable integer not null primary key autoincrement,
	if pTD.IDName == "" {
	   panic("pTD.IDName is not set")
	}
	
	sb.WriteString(pTD.IDName +
		" integer not null primary key autoincrement, " +
		"-- NOTE: integer, not int. \n")

	for _, pCS := range pTD.ColumnSpecs {
		colName := pCS.StorName // column name in DB
		println("Creating column:", pCS.String())
		SFT := D.SemanticFieldType(pCS.Datatype)
		BDT := SFT.BasicDatatype()

		switch BDT {
		case D.BDT_TEXT: // maps to SQLITE_TEXT 3
			sb.WriteString(colName + " text not null,\n")
		case D.BDT_INTG: // maps to SQLITE_INTEGER filect
			// TODO: int not null check (filect >= 0) default 0
			sb.WriteString(colName + " int not null,\n")
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
		case D.BDT_DYTM: // maps to SQLYT_DATETIME 6
		     println("column DATE-TIME:",
		     	D.Datum(SFT.Descriptor()).String())
		// Use ISO-8601 strings ("YYYY-MM-DD HH:MM:SS. SSS")
		     sb.WriteString(colName + " text not null, -- ISO-8601 \n")

		case D.BDT_KEYY: // PRIMARY/FOREIGN/OTHER KEY (SQLite "INTEGER")
		switch SFT { 
		  case D.SFT_PRKEY: 
			panic("DUPE PRIMARY KEY: " +
				"explicit decl'm conflix w implicit")
		  case D.SFT_FRKEY: 
			//> D.ColumnSpec{D.FKEY, "idx_inbatch", "inbatch",
			//>  "Input batch of imported content"},
			// referencing fields's name is idx_inbatch
			refgField := colName
			// referenced table's name is inbatch
			refdTable := pCS.DispName

			// Count up underscores (per comment above)
			ss := S.Split(refgField, "_")
			fmt.Printf("FKey name split: %#v \n", ss)
			
			switch len(ss) {

			case 2: // normal case, for example:
				//> idx_inbatch integer not null ref's inbatch,
				//> foreign key(idx_inbatch) references
				//>     inbatch(idx_inbatch)
				sb.WriteString(fmt.Sprintf(
					"%s integer not null references %s,\n",
					refgField, refdTable))
				sb.WriteString(fmt.Sprintf(
					"foreign key(%s) references %s(%s),\n",
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
					sb.WriteString(fmt.Sprintf(refgField+
						" integer not null "+
						"references %s,\n", refdTable))
					sb.WriteString(fmt.Sprintf("foreign "+
						"key(%s) references %s(%s),\n",
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
			panic(pCS.Datatype)
		}
	}
	// trim off final ",\n"
	ss := sb.String()
	// and add STRICT 
	sb3 := ss[0:len(ss)-2] + "\n) STRICT;\n"
	return sb3, nil
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
