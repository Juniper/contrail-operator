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
			name:           "Should return empty string for empty list.",
			inputIpList:    []string{},
			expectedOutput: []string{},
		},
		{
			name:           "Shouldn't add commas to single ip and should add quotes",
			inputIpList:    []string{"1.1.1.1"},
			expectedOutput: []string{},
		},
		{
			name:           "Should add commas between multiple ips",
			inputIpList:    []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			expectedOutput: "'1.1.1.1','2.2.2.2','3.3.3.3'",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualOutput := EndpointList(test.inputIpList, test.inputPort)
			assert.Equal(t, test.expectedOutput, actualOutput)
		})
	}
}
