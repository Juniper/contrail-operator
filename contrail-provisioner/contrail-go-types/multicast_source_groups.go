//
// Automatically generated. DO NOT EDIT.
//

package types

type MulticastSourceGroup struct {
	SourceAddress string `json:"source_address,omitempty"`
	GroupAddress string `json:"group_address,omitempty"`
	Action string `json:"action,omitempty"`
}

type MulticastSourceGroups struct {
	MulticastSourceGroup []MulticastSourceGroup `json:"multicast_source_group,omitempty"`
}

func (obj *MulticastSourceGroups) AddMulticastSourceGroup(value *MulticastSourceGroup) {
        obj.MulticastSourceGroup = append(obj.MulticastSourceGroup, *value)
}
