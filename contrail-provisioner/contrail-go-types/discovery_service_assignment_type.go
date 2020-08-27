//
// Automatically generated. DO NOT EDIT.
//

package types

type DiscoveryPubSubEndPointType struct {
	EpType    string      `json:"ep_type,omitempty"`
	EpId      string      `json:"ep_id,omitempty"`
	EpPrefix  *SubnetType `json:"ep_prefix,omitempty"`
	EpVersion string      `json:"ep_version,omitempty"`
}

type DiscoveryServiceAssignmentType struct {
	Publisher  *DiscoveryPubSubEndPointType  `json:"publisher,omitempty"`
	Subscriber []DiscoveryPubSubEndPointType `json:"subscriber,omitempty"`
}

func (obj *DiscoveryServiceAssignmentType) AddSubscriber(value *DiscoveryPubSubEndPointType) {
	obj.Subscriber = append(obj.Subscriber, *value)
}
