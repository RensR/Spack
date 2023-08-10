package config

import (
	"encoding/json"
	"io"
	"os"
)

type GlobalConfig struct {
	PrintingConfig PrintingConfig
}

type PrintingConfig struct {
	EnableStructPackingComments bool
	StripComments               bool
	CharacterConfig             CharacterConfig
}

type CharacterConfig struct {
	HorizontalLineChar string
	VerticalLineChar   string

	TopCapChar          string
	BottomCapChar       string
	UnpackedSlotCapChar string
	UnpackedLineChar    string
	EmptySpaceChar      string
}

func LoadConfigFromFile(configFile string) (GlobalConfig, error) {
	fileContent, err := os.Open(configFile)

	if err != nil {
		return GlobalConfig{}, err
	}

	defer fileContent.Close()

	byteResult, err := io.ReadAll(fileContent)
	if err != nil {
		return GlobalConfig{}, err
	}

	var globalConfig GlobalConfig
	err = json.Unmarshal(byteResult, &globalConfig)
	if err != nil {
		return GlobalConfig{}, err
	}

	return globalConfig, nil
}
