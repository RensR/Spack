package config

type GlobalConfig struct {
	PrintingConfig PrintingConfig
}

type PrintingConfig struct {
	EnableStructPackingComments bool
	StripComments               bool
	PrintingCharacterConfig     PrintingCharacterConfig
}

type PrintingCharacterConfig struct {
	HorizontalLineChar string
	VerticalLineChar   string

	TopCapChar          string
	BottomCapChar       string
	UnpackedSlotCapChar string
	UnpackedLineChar    string
	EmptySpaceChar      string
}
