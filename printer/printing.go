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
		for i, field := range slot.Fields {
			capChar := p.PrintingCharacterConfig.UnpackedSlotCapChar
			fillerChar := p.PrintingCharacterConfig.UnpackedLineChar

			if len(slot.Fields) > 1 {
				if i == 0 {
					capChar = p.PrintingCharacterConfig.TopCapChar
					fillerChar = p.PrintingCharacterConfig.HorizontalLineChar
				} else if i == len(slot.Fields)-1 {
					capChar = p.PrintingCharacterConfig.BottomCapChar
					fillerChar = p.PrintingCharacterConfig.HorizontalLineChar
				} else {
					capChar = p.PrintingCharacterConfig.VerticalLineChar
					fillerChar = p.PrintingCharacterConfig.EmptySpaceChar
				}
			}

			spacingAndFieldDef := fmt.Sprintf("   %s; // ", field.ToString())
			fillerCount := maxFieldNameSize - field.FieldNameLength() + 2
			structPackingComment := fmt.Sprintf("%s%s%s\n", strings.Repeat(fillerChar, fillerCount), capChar, field.Comment)
			output += spacingAndFieldDef + structPackingComment
		}
	}

	return output
}

func (p *Printer) PrintSolidityStruct(structDef solidity.Struct) string {
	return fmt.Sprintf("struct %s {\n%s}\n", structDef.Name, p.printStorageSlots(structDef.StorageSlots, structDef.MaxFieldNameLength()))
}
