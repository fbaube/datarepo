package sqlite

import (
	"fmt"
	D "github.com/fbaube/dsmnd"
	// FU "github.com/fbaube/fileutils"
	// DR "github.com/fbaube/datarepo"
	S "strings"
	// "time"
)

/* REF
https://www.sqlite.org/lang_createtable.html
Create Table [ If Not Exists ] [schemaname.]tablename
1) as select-statement
2) ( columndef,+ tableconstraint,* )
[ table-options ] ;
*/

/* REF
TABL Fundatype = "tabl" // This datum describes a table !
INTG = "1234" // SQLITE_INTEGER 1
FLOT = "1.0f" // SQLITE_FLOAT   2
TEXT = "AaZz" // SQLITE_TEXT    3
BLOB = "blob" // SQLITE_BLOB    4
NULL = "null" // SQLITE_NULL    5
PKEY = "pkey" // PRIMARY KEY (SQLite "INTEGER")
FKEY = "fkey" // FOREIGN KEY (SQLite "INTEGER")
LIST = "list" // table type, one per enumeration ??
OTHR = "othr"
*/

// BuildInsertStmt implements
// https://www.sqlite.org/lang_insert.html
// 
// INSERT INTO table [ ( column-name,+ ) ] VALUES(...);
//
// This creates one or more new rows in an existing table:
//  - If the column-name list after table-name is omitted, 
//    then the number of values inserted into each row must 
//    be the same as the number of columns in the table. 
//  - If a column-name list is specified, then the number 
//    of values in each term of the VALUE list must match 
//    the number of specified columns. 
//
// Table columns that do not appear in the column list are populated 
// with the default column value (specified as part of the CREATE 
// TABLE statement), or with NULL if no default value is specified.
// .
func (pSR *SqliteRepo) BuildInsertStmt(RM RowModeler) (string, error) {
	var sb, sb2 S.Builder
	panic("FIXME")
	var pTD = RM.TableDetails()
	sb.WriteString(fmt.Sprintf("CREATE TABLE %s(\n", pTD.StorName))
	sb2.WriteString("reposqlite.GenInsertStmt: ")
	for _, pCS := range pTD.ColumnSpecs {
		cnm := pCS.StorName // column name
		bdt := pCS.Datatype
		sb2.WriteString(fmt.Sprintf("%s:%s, ", cnm, bdt))
	}
	fmt.Printf(sb2.String() + "\n")

	// PRIMARY KEY IS ASSUMED - DO IT FIRST
	// idx_mytable integer not null primary key autoincrement,
	sb.WriteString(pTD.IDName +
		" integer not null primary key autoincrement, " +
		"-- NOTE: integer, not int. \n")

	for _, pCS := range pTD.ColumnSpecs {
		colName := pCS.StorName // column name in DB
		fmt.Sprintf("Creating column: %s \n", pCS.String())
		SFT := D.SemanticFieldType(pCS.Datatype)
                BDT := SFT.BasicDatatype()

                switch BDT {

		case "PRKEY": // D.PKEY:
			panic("DUPE PRIMARY KEY")

		case "FRKEY": // D.FKEY:
			//> D.ColumnSpec{D.FKEY, "idx_inbatch", "inbatch",
			//>  "Input batch of imported content"},
			// referencing fields's name is idx_inbatch
			refgField := colName
			// referenced table's name is inbatch
			refdTable := pCS.DispName

			// Count up underscores (per comment above)
			ss := S.Split(refgField, "_")
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
		case D.BDT_TEXT:
			sb.WriteString(colName + " text not null,\n")
		case D.BDT_INTG:
			// filect int not null check (filect >= 0) default 0
			sb.WriteString(colName + " int not null,\n")
		/* Unimplem'd:
		case D.FLOT:
		case D.BLOB:
		case D.NULL:
		case D.LIST:
		case D.OTHR:
		*/
		default:
			panic(pCS.Datatype)
		}
	}
	// trim off final ",\n"
	ss := sb.String()
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
