package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type DataType string

const (
	_string  DataType = "string"
	_bytes   DataType = "bytes"
	_bool    DataType = "bool"
	_int     DataType = "int"
	_uint    DataType = "uint"
	_address DataType = "address"
)

func (dt *DataType) Size() uint8 {
	var dataTypeSize = map[DataType]uint8{
		_string:  32,
		_bytes:   32,
		_bool:    1,
		_int:     32,
		_uint:    32,
		_address: 20,
	}

	if size, ok := dataTypeSize[*dt]; ok {
		return size
	}

	if strings.HasSuffix(string(*dt), "[]") {
		return 32
	}
	intRegex := regexp.MustCompile(`u?int(\d+)`)
	if intRegex.MatchString(string(*dt)) {
		size, _ := strconv.Atoi(intRegex.FindStringSubmatch(string(*dt))[1])
		return uint8(size / 8)
	}
	bytesRegex := regexp.MustCompile(`bytes(\d+)`)
	if bytesRegex.MatchString(string(*dt)) {
		size, _ := strconv.Atoi(bytesRegex.FindStringSubmatch(string(*dt))[1])
		return uint8(size)
	}

	fmt.Printf("Unknown data type %s. Assuming 32 bytes. If this is not correct add them with the correct size\n", *dt)
	return 32
}
