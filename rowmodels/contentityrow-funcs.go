package rowmodels

// Implement interface RowModeler

func (cro *ContentityRow) TableDetails() TableDetails {
     return TableDetailsCNT
}

// ColumnPtrs implements interface [Row].
// TODO: Still can't sort out the notation for ptr constraints ?!
func ColumnPtrsCNT(cro *ContentityRow) []any { // barfs on []db.PtrFields
	return []any{
		// &cro.Idx_Contentity,
		&cro.Idx_Inbatch,
		&cro.PathProps.RelFP, &cro.PathProps.AbsFP,
		// &cro.RelFP, &cro.AbsFP,
		// &cro.FUPP.RelFP, &cro.FUPP.AbsFP,
		&cro.Descr, &cro.T_Cre, &cro.T_Imp, &cro.T_Edt,
		&cro.PathProps.TypedRaw.MarkupType,
		&cro.PathAnalysis.ContypingInfo.MimeType,
		&cro.PathAnalysis.ContypingInfo.MType,
		&cro.PathProps.TypedRaw.Raw}
}

// ColumnPtrs implements interface [Row].
// TODO: Still can't sort out the notation for ptr constraints ?!
func (cro *ContentityRow) ColumnPtrs() []any { // barfs on []db.PtrFields
	return []any{
		// &cro.Idx_Contentity,
		&cro.Idx_Inbatch,
		&cro.PathProps.RelFP, &cro.PathProps.AbsFP,
		// &cro.RelFP, &cro.AbsFP,
		// &cro.FUPP.RelFP, &cro.FUPP.AbsFP,
		&cro.Descr, &cro.T_Cre, &cro.T_Imp, &cro.T_Edt,
		&cro.PathProps.TypedRaw.MarkupType,
		&cro.PathAnalysis.ContypingInfo.MimeType,
		&cro.PathAnalysis.ContypingInfo.MType,
		&cro.PathProps.TypedRaw.Raw}
}

