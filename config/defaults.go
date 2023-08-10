package config

func GetDefaultConfig() GlobalConfig {
	return GlobalConfig{
		PrintingConfig: GetDefaultPrintingConfig(),
	}
}

func GetDefaultPrintingConfig() PrintingConfig {
	return PrintingConfig{
		EnableStructPackingComments: true,
		StripComments:               true,
		CharacterConfig:             GetDefaultPrintingCharConfig(),
	}
}

func GetDefaultPrintingCharConfig() CharacterConfig {
	return CharacterConfig{
		HorizontalLineChar:  "─",
		VerticalLineChar:    "│",
		TopCapChar:          "╮",
		BottomCapChar:       "╯",
		UnpackedSlotCapChar: " ",
		UnpackedLineChar:    " ",
		EmptySpaceChar:      " ",
	}
}
