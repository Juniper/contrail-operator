//
// Automatically generated. DO NOT EDIT.
//

package types

type FirewallRuleEndpointType struct {
	Subnet *SubnetType `json:"subnet,omitempty"`
	VirtualNetwork string `json:"virtual_network,omitempty"`
	AddressGroup string `json:"address_group,omitempty"`
	Tags []string `json:"tags,omitempty"`
	TagIds []int `json:"tag_ids,omitempty"`
	Any bool `json:"any,omitempty"`
}

func (obj *FirewallRuleEndpointType) AddTags(value string) {
        obj.Tags = append(obj.Tags, value)
}

func (obj *FirewallRuleEndpointType) AddTagIds(value int) {
        obj.TagIds = append(obj.TagIds, value)
}
