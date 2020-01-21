//
// Automatically generated. DO NOT EDIT.
//

package types

type AuthenticationData struct {
	KeyType string `json:"key_type,omitempty"`
	KeyItems []AuthenticationKeyItem `json:"key_items,omitempty"`
}

func (obj *AuthenticationData) AddKeyItems(value *AuthenticationKeyItem) {
        obj.KeyItems = append(obj.KeyItems, *value)
}
