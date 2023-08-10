package printer

import (
	"fmt"
	"strings"

	"spack/config"
	"spack/solidity"
)

type Printer struct {
	config.PrintingConfig
}

func NewPrinter(printingConfig config.PrintingConfig) (Printer, error) {
	return Printer{PrintingConfig: printingConfig}, nil
}

func (p *Printer) printStorageSlots(storageSlots []solidity.StorageSlot, maxFieldNameSize int) string {
	var output string

	for _, slot := range storageSlots {
		for i, _ := range slot.Fields {
			output += p.PrintSingleLine(slot.Fields, i, maxFieldNameSize)
		}
	}

	return output
}

func (p *Printer) PrintSingleLine(fields []solidity.DataDef, position int, maxFieldNameSize int) string {
	field := fields[position]

	if !p.PrintingConfig.EnableStructPackingComments {
		if p.PrintingConfig.StripComments {
			return fmt.Sprintf("   %s;\n", field.ToString())
		}
		return fmt.Sprintf("   %s; // %s\n", field.ToString(), field.Comment)
	}

	capChar := p.CharacterConfig.UnpackedSlotCapChar
	fillerChar := p.CharacterConfig.UnpackedLineChar

	if len(fields) > 1 {
		if position == 0 {
			capChar = p.CharacterConfig.TopCapChar
			fillerChar = p.CharacterConfig.HorizontalLineChar
		} else if position == len(fields)-1 {
			capChar = p.CharacterConfig.BottomCapChar
			fillerChar = p.CharacterConfig.HorizontalLineChar
		} else {
			capChar = p.CharacterConfig.VerticalLineChar
			fillerChar = p.CharacterConfig.EmptySpaceChar
		}
	}
	fieldComment := ""
	if !p.PrintingConfig.StripComments {
		fieldComment = field.Comment
	}

	fillerCount := maxFieldNameSize - field.FieldNameLength() + p.CharacterConfig.MinLinePadding
	structPackingComment := fmt.Sprintf("%s%s%s\n", strings.Repeat(fillerChar, fillerCount), capChar, fieldComment)
	spacingAndFieldDef := fmt.Sprintf("   %s; // ", field.ToString())
	return spacingAndFieldDef + structPackingComment
}

func (p *Printer) PrintSolidityStruct(structDef solidity.Struct) string {
	return fmt.Sprintf("struct %s {\n%s}\n", structDef.Name, p.printStorageSlots(structDef.StorageSlots, structDef.MaxFieldNameLength()))
}
