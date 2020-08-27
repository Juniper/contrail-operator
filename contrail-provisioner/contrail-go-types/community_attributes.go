//
// Automatically generated. DO NOT EDIT.
//

package types

type CommunityAttributes struct {
	CommunityAttribute []string `json:"community_attribute,omitempty"`
}

func (obj *CommunityAttributes) AddCommunityAttribute(value string) {
	obj.CommunityAttribute = append(obj.CommunityAttribute, value)
}
