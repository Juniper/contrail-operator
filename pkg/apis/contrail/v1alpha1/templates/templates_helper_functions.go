package templates

import (
	"strconv"
	"strings"
)

func IPListCommaSeparated(ips []string) string {
	return strings.Join(ips, ",")
}

func IPListSpaceSeparated(ips []string) string {
	return strings.Join(ips, " ")
}

func IPListCommaSeparatedQuoted(ips []string) string {
	if len(ips) == 0 {
		return ""
	}
	endpointList := strings.Join(ips, "','")
	endpointList = "'" + endpointList + "'"
	return endpointList
}

func EndpointList(ips []string, port int) []string {
	portStr := strconv.Itoa(port)
	endpoints := []string{}
	for _, ip := range ips {
		endpoints = append(endpoints, ip+":"+portStr)
	}
	return endpoints
}

func EndpointListCommaSeparatedQuoted(ips []string, port int) string {
	if len(ips) == 0 {
		return ""
	}
	portStr := strconv.Itoa(port)
	endpointList := strings.Join(ips, ":"+portStr+"','")
	endpointList = "'" + endpointList + ":" + portStr + "'"
	return endpointList
}

func EndpointListCommaSeparated(ips []string, port int) string {
	return endpointListWithSeparator(ips, port, ",")
}

func EndpointListSpaceSeparated(ips []string, port int) string {
	return endpointListWithSeparator(ips, port, " ")
}

func endpointListWithSeparator(ips []string, port int, separator string) string {
	if len(ips) == 0 {
		return ""
	}
	portStr := strconv.Itoa(port)
	endpointList := strings.Join(ips, ":"+portStr+separator)
	endpointList = endpointList + ":" + portStr
	return endpointList
}
