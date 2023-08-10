package config

func GetDefaultConfig() GlobalConfig {
	return GlobalConfig{
		PrintingConfig: GetDefaultPrintingConfig(),
	}
}

func GetDefaultPrintingConfig() PrintingConfig {
	return PrintingConfig{
		EnableStructPackingComments: true,
		StripComments:               false,
		PrintingCharacterConfig:     GetDefaultPrintingCharConfig(),
	}
}

func GetDefaultPrintingCharConfig() PrintingCharacterConfig {
	return PrintingCharacterConfig{
		HorizontalLineChar:  "─",
		VerticalLineChar:    "│",
		TopCapChar:          "╮",
		BottomCapChar:       "╯",
		UnpackedSlotCapChar: " ",
		UnpackedLineChar:    " ",
		EmptySpaceChar:      " ",
	}
}
