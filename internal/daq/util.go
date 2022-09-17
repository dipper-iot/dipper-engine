package daq

import (
	"strconv"
	"strings"
)

type TypeDataPath int8
type TypeData int8

const (
	PathData  TypeDataPath = 1
	ArrayData TypeDataPath = 2
)

const (
	Number TypeData = 1
	String TypeData = 2
	Object TypeData = 3
	Array  TypeData = 4
)

func getPath(location string) []string {
	return strings.Split(location, ".")
}

func getPathType(item string) TypeDataPath {
	if strings.Contains(item, "[") && strings.Contains(item, "]") {
		return ArrayData
	}
	return PathData
}

func getIndexArray(item string) (int, error) {
	item = strings.ReplaceAll(item, "[", "")
	item = strings.ReplaceAll(item, "]", "")
	return strconv.Atoi(item)
}
