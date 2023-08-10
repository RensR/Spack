package main

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

type StorageSlot struct {
	Offset uint8
	Fields []DataDef
}

type StructDef struct {
	Name         string
	Fields       []DataDef
	StorageSlots []StorageSlot
}

func (sd *StructDef) PrintStats() string {
	return fmt.Sprintf("Struct %s has %d fields packed into %d slots\n\n", sd.Name, len(sd.Fields), len(sd.StorageSlots))
}

func (sd *StructDef) ToString() string {
	return fmt.Sprintf("struct %s {\n%s}\n", sd.Name, printStorageSlots(sd.StorageSlots, sd.MaxFieldNameLength()))
}

func (sd *StructDef) MaxFieldNameLength() int {
	size := 0
	for _, field := range sd.Fields {
		fieldLength := field.FieldNameLength()
		if fieldLength > size {
			size = fieldLength
		}
	}
	return size
}

type DataDef struct {
	Name    string
	Comment string
	Type    DataType
	Size    uint8
}

func (dd *DataDef) FieldNameLength() int {
	return len(dd.Name) + len(dd.Type) + 1
}

func (dd *DataDef) ToString() string {
	return fmt.Sprintf("%s %s", dd.Type, dd.Name)
}

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
