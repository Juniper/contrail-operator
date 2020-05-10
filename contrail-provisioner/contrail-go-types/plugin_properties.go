//
// Automatically generated. DO NOT EDIT.
//

package types

type PluginProperty struct {
	Property string `json:"property,omitempty"`
	Value string `json:"value,omitempty"`
}

type PluginProperties struct {
	PluginProperty []PluginProperty `json:"plugin_property,omitempty"`
}

func (obj *PluginProperties) AddPluginProperty(value *PluginProperty) {
        obj.PluginProperty = append(obj.PluginProperty, *value)
}
