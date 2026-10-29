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
	var sb strings.Builder

	length := strconv.Itoa(len(content))
	sb.WriteString("$")
	sb.WriteString(length)
	sb.WriteString(SEPARATOR)
	sb.WriteString(content)
	sb.WriteString(SEPARATOR)
	return sb.String()
}

func BuildSimpleString(content string) string {
	var sb strings.Builder

	sb.WriteString("+")
	sb.WriteString(content)
	sb.WriteString(SEPARATOR)
	return sb.String()
}

func BuildSimpleError(prefix, body string) string {
	var sb strings.Builder

	sb.WriteString("-ERR")
	sb.WriteString(" " + prefix)
	sb.WriteString(" " + body)
	sb.WriteString(SEPARATOR)
	return sb.String()
}

func BuildNULLValue() string {
	var sb strings.Builder

	sb.WriteString("$")
	sb.WriteString("-1")
	sb.WriteString(SEPARATOR)
	return sb.String()
}

func BuildInteger(val int) string {
	var sb strings.Builder

	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(val))
	sb.WriteString(SEPARATOR)
	return sb.String()
}
