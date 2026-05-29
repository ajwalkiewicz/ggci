package output

type View struct {
	Columns []Column
}

var UnixView = View{
	Columns: []Column{
		UnixModeColumn,
		UnixUserColumn,
		UnixGroupColumn,
		LastWriteTimeColumn,
		UnixSizeColumn,
		NameColumn,
	},
}

var LegacyView = View{
	Columns: []Column{
		LegacyModeColumn,
		LastWriteTimeColumn,
		LegacyLengthColumn,
		NameColumn,
	},
}
