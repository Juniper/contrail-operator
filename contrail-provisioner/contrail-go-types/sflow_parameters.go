//
// Automatically generated. DO NOT EDIT.
//

package types

type StatsCollectionFrequency struct {
	SampleRate      int    `json:"sample_rate,omitempty"`
	PollingInterval int    `json:"polling_interval,omitempty"`
	Direction       string `json:"direction,omitempty"`
}

type EnabledInterfaceParams struct {
	Name                     string                    `json:"name,omitempty"`
	StatsCollectionFrequency *StatsCollectionFrequency `json:"stats_collection_frequency,omitempty"`
}

type SflowParameters struct {
	StatsCollectionFrequency *StatsCollectionFrequency `json:"stats_collection_frequency,omitempty"`
	AgentId                  string                    `json:"agent_id,omitempty"`
	AdaptiveSampleRate       int                       `json:"adaptive_sample_rate,omitempty"`
	EnabledInterfaceType     string                    `json:"enabled_interface_type,omitempty"`
	EnabledInterfaceParams   []EnabledInterfaceParams  `json:"enabled_interface_params,omitempty"`
}

func (obj *SflowParameters) AddEnabledInterfaceParams(value *EnabledInterfaceParams) {
	obj.EnabledInterfaceParams = append(obj.EnabledInterfaceParams, *value)
}
