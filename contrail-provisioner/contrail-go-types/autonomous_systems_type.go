//
// Automatically generated. DO NOT EDIT.
//

package types

type AutonomousSystemsType struct {
	Asn []int `json:"asn,omitempty"`
}

func (obj *AutonomousSystemsType) AddAsn(value int) {
	obj.Asn = append(obj.Asn, value)
}
