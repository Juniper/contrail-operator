//
// Automatically generated. DO NOT EDIT.
//

package types

type GracefulRestartParametersType struct {
	Enable bool `json:"enable,omitempty"`
	RestartTime int `json:"restart_time,omitempty"`
	LongLivedRestartTime int `json:"long_lived_restart_time,omitempty"`
	EndOfRibTimeout int `json:"end_of_rib_timeout,omitempty"`
	BgpHelperEnable bool `json:"bgp_helper_enable,omitempty"`
	XmppHelperEnable bool `json:"xmpp_helper_enable,omitempty"`
}
