package proto

import (
	"strconv"
	"strings"
)

const SEPARATOR string = "\r\n"

func BuildArray(content []string) string {
	var arrayBuilder strings.Builder

	arrayBuilder.WriteString("*")
	arrayBuilder.WriteString(strconv.Itoa(len(content)))
	arrayBuilder.WriteString(SEPARATOR)
	for _, str := range content {
		arrayBuilder.WriteString("$")
		length := strconv.Itoa(len(str))
		arrayBuilder.WriteString(length)
		arrayBuilder.WriteString(SEPARATOR)
		arrayBuilder.WriteString(str)
		arrayBuilder.WriteString(SEPARATOR)
	}

	return arrayBuilder.String()
}

func BuildBulkString(content string) string {
	length := strconv.Itoa(len(content))
	bs := "$"
	bs += length
	bs += SEPARATOR
	bs += content
	bs += SEPARATOR
	return bs
}

func BuildSimpleString(content string) string {
	ss := "+"
	ss += content
	ss += SEPARATOR
	return ss
}

func BuildSimpleError(prefix, body string) string {
	se := "-ERR"
	se += " " + prefix
	se += " " + body
	se += SEPARATOR
	return se
}

func BuildNULLValue() string {
	null := "$"
	null += "-1"
	null += SEPARATOR
	return null
}

func BuildInteger(val int) string {
	integer := ":"
	integer += strconv.Itoa(val)
	integer += SEPARATOR
	return integer
}
