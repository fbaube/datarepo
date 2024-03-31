package sqlite

import (
	"fmt"
	// "os"
	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	S "strings"
	CT "github.com/fbaube/ctoken"
	// "time"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

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
	/*
	colPtrs = pTD.ColumnPtrsFunc() // should be a no-arg func 
	fmt.Fprintf(os.Stderr, "LENS: ColSpex<%d> ColPtrs<%d> \n",
		len(pTD.ColumnSpecs), len(colPtrs))
	*/
	var sb S.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(pTD.TableSummary.StorName)
	sb.WriteString("(")
	sb.WriteString(pTD.ColumnNamesCSV)
	sb.WriteString(") VALUES(")
	// var sft D.SemanticFieldType
	var sn string
	var dt D.SemanticFieldType
	fmt.Fprintf(pSR.w, "=== %s.NewInsStmt.ColSpex ===\n")
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
	    fmt.Fprintf(pSR.w, "[%d] %s / %s / %T \n", iii, sn, dt, cp)
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
	fmt.Fprintf(pSR.w, "INSERT STMT: %.60s[...] \n", stmt2) // sb.String())
	return stmt2, nil
}

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

	// fmt.Fprintf(os.Stderr, "LENS: ColSpex<%d> ColPtrs<%d> \n",
	//	len(pTD.ColumnSpecs), len(colPtrs))

	// Add some log info 
	fmt.Fprintf(pSR.w, "=== %s.NewInsStmtGenc.ColSpex ===\n", pTD.StorName)
	// Build the statement  
	var sqlBldr S.Builder
	sqlBldr.WriteString("INSERT INTO ")
	sqlBldr.WriteString(pTD.TableSummary.StorName)
	sqlBldr.WriteString("(")
	sqlBldr.WriteString(pTD.ColumnNamesCSV)
	sqlBldr.WriteString(") VALUES(")
	// var sft D.SemanticFieldType
	var sn string
	var dt D.SemanticFieldType

	// ======================
	//  FIRST, MAKE THE STMT
	//   USING PLACEHOLDERS 
	// ======================

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
	    // Add some log info 
	    fmt.Fprintf(pSR.w, "[%d] %s / %s / %T \n", iii, sn, dt, cp)
	    // sft = D.SemanticFieldType(ppp.Datatype)
	    switch cp.(type) {
	    	   case *string:
		   	var pS *string
			var sS string
			pS = cp.(*string)
			sS = *pS
		   	sqlBldr.WriteString(fmt.Sprintf("'%s', ", sS))
	    	   case *FU.AbsFilePath:
		   	sqlBldr.WriteString(fmt.Sprintf("'%s', ", AFPval(cp)))
	    	   case *SU.MarkupType:
		   	sqlBldr.WriteString(fmt.Sprintf("'%s', ", MTval(cp)))
	    	   case *CT.Raw:
		   	sqlBldr.WriteString(fmt.Sprintf("'%s', ", CTRval(cp)))
		   case *int:
		   	var pI *int
			pI = cp.(*int)
		   	sqlBldr.WriteString(fmt.Sprintf("%d, ", *pI))
	    }
	}
	var stmt string 
	stmt = sqlBldr.String()
	stmt2 := stmt[:len(stmt)-2] + ") RETURNING IDX_" + pTD.StorName + ";"
	// sqlBldr.WriteString(") RETURNING IDX_" + pTD.StorName + ";")
	fmt.Fprintf(pSR.w, "=== %s.NewInsStmtGenc.SQL ===\n%s\n",
		pTD.StorName, stmt2)
	return stmt2, nil
}

	// ================================
	//  THEN, FILL IN THE PLACEHOLDERS
	//  (so that we get the defense
	//   against SQL injection attacks!)
	// ================================
