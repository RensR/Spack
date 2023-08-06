package main

import (
	"fmt"
	"strings"
)

func printStorageSlots(storageSlots []StorageSlot, maxFieldNameSize int) string {
	var output string

	for _, slot := range storageSlots {
		for i, field := range slot.Fields {
			capChar := " "
			fillerChar := " "

			if len(slot.Fields) > 1 {
				if i == 0 {
					capChar = "┐"
					fillerChar = "─"
				} else if i == len(slot.Fields)-1 {
					capChar = "┘"
					fillerChar = "─"
				} else {
					capChar = "│"
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
