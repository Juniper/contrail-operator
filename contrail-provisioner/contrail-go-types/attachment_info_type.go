//
// Automatically generated. DO NOT EDIT.
//

package types

type ProtocolStaticType struct {
	Route []string `json:"route,omitempty"`
}

func (obj *ProtocolStaticType) AddRoute(value string) {
	obj.Route = append(obj.Route, value)
}

type AttachmentInfoType struct {
	Static *ProtocolStaticType `json:"_static,omitempty"`
	Bgp    *ProtocolBgpType    `json:"bgp,omitempty"`
	Ospf   *ProtocolOspfType   `json:"ospf,omitempty"`
	State  string              `json:"state,omitempty"`
}
