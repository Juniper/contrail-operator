//
// Automatically generated. DO NOT EDIT.
//

package types

type KeyValuePair struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type KeyValuePairs struct {
	KeyValuePair []KeyValuePair `json:"key_value_pair,omitempty"`
}

func (obj *KeyValuePairs) AddKeyValuePair(value *KeyValuePair) {
	obj.KeyValuePair = append(obj.KeyValuePair, *value)
}
