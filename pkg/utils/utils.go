package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// NextUniqueName gets next unique name from file name,
// if while adding, file with such name was added before.
// test_file (2).txt tells, that this is third file with name "test_file.txt"
// result of this function will be: test_file (3).txt
// if file doesn't have such id, function will return new name with id = 1
func NextUniqueName(fileName string) string {
	values := strings.Split(fileName, " ")
	if len(values) < 1 {
		return generateWithNonExistingID(fileName)
	}

	suffixValue := values[len(values)-1] // last value: test_file (2).txt -> (2).txt
	fmt.Println(suffixValue)

	withoutFormatValues := strings.Split(suffixValue, ".") // (2).txt -> [(2), txt]
	if len(withoutFormatValues) < 2 {
		return generateWithNonExistingID(fileName)
	}

	uniqueID := withoutFormatValues[0] // [(2), txt] -> (2)
	fmt.Println(uniqueID)

	uniqueIDWithoutPrefix, found := strings.CutPrefix(uniqueID, "(") // (2) -> 2)
	if !found {                                                      // tells that this name doesn't have suffix like (n)
		return generateWithNonExistingID(fileName)
	}

	nowUniqueIDString, found := strings.CutSuffix(uniqueIDWithoutPrefix, ")") // 2) -> 2
	if !found {
		return generateWithNonExistingID(fileName)
	}

	nowUniqueID, err := strconv.Atoi(nowUniqueIDString)
	if err != nil {
		return generateWithNonExistingID(fileName)
	}

	fmt.Println(nowUniqueID)

	nextUniqueName := fmt.Sprintf(
		"%s (%d).%s",
		strings.Join(values[:len(values)-1], " "),
		nowUniqueID+1,
		strings.Join(withoutFormatValues[1:], "."),
	)

	return nextUniqueName
}

func generateWithNonExistingID(name string) string {
	values := strings.Split(name, ".")               // test.txt -> [test, txt]
	values[0] = fmt.Sprintf("%s (%d)", values[0], 1) // [test, txt] -> [test (1), txt]

	return strings.Join(values, ".") // [test (1), txt] -> test (1).txt
}
