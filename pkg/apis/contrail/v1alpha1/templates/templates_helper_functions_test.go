package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinListWithSeparator(t *testing.T) {
	tests := []struct {
		name           string
		inputIpList    []string
		inputSeparator string
		expectedOutput string
	}{
		{
			name:           "Should return empty string for empty list.",
			inputIpList:    []string{},
			inputSeparator: ",",
			expectedOutput: "",
		},
		{
			name:           "Shouldn't add separators to single ip.",
			inputIpList:    []string{"1.1.1.1"},
			inputSeparator: ",",
			expectedOutput: "1.1.1.1",
		},
		{
			name:           "Should add separators between multiple ips",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			inputSeparator: ",",
			expectedOutput: "1.1.1.1,2.2.2.2,3.3.3.3",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := JoinListWithSeparator(test.inputIpList, test.inputSeparator)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}

func TestJoinListWithSeparatorAndSingleQuotes(t *testing.T) {
	tests := []struct {
		name           string
		inputIpList    []string
		inputSeparator string
		expectedOutput string
	}{
		{
			name:           "Should return empty string for empty list.",
			inputIpList:    []string{},
			inputSeparator: ",",
		},
		{
			name:           "Shouldn't add separators to single ip and should add quotes",
			inputIpList:    []string{"1.1.1.1"},
			inputSeparator: ",",
			expectedOutput: "'1.1.1.1'",
		},
		{
			name:           "Should add separators between multiple ips",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			inputSeparator: ",",
			expectedOutput: "'1.1.1.1','2.2.2.2','3.3.3.3'",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := JoinListWithSeparatorAndSingleQuotes(test.inputIpList, test.inputSeparator)
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
