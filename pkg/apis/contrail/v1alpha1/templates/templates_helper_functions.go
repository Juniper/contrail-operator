package templates

import (
	"strconv"
	"strings"
)

func JoinListWithSeparator(items []string, separator string) string {
	return strings.Join(items, separator)
}

func JoinListWithSeparatorAndSingleQuotes(items []string, separator string) string {
	if len(items) == 0 {
		return ""
	}
	joinedList := JoinListWithSeparator(items, "'"+separator+"'")
	joinedList = "'" + joinedList + "'"
	return joinedList
}

func EndpointList(ips []string, port int) []string {
	portStr := strconv.Itoa(port)
	endpoints := []string{}
	for _, ip := range ips {
		endpoints = append(endpoints, ip+":"+portStr)
	}
	return endpoints
}
