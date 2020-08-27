//
// Automatically generated. DO NOT EDIT.
//

package types

type LoadbalancerType struct {
	Status             string `json:"status,omitempty"`
	ProvisioningStatus string `json:"provisioning_status,omitempty"`
	OperatingStatus    string `json:"operating_status,omitempty"`
	VipSubnetId        string `json:"vip_subnet_id,omitempty"`
	VipAddress         string `json:"vip_address,omitempty"`
	AdminState         bool   `json:"admin_state,omitempty"`
}
