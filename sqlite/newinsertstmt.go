package sqlite

import (
	"fmt"
	"os"
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	S "strings"
	CT "github.com/fbaube/ctoken"
	// "time"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

// case *FU.AbsFilePath, *SU.MarkupType, *CT.Raw:

// func AFPval(p *FU.AbsFilePath) string {
func AFPval(p any) string {
     var pAFP *FU.AbsFilePath
     pAFP = p.(*FU.AbsFilePath)
     var AFP FU.AbsFilePath
     AFP = *pAFP
     return string(AFP) 
}
// func MTval(p *SU.MarkupType) string {
func MTval(p any) string {
     var pAFP *SU.MarkupType
     pAFP = p.(*SU.MarkupType)
     var AFP SU.MarkupType
     AFP = *pAFP
     return string(AFP) 
}
// func CTRval(p *CT.Raw) string {
func CTRval(p any) string {
     var pAFP *CT.Raw
     pAFP = p.(*CT.Raw)
     var AFP CT.Raw
     AFP = *pAFP
     return string(AFP) 
}

// NewInsertStmt implements
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
func (pSR *SqliteRepo) NewInsertStmt(pRM DRM.RowModel) (string, error) {
	// So, we have to
	//  1) Check the struct-type of the RowModeler-instance
	//  2) Fetch the FieldPtrs
	//  3) Write out the fields to the stmt
	//  4) Write out the values to the stmt
	//  5) Add RETURNING
	//  6) (the stmt's user) Use that returned that ID

	var colPtrs []any
	var pTD DRM.TableDetails
	var now = SU.Now() // time.Now().UTC().Format(time.RFC3339)
	switch pRM.(type) {
	case *DRM.ContentityRow:
	     var pCR *DRM.ContentityRow 
	     pCR = pRM.(*DRM.ContentityRow)
	     pTD = pCR.TableDetails()
	     colPtrs = DRM.ColumnPtrsFuncCNT(pCR, false)
	     pCR.T_Cre = now
	     pCR.T_Imp = now
	     pCR.T_Edt = now
	case *DRM.InbatchRow:
	     var pIR *DRM.InbatchRow
	     pIR = pRM.(*DRM.InbatchRow)
	     pTD = pIR.TableDetails()
	     colPtrs = DRM.ColumnPtrsFuncINB(pIR, false) // not PK
	     pIR.T_Cre = now
	case *DRM.TopicrefRow:
	     var pTR *DRM.TopicrefRow
	     pTR = pRM.(*DRM.TopicrefRow)
	     pTD = pTR.TableDetails()
	     colPtrs = DRM.ColumnPtrsFuncTRF(pTR, false)
	}
	// colPtrs = pTD.ColumnPtrsFunc() // should be a no-arg func 
	fmt.Fprintf(os.Stderr, "LENS: ColSpex<%d> ColPtrs<%d> \n",
		len(pTD.ColumnSpecs), len(colPtrs))
	var sb S.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(pTD.TableSummary.StorName)
	sb.WriteString("(")
	sb.WriteString(pTD.ColumnNamesCSV)
	sb.WriteString(") VALUES(")
	// var sft D.SemanticFieldType
	var sn string
	var dt D.SemanticFieldType 
	for iii, cp := range colPtrs {
	    /*
	    if iii == 0 {
	       sn = "priKey"
	       dt = D.SFT_PRKEY
	       // This is an INSERT, so we do
	       // not write out the primary key!
	       continue
	    }
	    */
	    sn = pTD.ColumnSpecs[iii].StorName
	    dt = D.SemanticFieldType(pTD.ColumnSpecs[iii].Datatype)
	    fmt.Printf("[%d] %s / %s / %T \n", iii, sn, dt, cp)
	    // sft = D.SemanticFieldType(ppp.Datatype)
	    switch cp.(type) {
	    	   case *string:
		   	var pS *string
			var sS string
			pS = cp.(*string)
			sS = *pS
		   	sb.WriteString(fmt.Sprintf("'%s', ", sS))
	    	   case *FU.AbsFilePath:
		   	sb.WriteString(fmt.Sprintf("'%s', ", AFPval(cp)))
	    	   case *SU.MarkupType:
		   	sb.WriteString(fmt.Sprintf("'%s', ", MTval(cp)))
	    	   case *CT.Raw:
		   	sb.WriteString(fmt.Sprintf("'%s', ", CTRval(cp)))
		   case *int:
		   	var pI *int
			pI = cp.(*int)
		   	sb.WriteString(fmt.Sprintf("%d, ", *pI))
	    }
	}
	var stmt string 
	stmt = sb.String()
	stmt2 := stmt[:len(stmt)-2] + ") RETURNING IDX_" + pTD.StorName + ";"
	// sb.WriteString(") RETURNING IDX_" + pTD.StorName + ";")
	println("INSERT STMT:", stmt2) // sb.String())
	return stmt2, nil
}

/*
// func Deref[T any](ptr *T) T { 
func StringFromPtr[SS ~string](ss *SS) string {

     // cannot use *ss (variable of type SS constrained
     // by ~string) as string value in return statement
     return *ss

     var pS *string
     // cannot convert ss (variable of type *SS) to type *string 
     pS = (*string)(ss)
     return *pS

     var SSS string
     // cannot use *ss (variable of type SS constrained
     // by ~string) as string value in assignment
     SSS = *ss
     return SSS
}
*/

/*
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
		* /
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

func NewInsertStmtGnrcFunc[T DRM.RowModel](pSR *SqliteRepo, pRM T) (string, error) {
	// So, we have to
	//  1) Check the struct-type of the RowModeler-instance
	//  2) Fetch the FieldPtrs
	//  3) Write out the fields to the stmt
	//  4) Write out the values to the stmt
	//  5) Add RETURNING
	//  6) (the stmt's user) Use that returned that ID

	var colPtrs []any
	var pTD DRM.TableDetails
	// TMP var now = SU.Now() // time.Now().UTC().Format(time.RFC3339)

	pTD = pRM.TableDetails()
	colPtrs = pTD.ColumnPtrsFunc(pRM, false)
	/* TMP
	pRM.T_Cre = now
	pRM.T_Imp = now
	pRM.T_Edt = now
	*/

	fmt.Fprintf(os.Stderr, "LENS: ColSpex<%d> ColPtrs<%d> \n",
		len(pTD.ColumnSpecs), len(colPtrs))
	var sb S.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(pTD.TableSummary.StorName)
	sb.WriteString("(")
	sb.WriteString(pTD.ColumnNamesCSV)
	sb.WriteString(") VALUES(")
	// var sft D.SemanticFieldType
	var sn string
	var dt D.SemanticFieldType 
	for iii, cp := range colPtrs {
	    /*
	    if iii == 0 {
	       sn = "priKey"
	       dt = D.SFT_PRKEY
	       // This is an INSERT, so we do
	       // not write out the primary key!
	       continue
	    }
	    */
	    sn = pTD.ColumnSpecs[iii].StorName
	    dt = D.SemanticFieldType(pTD.ColumnSpecs[iii].Datatype)
	    fmt.Printf("[%d] %s / %s / %T \n", iii, sn, dt, cp)
	    // sft = D.SemanticFieldType(ppp.Datatype)
	    switch cp.(type) {
	    	   case *string:
		   	var pS *string
			var sS string
			pS = cp.(*string)
			sS = *pS
		   	sb.WriteString(fmt.Sprintf("'%s', ", sS))
	    	   case *FU.AbsFilePath:
		   	sb.WriteString(fmt.Sprintf("'%s', ", AFPval(cp)))
	    	   case *SU.MarkupType:
		   	sb.WriteString(fmt.Sprintf("'%s', ", MTval(cp)))
	    	   case *CT.Raw:
		   	sb.WriteString(fmt.Sprintf("'%s', ", CTRval(cp)))
		   case *int:
		   	var pI *int
			pI = cp.(*int)
		   	sb.WriteString(fmt.Sprintf("%d, ", *pI))
	    }
	}
	var stmt string 
	stmt = sb.String()
	stmt2 := stmt[:len(stmt)-2] + ") RETURNING IDX_" + pTD.StorName + ";"
	// sb.WriteString(") RETURNING IDX_" + pTD.StorName + ";")
	println("INSERT STMT:", stmt2) // sb.String())
	return stmt2, nil
}

