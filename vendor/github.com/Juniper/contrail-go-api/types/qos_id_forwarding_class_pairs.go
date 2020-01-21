//
// Automatically generated. DO NOT EDIT.
//

package types

type QosIdForwardingClassPair struct {
	Key int `json:"key,omitempty"`
	ForwardingClassId int `json:"forwarding_class_id,omitempty"`
}

type QosIdForwardingClassPairs struct {
	QosIdForwardingClassPair []QosIdForwardingClassPair `json:"qos_id_forwarding_class_pair,omitempty"`
}

func (obj *QosIdForwardingClassPairs) AddQosIdForwardingClassPair(value *QosIdForwardingClassPair) {
        obj.QosIdForwardingClassPair = append(obj.QosIdForwardingClassPair, *value)
}
