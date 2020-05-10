//
// Automatically generated. DO NOT EDIT.
//

package types

type FirewallRuleMatchTagsTypeIdList struct {
	TagType []int `json:"tag_type,omitempty"`
}

func (obj *FirewallRuleMatchTagsTypeIdList) AddTagType(value int) {
        obj.TagType = append(obj.TagType, value)
}
