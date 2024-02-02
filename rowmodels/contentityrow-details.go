package rowmodels

// Implement interface RowModeler

// TableDetailsCNT specifies 11 DB columns,
// incl primary key (assumed) and one foreign key, "inbatch".
var TableDetailsCNT = TableDetails{
        TableSummaryCNT, 
	"idx_contentity", // IDName
	"IDX_inbatch, RelFP, AbsFP, Descr, T_Cre, T_Imp, T_Edt, " +
		"RawMT, Mimtp, MType, Contt", // ColumnNames
	// One foreign key: "inbatch"
	ColumnSpecsCNT, // []D.ColumnSpecs
}

// TableDetails returns the table
// detail info, given any instance. 
func (cro *ContentityRow) TableDetails() TableDetails {
     return TableDetailsCNT
}

