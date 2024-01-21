package rowmodels

import (
	"fmt"
	DRU "github.com/fbaube/datarepo/utils"
)

// TODO Write col desc's using Desmond !
// TODO Generate ColNames from ColumnSpecs_InbatchRow

// String implements Stringer. FIXME
func (p *InbatchRow) String() string {
	return fmt.Sprintf("inbatchrow FIXME")
	// p.PathProps.String(), p.PathAnalysis.String())
}

// TableDetailsINB TBS has no foreign keys.
var TableDetailsINB = DRU.TableDetails{
        TableSummaryINB,
	"idx_inbatch", // IDName
	ColumnNamesCsvINB, 
	ColumnSpecsINB, // []D.ColumnSpec
	// /* ColumnPtrsFunc */ ColumnPtrsINB,
}

