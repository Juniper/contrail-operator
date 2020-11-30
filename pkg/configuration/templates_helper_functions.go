package configuration

import (
	"strconv"
	"strings"
)

// JoinListWithSeparator joins a slice into a string using the given separator.
func JoinListWithSeparator(items []string, separator string) string {
	return strings.Join(items, separator)
}

// JoinListWithSeparatorAndSingleQuotes joins a slice into a string using the
// given separator and surrounds each slice item with single quotes.
func JoinListWithSeparatorAndSingleQuotes(items []string, separator string) string {
	if len(items) == 0 {
		return ""
	}
	joinedList := JoinListWithSeparator(items, "'"+separator+"'")
	joinedList = "'" + joinedList + "'"
	return joinedList
}

// EndpointList creates a new slice in which each item is an ip and port joined
// with a colon.
func EndpointList(ips []string, port int) []string {
	portStr := strconv.Itoa(port)
	endpoints := []string{}
	for _, ip := range ips {
		endpoints = append(endpoints, ip+":"+portStr)
	}
	return endpoints
}
