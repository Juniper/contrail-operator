//
// Automatically generated. DO NOT EDIT.
//

package types

type FieldNamesList struct {
	FieldNames []string `json:"field_names,omitempty"`
}

func (obj *FieldNamesList) AddFieldNames(value string) {
	obj.FieldNames = append(obj.FieldNames, value)
}
