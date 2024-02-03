package rowmodels

import (
	"fmt"
	D "github.com/fbaube/dsmnd"
	// FU "github.com/fbaube/fileutils"
)

// SqlarRow gas structure defined by SQLite.
type SqlarRow struct {
	/* OBS
	Idx_Contentity int
	Idx_Inbatch    int // NOTE: Maybe rename to FILESET. Could be multiple?
	Descr          string
	// FU.PathProps
	Times
	FU.PathProps
	// PathAnalysis is a ptr, so that we get a
	// NPE if it is not initialized properly.
	*FU.PathAnalysis
	*/
}

// String implements Stringer.
func (p *SqlarRow) String() string {
	return fmt.Sprintf("PP<%s> AR <%s>", "", "")
}

// ColumnSpecs_SqlarRow specifies TBS.
var ColumnSpecs_SqlarRow = []D.ColumnSpec{
	/* OBS
	D.DD_RelFP,
	D.DD_AbsFP,
	D.DD_T_Cre,
	D.DD_T_Imp,
	D.DD_T_Edt,
	*/
}

// TableDetails_SqlarRow specifies TBS.
var TableDetails_SqlarRow = TableDetails{
    	D.TableSummary{}, // FIXME
	"idx_sqlar", // IDName
	"FIXME",     // ColumnNames
	//
	ColumnSpecs_SqlarRow, // []D.ColumnSpec
	nil, 
}
