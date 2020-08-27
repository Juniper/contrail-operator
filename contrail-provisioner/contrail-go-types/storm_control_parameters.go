//
// Automatically generated. DO NOT EDIT.
//

package types

type StormControlParameters struct {
	StormControlActions     []string `json:"storm_control_actions,omitempty"`
	RecoveryTimeout         int      `json:"recovery_timeout,omitempty"`
	NoUnregisteredMulticast bool     `json:"no_unregistered_multicast,omitempty"`
	NoRegisteredMulticast   bool     `json:"no_registered_multicast,omitempty"`
	NoUnknownUnicast        bool     `json:"no_unknown_unicast,omitempty"`
	NoMulticast             bool     `json:"no_multicast,omitempty"`
	NoBroadcast             bool     `json:"no_broadcast,omitempty"`
	BandwidthPercent        int      `json:"bandwidth_percent,omitempty"`
}

func (obj *StormControlParameters) AddStormControlActions(value string) {
	obj.StormControlActions = append(obj.StormControlActions, value)
}
