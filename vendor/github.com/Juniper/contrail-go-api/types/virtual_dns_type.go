//
// Automatically generated. DO NOT EDIT.
//

package types

type DnsSoaRecordType struct {
	NegativeCacheTtlSeconds int `json:"negative_cache_ttl_seconds,omitempty"`
}

type VirtualDnsType struct {
	DomainName string `json:"domain_name,omitempty"`
	DynamicRecordsFromClient bool `json:"dynamic_records_from_client,omitempty"`
	RecordOrder string `json:"record_order,omitempty"`
	DefaultTtlSeconds int `json:"default_ttl_seconds,omitempty"`
	NextVirtualDns string `json:"next_virtual_DNS,omitempty"`
	FloatingIpRecord string `json:"floating_ip_record,omitempty"`
	ExternalVisible bool `json:"external_visible,omitempty"`
	ReverseResolution bool `json:"reverse_resolution,omitempty"`
	SoaRecord *DnsSoaRecordType `json:"soa_record,omitempty"`
}
