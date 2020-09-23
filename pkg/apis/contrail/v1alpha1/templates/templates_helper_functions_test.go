package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPListCommaSeparated(t *testing.T) {
	tests := []struct {
		name           string
		inputIpList    []string
		expectedOutput string
	}{
		{
			name:           "Should return empty string for empty list.",
			inputIpList:    []string{},
			expectedOutput: "",
		},
		{
			name:           "Shouldn't add commas to single ip.",
			inputIpList:    []string{"1.1.1.1"},
			expectedOutput: "1.1.1.1",
		},
		{
			name:           "Should add commas between multiple ips",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			expectedOutput: "1.1.1.1,2.2.2.2,3.3.3.3",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := IPListCommaSeparated(test.inputIpList)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}

func TestIPListCommaSeparatedQuoted(t *testing.T) {
	tests := []struct {
		name           string
		inputIpList    []string
		expectedOutput string
	}{
		{
			name:        "Should return empty string for empty list.",
			inputIpList: []string{},
		},
		{
			name:           "Shouldn't add commas to single ip and should add quotes",
			inputIpList:    []string{"1.1.1.1"},
			expectedOutput: "'1.1.1.1'",
		},
		{
			name:           "Should add commas between multiple ips",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			expectedOutput: "'1.1.1.1','2.2.2.2','3.3.3.3'",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := IPListCommaSeparatedQuoted(test.inputIpList)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}

func TestEndpointList(t *testing.T) {
	tests := []struct {
		name           string
		inputIpList    []string
		inputPort      int
		expectedOutput []string
	}{
		{
			name:           "Should return empty slice for empty input slice.",
			inputIpList:    []string{},
			inputPort:      1234,
			expectedOutput: []string{},
		},
		{
			name:           "Should return slice with single endpoint for single ip in input slice.",
			inputIpList:    []string{"1.1.1.1"},
			inputPort:      1234,
			expectedOutput: []string{"1.1.1.1:1234"},
		},
		{
			name:           "Should return slice with multiple endpoints for multiple ips in input slice.",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			inputPort:      1234,
			expectedOutput: []string{"1.1.1.1:1234", "2.2.2.2:1234", "3.3.3.3:1234"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := EndpointList(test.inputIpList, test.inputPort)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}

func TestEndpointListCommaSeparated(t *testing.T) {
	tests := []struct {
		name           string
		inputIpList    []string
		inputPort      int
		expectedOutput string
	}{
		{
			name:           "Should return empty string for empty list.",
			inputIpList:    []string{},
			inputPort:      1234,
			expectedOutput: "",
		},
		{
			name:           "Should return single endpoint for single IP, without commas",
			inputIpList:    []string{"1.1.1.1"},
			inputPort:      1234,
			expectedOutput: "1.1.1.1:1234",
		},
		{
			name:           "Should comma separated endpoints for multiple IPs",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			inputPort:      1234,
			expectedOutput: "1.1.1.1:1234,2.2.2.2:1234,3.3.3.3:1234",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := EndpointListCommaSeparated(test.inputIpList, test.inputPort)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}

func TestEndpointListCommaSeparatedQuoted(t *testing.T) {
	tests := []struct {
		name           string
		inputIpList    []string
		inputPort      int
		expectedOutput string
	}{
		{
			name:           "Should return empty string for empty list.",
			inputIpList:    []string{},
			inputPort:      1234,
			expectedOutput: "",
		},
		{
			name:           "Should return single endpoint for single IP, without commas, but quoted.",
			inputIpList:    []string{"1.1.1.1"},
			inputPort:      1234,
			expectedOutput: "'1.1.1.1:1234'",
		},
		{
			name:           "Should comma separated and quoted endpoints for multiple IPs",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			inputPort:      1234,
			expectedOutput: "'1.1.1.1:1234','2.2.2.2:1234','3.3.3.3:1234'",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := EndpointListCommaSeparatedQuoted(test.inputIpList, test.inputPort)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}

func TestEndpointListSpaceSeparated(t *testing.T) {
	tests := []struct {
		name           string
		inputIpList    []string
		inputPort      int
		expectedOutput string
	}{
		{
			name:           "Should return empty string for empty list.",
			inputIpList:    []string{},
			inputPort:      1234,
			expectedOutput: "",
		},
		{
			name:           "Should return single endpoint for single IP, without spaces",
			inputIpList:    []string{"1.1.1.1"},
			inputPort:      1234,
			expectedOutput: "1.1.1.1:1234",
		},
		{
			name:           "Should space separated endpoints for multiple IPs",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			inputPort:      1234,
			expectedOutput: "1.1.1.1:1234 2.2.2.2:1234 3.3.3.3:1234",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := EndpointListSpaceSeparated(test.inputIpList, test.inputPort)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}
