package main

import (
	"fmt"
	"os"
)

func main() {
	filePath := os.Args[1]
	dat, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	parseAndPrint(string(dat))
}

func parseAndPrint(input string) {
	dataDef, err := ParseStruct(input)
	if err != nil {
		fmt.Printf("Error parsing struct: %s\n", err)
		return
	}

	dataDef.StorageSlots = packStructCurrentFieldOrder(dataDef.Fields)

	fmt.Print("\nStruct prior to packing:\n")
	fmt.Print(dataDef.PrintStats())
	fmt.Print(dataDef.ToString())

	dataDef.packStructOptimal()

	fmt.Print("\nStruct after packing:\n")
	fmt.Print(dataDef.PrintStats())
	fmt.Print(dataDef.ToString())
}

type StructDef struct {
	Name         string
	Fields       []DataDef
	StorageSlots []StorageSlot
}

type StorageSlot struct {
	Offset uint8
	Fields []DataDef
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
