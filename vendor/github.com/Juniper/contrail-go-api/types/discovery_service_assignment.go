//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	discovery_service_assignment_id_perms = iota
	discovery_service_assignment_perms2
	discovery_service_assignment_annotations
	discovery_service_assignment_display_name
	discovery_service_assignment_dsa_rules
	discovery_service_assignment_max_
)

type DiscoveryServiceAssignment struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	dsa_rules contrail.ReferenceList
        valid [discovery_service_assignment_max_] bool
        modified [discovery_service_assignment_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *DiscoveryServiceAssignment) GetType() string {
        return "discovery-service-assignment"
}

func (obj *DiscoveryServiceAssignment) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *DiscoveryServiceAssignment) GetDefaultParentType() string {
        return ""
}

func (obj *DiscoveryServiceAssignment) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *DiscoveryServiceAssignment) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *DiscoveryServiceAssignment) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *DiscoveryServiceAssignment) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *DiscoveryServiceAssignment) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *DiscoveryServiceAssignment) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *DiscoveryServiceAssignment) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[discovery_service_assignment_id_perms] = true
}

func (obj *DiscoveryServiceAssignment) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *DiscoveryServiceAssignment) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[discovery_service_assignment_perms2] = true
}

func (obj *DiscoveryServiceAssignment) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *DiscoveryServiceAssignment) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[discovery_service_assignment_annotations] = true
}

func (obj *DiscoveryServiceAssignment) GetDisplayName() string {
        return obj.display_name
}

func (obj *DiscoveryServiceAssignment) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[discovery_service_assignment_display_name] = true
}

func (obj *DiscoveryServiceAssignment) readDsaRules() error {
        if !obj.IsTransient() &&
                !obj.valid[discovery_service_assignment_dsa_rules] {
                err := obj.GetField(obj, "dsa_rules")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DiscoveryServiceAssignment) GetDsaRules() (
        contrail.ReferenceList, error) {
        err := obj.readDsaRules()
        if err != nil {
                return nil, err
        }
        return obj.dsa_rules, nil
}

func (obj *DiscoveryServiceAssignment) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[discovery_service_assignment_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[discovery_service_assignment_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[discovery_service_assignment_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[discovery_service_assignment_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *DiscoveryServiceAssignment) UnmarshalJSON(body []byte) error {
        var m map[string]json.RawMessage
        err := json.Unmarshal(body, &m)
        if err != nil {
                return err
        }
        err = obj.UnmarshalCommon(m)
        if err != nil {
                return err
        }
        for key, value := range m {
                switch key {
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[discovery_service_assignment_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[discovery_service_assignment_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[discovery_service_assignment_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[discovery_service_assignment_display_name] = true
                        }
                        break
                case "dsa_rules":
                        err = json.Unmarshal(value, &obj.dsa_rules)
                        if err == nil {
                                obj.valid[discovery_service_assignment_dsa_rules] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DiscoveryServiceAssignment) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[discovery_service_assignment_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[discovery_service_assignment_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[discovery_service_assignment_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[discovery_service_assignment_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *DiscoveryServiceAssignment) UpdateReferences() error {

        return nil
}

func DiscoveryServiceAssignmentByName(c contrail.ApiClient, fqn string) (*DiscoveryServiceAssignment, error) {
    obj, err := c.FindByName("discovery-service-assignment", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*DiscoveryServiceAssignment), nil
}

func DiscoveryServiceAssignmentByUuid(c contrail.ApiClient, uuid string) (*DiscoveryServiceAssignment, error) {
    obj, err := c.FindByUuid("discovery-service-assignment", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*DiscoveryServiceAssignment), nil
}
