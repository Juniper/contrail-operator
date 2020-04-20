//
// Automatically generated. DO NOT EDIT.
//

package types

type BgpFamilyAttributes struct {
	AddressFamily string `json:"address_family,omitempty"`
	LoopCount int `json:"loop_count,omitempty"`
	PrefixLimit *BgpPrefixLimit `json:"prefix_limit,omitempty"`
	DefaultTunnelEncap []string `json:"default_tunnel_encap,omitempty"`
}

func (obj *BgpFamilyAttributes) AddDefaultTunnelEncap(value string) {
        obj.DefaultTunnelEncap = append(obj.DefaultTunnelEncap, value)
}
