//
// Automatically generated. DO NOT EDIT.
//

package types

type EncapsulationPrioritiesType struct {
	Encapsulation []string `json:"encapsulation,omitempty"`
}

func (obj *EncapsulationPrioritiesType) AddEncapsulation(value string) {
	obj.Encapsulation = append(obj.Encapsulation, value)
}
