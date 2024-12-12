package utils

import (
	"log"
	"strconv"
	"strings"
)

// UniqueIDFromFileName gets id from file name,
// if while adding, file with such name was added before.
// test_file.txt (2) tells, that this is third file with name "test_file.txt"
// result of this function will be: 2
// if file doesn't have such id, function will return 0
func UniqueIDFromFileName(fileName string) int {
	values := strings.Split(fileName, " ")

	suffixValue := values[len(values)-1] // last value

	suffixWithoutPrefix, found := strings.CutPrefix(suffixValue, "(")
	if !found { // tells that this name doesn't have suffix like (n)
		return 0
	}

	suffix, found := strings.CutSuffix(suffixWithoutPrefix, ")")
	if !found {
		return 0
	}

	id, err := strconv.Atoi(suffix)
	if err != nil {
		log.Printf("error convertins suffix to integer: %v", err)
		return 0
	}

	return id
}
