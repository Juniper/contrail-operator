//
// Automatically generated. DO NOT EDIT.
//

package types

type NamespaceValue struct {
	Ipv4Cidr   *SubnetListType        `json:"ipv4_cidr,omitempty"`
	Asn        *AutonomousSystemsType `json:"asn,omitempty"`
	MacAddr    *MacAddressesType      `json:"mac_addr,omitempty"`
	AsnRanges  []AsnRangeType         `json:"asn_ranges,omitempty"`
	SerialNums []string               `json:"serial_nums,omitempty"`
}

func (obj *NamespaceValue) AddAsnRanges(value *AsnRangeType) {
	obj.AsnRanges = append(obj.AsnRanges, *value)
}

func (obj *NamespaceValue) AddSerialNums(value string) {
	obj.SerialNums = append(obj.SerialNums, value)
}
