package utils

import "strconv"

// $<length>\r\n<data>\r\n
func GetBulkString(input string) string {
	return "$" + strconv.Itoa(len(input)) + CLRF + input + CLRF
}

func GetArray(input []string) string {
	op := "*" + strconv.Itoa(len(input)) + CLRF
	for _, val := range input {
		op += GetBulkString(val)
	}
	return op
}
