//
// Automatically generated. DO NOT EDIT.
//

package types

type BgpSession struct {
	Uuid       string                 `json:"uuid,omitempty"`
	Attributes []BgpSessionAttributes `json:"attributes,omitempty"`
}

func (obj *BgpSession) AddAttributes(value *BgpSessionAttributes) {
	obj.Attributes = append(obj.Attributes, *value)
}

type BgpPeeringAttributes struct {
	Session []BgpSession `json:"session,omitempty"`
}

func (obj *BgpPeeringAttributes) AddSession(value *BgpSession) {
	obj.Session = append(obj.Session, *value)
}
