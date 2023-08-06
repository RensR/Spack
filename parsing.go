package main

import (
	"errors"
	"regexp"
)

func ParseStruct(structDefString string) (StructDef, error) {
	structName, err := parseStructName(structDefString)
	if err != nil {
		return StructDef{}, err
	}

	return StructDef{
		Name:   structName,
		Fields: parseStructFields(structDefString),
	}, nil
}

func parseStructName(structString string) (string, error) {
	structNameRegex, err := regexp.Compile(`struct ([a-zA-Z0-9]*) {`)
	if err != nil {
		return "", err
	}

	nameMatch := structNameRegex.FindStringSubmatch(structString)
	if len(nameMatch) < 2 {
		return "", errors.New("could not find struct name")
	}

	return nameMatch[1], nil
}

func parseStructFields(structString string) []DataDef {
	componentRegex := regexp.MustCompile(`\s*([a-zA-Z0-9\][.]+)\s+([a-zA-Z0-9_]+)\s*;[ \t]*(?://)?(.*)\n`)
	matches := componentRegex.FindAllStringSubmatch(structString, -1)

	var fields []DataDef
	for _, match := range matches {
		dataType := DataType(match[1])
		newField := DataDef{
			Name:    match[2],
			Type:    dataType,
			Comment: match[3],
			Size:    dataType.Size(),
		}

		fields = append(fields, newField)
	}
	return fields
}
