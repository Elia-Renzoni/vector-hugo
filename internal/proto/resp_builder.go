package proto

import (
	"bytes"
	"strconv"
)

const SEPARATOR string = "\r\n"

func BuildArray(content []string) []byte {
	var arrayBuilder bytes.Buffer

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

	return arrayBuilder.Bytes()
}

func BuildBulkString(content string) []byte {
	var sb bytes.Buffer

	length := strconv.Itoa(len(content))
	sb.WriteString("$")
	sb.WriteString(length)
	sb.WriteString(SEPARATOR)
	sb.WriteString(content)
	sb.WriteString(SEPARATOR)
	return sb.Bytes()
}

func BuildSimpleString(content string) []byte {
	var sb bytes.Buffer

	sb.WriteString("+")
	sb.WriteString(content)
	sb.WriteString(SEPARATOR)
	return sb.Bytes()
}

func BuildSimpleError(prefix, body string) []byte {
	var sb bytes.Buffer

	sb.WriteString("-ERR")
	sb.WriteString(" " + prefix)
	sb.WriteString(" " + body)
	sb.WriteString(SEPARATOR)
	return sb.Bytes()
}

func BuildNULLValue() []byte {
	var sb bytes.Buffer

	sb.WriteString("$")
	sb.WriteString("-1")
	sb.WriteString(SEPARATOR)
	return sb.Bytes()
}

func BuildInteger(val int) []byte {
	var sb bytes.Buffer

	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(val))
	sb.WriteString(SEPARATOR)
	return sb.Bytes()
}
