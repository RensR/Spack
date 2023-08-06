package main

import (
	"sort"
)

func packStructCurrentFieldOrder(fields []DataDef) []StorageSlot {
	var storageSlots []StorageSlot

	for _, field := range fields {
		// If we have no packed fields yet, add the first field as a packed field
		if len(storageSlots) == 0 {
			storageSlots = []StorageSlot{{Fields: []DataDef{field}, Offset: field.Size}}
			continue
		}

		lastSlot := storageSlots[len(storageSlots)-1]
		if lastSlot.Offset+field.Size <= 32 {
			storageSlots[len(storageSlots)-1].Fields = append(lastSlot.Fields, field)
			storageSlots[len(storageSlots)-1].Offset += field.Size
			continue
		}

		storageSlots = append(storageSlots, StorageSlot{Fields: []DataDef{field}, Offset: field.Size})
	}

	return storageSlots
}

func (sd *StructDef) packStructOptimal() {
	sort.Slice(sd.Fields, func(i, j int) bool {
		return sd.Fields[i].Size > sd.Fields[j].Size
	})

	sd.StorageSlots = binPacking(sd.Fields, []StorageSlot{})
}

func binPacking(fields []DataDef, existingSlots []StorageSlot) []StorageSlot {
	if len(fields) == 0 {
		return existingSlots
	}

	currentItem := fields[0]
	var packingOptions [][]StorageSlot

	for i, slot := range existingSlots {
		// The field doesn't fit into the slot, so skip it
		if slot.Offset+currentItem.Size > 32 {
			continue
		}

		// It the field does fit, make a copy of the existing slots and add the field to the slot
		slotsCopy := make([]StorageSlot, len(existingSlots))
		copy(slotsCopy, existingSlots)

		// Make sure we copy slot.Fields to avoid modifying the original slice
		slotsCopy[i].Fields = AddToCopyOfList(slot.Fields, currentItem)
		slotsCopy[i].Offset += currentItem.Size

		// Recursively call binPacking with the remaining fields and the new slot state
		packingOptions = append(packingOptions, binPacking(fields[1:], slotsCopy))
	}

	// If we can't fit the field into any existing slots, create a new slot and recursively call binPacking
	packingOptions = append(packingOptions, binPacking(fields[1:], append(existingSlots, StorageSlot{
		Fields: []DataDef{currentItem},
		Offset: currentItem.Size,
	})))

	return findOptimalPacking(packingOptions)
}

func findOptimalPacking(options [][]StorageSlot) []StorageSlot {
	if len(options) == 0 {
		return []StorageSlot{}
	}

	// First find the options with the least amount of slots
	// This gives the first option twice if it happens to be the best
	// This doesn't matter much since we only return a single option anyway
	leastAmountOfSlots := [][]StorageSlot{options[0]}

	for _, option := range options {
		if len(option) < len(leastAmountOfSlots[0]) {
			leastAmountOfSlots = [][]StorageSlot{option}
		} else if len(option) == len(leastAmountOfSlots[0]) {
			leastAmountOfSlots = append(leastAmountOfSlots, option)
		}
	}

	// Then find the options with the most completely filled slots
	bestSlots := leastAmountOfSlots[0]
	mostFilledStorageSlots := 0

	for _, option := range leastAmountOfSlots {
		filledSlots := 0
		for _, slot := range option {
			if slot.Offset == 32 {
				filledSlots++
			}
		}

		if filledSlots > mostFilledStorageSlots {
			bestSlots = option
			mostFilledStorageSlots = filledSlots
		}
	}

	return bestSlots
}

func AddToCopyOfList[T any](list []T, item T) []T {
	listCopy := make([]T, len(list))
	copy(listCopy, list)
	return append(listCopy, item)
}
