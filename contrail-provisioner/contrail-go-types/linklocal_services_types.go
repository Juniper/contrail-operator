//
// Automatically generated. DO NOT EDIT.
//

package types

type LinklocalServiceEntryType struct {
	LinklocalServiceName   string   `json:"linklocal_service_name,omitempty"`
	LinklocalServiceIp     string   `json:"linklocal_service_ip,omitempty"`
	LinklocalServicePort   int      `json:"linklocal_service_port,omitempty"`
	IpFabricDnsServiceName string   `json:"ip_fabric_DNS_service_name,omitempty"`
	IpFabricServicePort    int      `json:"ip_fabric_service_port,omitempty"`
	IpFabricServiceIp      []string `json:"ip_fabric_service_ip,omitempty"`
}

func (obj *LinklocalServiceEntryType) AddIpFabricServiceIp(value string) {
	obj.IpFabricServiceIp = append(obj.IpFabricServiceIp, value)
}

type LinklocalServicesTypes struct {
	LinklocalServiceEntry []LinklocalServiceEntryType `json:"linklocal_service_entry,omitempty"`
}

func (obj *LinklocalServicesTypes) AddLinklocalServiceEntry(value *LinklocalServiceEntryType) {
	obj.LinklocalServiceEntry = append(obj.LinklocalServiceEntry, *value)
}
