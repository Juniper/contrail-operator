//
// Automatically generated. DO NOT EDIT.
//

package types

type PortTranslationPool struct {
	Protocol string `json:"protocol,omitempty"`
	PortRange *PortType `json:"port_range,omitempty"`
	PortCount string `json:"port_count,omitempty"`
}

type PortTranslationPools struct {
	PortTranslationPool []PortTranslationPool `json:"port_translation_pool,omitempty"`
}

func (obj *PortTranslationPools) AddPortTranslationPool(value *PortTranslationPool) {
        obj.PortTranslationPool = append(obj.PortTranslationPool, *value)
}
