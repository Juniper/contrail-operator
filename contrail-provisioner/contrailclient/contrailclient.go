package contrailclient

import (
	"errors"

	contrail "github.com/Juniper/contrail-go-api"

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
)

// ApiClient interface extends contrail.ApiClient by a missing ReadListResult
// to enable passing ApiClient interface instead of the struct to ease
// mocking in unit test
type ApiClient interface {
	contrail.ApiClient
	ReadListResult(string, *contrail.ListResult) (contrail.IObject, error)
}

func HasRequiredAnnotations(actualAnnotations map[string]string, requiredAnnotations map[string]string) bool {
	hasRequiredAnnotations := true
	for reqKey, regVal := range requiredAnnotations {
		if actualVal, ok := actualAnnotations[reqKey]; !ok || actualVal != regVal {
			hasRequiredAnnotations = false
		}
	}
	return hasRequiredAnnotations
}

func ConvertContrailKeyValuePairsToMap(keyValPairs contrailtypes.KeyValuePairs) map[string]string {
	output := map[string]string{}
	for _, annotation := range keyValPairs.KeyValuePair {
		output[annotation.Key] = annotation.Value
	}
	return output
}

func ConvertMapToContrailKeyValuePairs(keyValMap map[string]string) contrailtypes.KeyValuePairs {
	keyVals := []contrailtypes.KeyValuePair{}
	for key, val := range keyValMap {
		keyVals = append(keyVals, contrailtypes.KeyValuePair{Key: key, Value: val})
	}
	return contrailtypes.KeyValuePairs{KeyValuePair: keyVals}
}

func GetContrailObjectByName(contrailClient ApiClient, contrailType string, requiredName string) (contrail.IObject, error) {
	listResults, err := contrailClient.List(contrailType)
	if err != nil {
		return nil, err
	}
	for _, listResult := range listResults {
		obj, err := contrailClient.ReadListResult(contrailType, &listResult)
		if err != nil {
			return nil, err
		}
		if obj.GetName() == requiredName {
			return obj, nil
		}
	}
	return nil, errors.New(contrailType + " " + requiredName + " not found.")
}
