//
// Automatically generated. DO NOT EDIT.
//

package types

type AddressFamilies struct {
	Family []string `json:"family,omitempty"`
}

func (obj *AddressFamilies) AddFamily(value string) {
        obj.Family = append(obj.Family, value)
}
