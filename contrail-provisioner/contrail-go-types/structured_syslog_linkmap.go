//
// Automatically generated. DO NOT EDIT.
//

package types

type StructuredSyslogLinkType struct {
	Overlay  string `json:"overlay,omitempty"`
	Underlay string `json:"underlay,omitempty"`
}

type StructuredSyslogLinkmap struct {
	Links []StructuredSyslogLinkType `json:"links,omitempty"`
}

func (obj *StructuredSyslogLinkmap) AddLinks(value *StructuredSyslogLinkType) {
	obj.Links = append(obj.Links, *value)
}
