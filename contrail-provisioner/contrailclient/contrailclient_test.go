package contrailclient

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
)

func TestHasRequiredAnnotations(t *testing.T) {
	tests := []struct {
		name                string
		actualAnnotations   map[string]string
		requiredAnnotations map[string]string
		expectedOutput      bool
	}{
		{
			name: "should return true when required annotations are equal to actual annotations",
			actualAnnotations: map[string]string{
				"test": "annotation",
			},
			requiredAnnotations: map[string]string{
				"test": "annotation",
			},
			expectedOutput: true,
		},
		{
			name: "should return false when actual annotations have the same key but different value",
			actualAnnotations: map[string]string{
				"test": "different value",
			},
			requiredAnnotations: map[string]string{
				"test": "annotation",
			},
			expectedOutput: false,
		},
		{
			name: "should return true when required annotations is a subset of actual annotations",
			actualAnnotations: map[string]string{
				"test":  "annotation",
				"test2": "annotation2",
			},
			requiredAnnotations: map[string]string{
				"test": "annotation",
			},
			expectedOutput: true,
		},
		{
			name: "should return true when required annotations is an empty map",
			actualAnnotations: map[string]string{
				"test":  "annotation",
				"test2": "annotation2",
			},
			requiredAnnotations: map[string]string{},
			expectedOutput:      true,
		},
		{
			name:                "should return true when both required and actual annotations are an empty map",
			actualAnnotations:   map[string]string{},
			requiredAnnotations: map[string]string{},
			expectedOutput:      true,
		},
		{
			name:              "should return false when actual annotations is an empty map, but required annotations is not",
			actualAnnotations: map[string]string{},
			requiredAnnotations: map[string]string{
				"test":  "annotation",
				"test2": "annotation2",
			},
			expectedOutput: false,
		},
		{
			name: "should return false when one key value pair matches but the other does not",
			actualAnnotations: map[string]string{
				"test":  "annotation",
				"test2": "annotation3",
			},
			requiredAnnotations: map[string]string{
				"test":  "annotation",
				"test2": "annotation2",
			},
			expectedOutput: false,
		},
		{
			name:              "should false when required annotations have empty values and do not match the actual annotations",
			actualAnnotations: map[string]string{},
			requiredAnnotations: map[string]string{
				"test": "",
			},
			expectedOutput: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualResult := HasRequiredAnnotations(test.actualAnnotations, test.requiredAnnotations)
			assert.Equal(t, test.expectedOutput, actualResult)
		})
	}
}

func TestConvertContrailKeyValuePairsToMap(t *testing.T) {
	tests := []struct {
		name           string
		input          contrailtypes.KeyValuePairs
		expectedOutput map[string]string
	}{
		{
			name:           "empty input should result in empty output map",
			input:          contrailtypes.KeyValuePairs{},
			expectedOutput: map[string]string{},
		},
		{
			name: "multiple key value pairs should result in map with those pairs",
			input: contrailtypes.KeyValuePairs{
				KeyValuePair: []contrailtypes.KeyValuePair{
					{Key: "testkey1", Value: "testval1"},
					{Key: "testkey2", Value: "testval2"},
				},
			},
			expectedOutput: map[string]string{
				"testkey1": "testval1",
				"testkey2": "testval2",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualResult := ConvertContrailKeyValuePairsToMap(test.input)
			assert.True(t, reflect.DeepEqual(actualResult, test.expectedOutput), "%v not equal to %v", actualResult, test.expectedOutput)
		})
	}
}

func TestConvertMapToContrailKeyValuePairs(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]string
		expectedOutput contrailtypes.KeyValuePairs
	}{
		{
			name:  "empty input should result in KeyValuePairs with empty slice",
			input: map[string]string{},
			expectedOutput: contrailtypes.KeyValuePairs{
				KeyValuePair: []contrailtypes.KeyValuePair{},
			},
		},
		{
			name: "key value pairs should result in KeyValuePairs with those pairs",
			input: map[string]string{
				"testkey1": "testval1",
			},
			expectedOutput: contrailtypes.KeyValuePairs{
				KeyValuePair: []contrailtypes.KeyValuePair{
					{Key: "testkey1", Value: "testval1"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualResult := ConvertMapToContrailKeyValuePairs(test.input)
			assert.Equal(t, test.expectedOutput, actualResult)
		})
	}
}
