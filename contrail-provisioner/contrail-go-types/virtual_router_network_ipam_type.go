//
// Automatically generated. DO NOT EDIT.
//

package types

type VirtualRouterNetworkIpamType struct {
	AllocationPools []AllocationPoolType `json:"allocation_pools,omitempty"`
	Subnet          []SubnetType         `json:"subnet,omitempty"`
}

func (obj *VirtualRouterNetworkIpamType) AddAllocationPools(value *AllocationPoolType) {
	obj.AllocationPools = append(obj.AllocationPools, *value)
}

func (obj *VirtualRouterNetworkIpamType) AddSubnet(value *SubnetType) {
	obj.Subnet = append(obj.Subnet, *value)
}
