package utility

import (
	"strings"
)

//RemoveDuplicatesInt ...this is used to remove duplicates in integers
func RemoveDuplicatesInt(elements []int) []int {
	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

//RemoveDuplicatesUnordered .. this is used to remove duplicates in string without ordering them
func RemoveDuplicatesUnordered(elements ...string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

//StringExistInSlice ....this is to check if a string exist in a slice
func StringExistInSlice(str string, elements ...string) bool {
	for _, n := range elements {
		if str == n {
			return true
		}
	}
	return false
}

//StringSliceToString ....this is used to convert a string slice to string with comma as a seperator
func StringSliceToString(slice []string) string {
	return strings.Join(slice, ",")
}

//StringSubString ...this is for string substring
func StringSubString(str string, startAt, endsAt int) string {
	return str[startAt:endsAt]
}

// https://github.com/DaddyOh/golang-samples/blob/master/pad.go

// RightPad2Len https://github.com/DaddyOh/golang-samples/blob/master/pad.go
func RightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

// LeftPad2Len https://github.com/DaddyOh/golang-samples/blob/master/pad.go
func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}
