//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	device_image_device_image_vendor_name = iota
	device_image_device_image_device_family
	device_image_device_image_supported_platforms
	device_image_device_image_os_version
	device_image_device_image_file_uri
	device_image_device_image_size
	device_image_device_image_md5
	device_image_device_image_sha1
	device_image_id_perms
	device_image_perms2
	device_image_annotations
	device_image_display_name
	device_image_hardware_refs
	device_image_tag_refs
	device_image_physical_router_back_refs
	device_image_max_
)

type DeviceImage struct {
        contrail.ObjectBase
	device_image_vendor_name string
	device_image_device_family string
	device_image_supported_platforms DevicePlatformListType
	device_image_os_version string
	device_image_file_uri string
	device_image_size int
	device_image_md5 string
	device_image_sha1 string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	hardware_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
        valid [device_image_max_] bool
        modified [device_image_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *DeviceImage) GetType() string {
        return "device-image"
}

func (obj *DeviceImage) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *DeviceImage) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *DeviceImage) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *DeviceImage) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *DeviceImage) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *DeviceImage) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *DeviceImage) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *DeviceImage) GetDeviceImageVendorName() string {
        return obj.device_image_vendor_name
}

func (obj *DeviceImage) SetDeviceImageVendorName(value string) {
        obj.device_image_vendor_name = value
        obj.modified[device_image_device_image_vendor_name] = true
}

func (obj *DeviceImage) GetDeviceImageDeviceFamily() string {
        return obj.device_image_device_family
}

func (obj *DeviceImage) SetDeviceImageDeviceFamily(value string) {
        obj.device_image_device_family = value
        obj.modified[device_image_device_image_device_family] = true
}

func (obj *DeviceImage) GetDeviceImageSupportedPlatforms() DevicePlatformListType {
        return obj.device_image_supported_platforms
}

func (obj *DeviceImage) SetDeviceImageSupportedPlatforms(value *DevicePlatformListType) {
        obj.device_image_supported_platforms = *value
        obj.modified[device_image_device_image_supported_platforms] = true
}

func (obj *DeviceImage) GetDeviceImageOsVersion() string {
        return obj.device_image_os_version
}

func (obj *DeviceImage) SetDeviceImageOsVersion(value string) {
        obj.device_image_os_version = value
        obj.modified[device_image_device_image_os_version] = true
}

func (obj *DeviceImage) GetDeviceImageFileUri() string {
        return obj.device_image_file_uri
}

func (obj *DeviceImage) SetDeviceImageFileUri(value string) {
        obj.device_image_file_uri = value
        obj.modified[device_image_device_image_file_uri] = true
}

func (obj *DeviceImage) GetDeviceImageSize() int {
        return obj.device_image_size
}

func (obj *DeviceImage) SetDeviceImageSize(value int) {
        obj.device_image_size = value
        obj.modified[device_image_device_image_size] = true
}

func (obj *DeviceImage) GetDeviceImageMd5() string {
        return obj.device_image_md5
}

func (obj *DeviceImage) SetDeviceImageMd5(value string) {
        obj.device_image_md5 = value
        obj.modified[device_image_device_image_md5] = true
}

func (obj *DeviceImage) GetDeviceImageSha1() string {
        return obj.device_image_sha1
}

func (obj *DeviceImage) SetDeviceImageSha1(value string) {
        obj.device_image_sha1 = value
        obj.modified[device_image_device_image_sha1] = true
}

func (obj *DeviceImage) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *DeviceImage) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[device_image_id_perms] = true
}

func (obj *DeviceImage) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *DeviceImage) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[device_image_perms2] = true
}

func (obj *DeviceImage) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *DeviceImage) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[device_image_annotations] = true
}

func (obj *DeviceImage) GetDisplayName() string {
        return obj.display_name
}

func (obj *DeviceImage) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[device_image_display_name] = true
}

func (obj *DeviceImage) readHardwareRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[device_image_hardware_refs] {
                err := obj.GetField(obj, "hardware_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DeviceImage) GetHardwareRefs() (
        contrail.ReferenceList, error) {
        err := obj.readHardwareRefs()
        if err != nil {
                return nil, err
        }
        return obj.hardware_refs, nil
}

func (obj *DeviceImage) AddHardware(
        rhs *Hardware) error {
        err := obj.readHardwareRefs()
        if err != nil {
                return err
        }

        if !obj.modified[device_image_hardware_refs] {
                obj.storeReferenceBase("hardware", obj.hardware_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.hardware_refs = append(obj.hardware_refs, ref)
        obj.modified[device_image_hardware_refs] = true
        return nil
}

func (obj *DeviceImage) DeleteHardware(uuid string) error {
        err := obj.readHardwareRefs()
        if err != nil {
                return err
        }

        if !obj.modified[device_image_hardware_refs] {
                obj.storeReferenceBase("hardware", obj.hardware_refs)
        }

        for i, ref := range obj.hardware_refs {
                if ref.Uuid == uuid {
                        obj.hardware_refs = append(
                                obj.hardware_refs[:i],
                                obj.hardware_refs[i+1:]...)
                        break
                }
        }
        obj.modified[device_image_hardware_refs] = true
        return nil
}

func (obj *DeviceImage) ClearHardware() {
        if obj.valid[device_image_hardware_refs] &&
           !obj.modified[device_image_hardware_refs] {
                obj.storeReferenceBase("hardware", obj.hardware_refs)
        }
        obj.hardware_refs = make([]contrail.Reference, 0)
        obj.valid[device_image_hardware_refs] = true
        obj.modified[device_image_hardware_refs] = true
}

func (obj *DeviceImage) SetHardwareList(
        refList []contrail.ReferencePair) {
        obj.ClearHardware()
        obj.hardware_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.hardware_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *DeviceImage) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[device_image_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DeviceImage) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *DeviceImage) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[device_image_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[device_image_tag_refs] = true
        return nil
}

func (obj *DeviceImage) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[device_image_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        for i, ref := range obj.tag_refs {
                if ref.Uuid == uuid {
                        obj.tag_refs = append(
                                obj.tag_refs[:i],
                                obj.tag_refs[i+1:]...)
                        break
                }
        }
        obj.modified[device_image_tag_refs] = true
        return nil
}

func (obj *DeviceImage) ClearTag() {
        if obj.valid[device_image_tag_refs] &&
           !obj.modified[device_image_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[device_image_tag_refs] = true
        obj.modified[device_image_tag_refs] = true
}

func (obj *DeviceImage) SetTagList(
        refList []contrail.ReferencePair) {
        obj.ClearTag()
        obj.tag_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.tag_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *DeviceImage) readPhysicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[device_image_physical_router_back_refs] {
                err := obj.GetField(obj, "physical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DeviceImage) GetPhysicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_back_refs, nil
}

func (obj *DeviceImage) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[device_image_device_image_vendor_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_vendor_name)
                if err != nil {
                        return nil, err
                }
                msg["device_image_vendor_name"] = &value
        }

        if obj.modified[device_image_device_image_device_family] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_device_family)
                if err != nil {
                        return nil, err
                }
                msg["device_image_device_family"] = &value
        }

        if obj.modified[device_image_device_image_supported_platforms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_supported_platforms)
                if err != nil {
                        return nil, err
                }
                msg["device_image_supported_platforms"] = &value
        }

        if obj.modified[device_image_device_image_os_version] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_os_version)
                if err != nil {
                        return nil, err
                }
                msg["device_image_os_version"] = &value
        }

        if obj.modified[device_image_device_image_file_uri] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_file_uri)
                if err != nil {
                        return nil, err
                }
                msg["device_image_file_uri"] = &value
        }

        if obj.modified[device_image_device_image_size] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_size)
                if err != nil {
                        return nil, err
                }
                msg["device_image_size"] = &value
        }

        if obj.modified[device_image_device_image_md5] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_md5)
                if err != nil {
                        return nil, err
                }
                msg["device_image_md5"] = &value
        }

        if obj.modified[device_image_device_image_sha1] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_sha1)
                if err != nil {
                        return nil, err
                }
                msg["device_image_sha1"] = &value
        }

        if obj.modified[device_image_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[device_image_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[device_image_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[device_image_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.hardware_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.hardware_refs)
                if err != nil {
                        return nil, err
                }
                msg["hardware_refs"] = &value
        }

        if len(obj.tag_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.tag_refs)
                if err != nil {
                        return nil, err
                }
                msg["tag_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *DeviceImage) UnmarshalJSON(body []byte) error {
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
                case "device_image_vendor_name":
                        err = json.Unmarshal(value, &obj.device_image_vendor_name)
                        if err == nil {
                                obj.valid[device_image_device_image_vendor_name] = true
                        }
                        break
                case "device_image_device_family":
                        err = json.Unmarshal(value, &obj.device_image_device_family)
                        if err == nil {
                                obj.valid[device_image_device_image_device_family] = true
                        }
                        break
                case "device_image_supported_platforms":
                        err = json.Unmarshal(value, &obj.device_image_supported_platforms)
                        if err == nil {
                                obj.valid[device_image_device_image_supported_platforms] = true
                        }
                        break
                case "device_image_os_version":
                        err = json.Unmarshal(value, &obj.device_image_os_version)
                        if err == nil {
                                obj.valid[device_image_device_image_os_version] = true
                        }
                        break
                case "device_image_file_uri":
                        err = json.Unmarshal(value, &obj.device_image_file_uri)
                        if err == nil {
                                obj.valid[device_image_device_image_file_uri] = true
                        }
                        break
                case "device_image_size":
                        err = json.Unmarshal(value, &obj.device_image_size)
                        if err == nil {
                                obj.valid[device_image_device_image_size] = true
                        }
                        break
                case "device_image_md5":
                        err = json.Unmarshal(value, &obj.device_image_md5)
                        if err == nil {
                                obj.valid[device_image_device_image_md5] = true
                        }
                        break
                case "device_image_sha1":
                        err = json.Unmarshal(value, &obj.device_image_sha1)
                        if err == nil {
                                obj.valid[device_image_device_image_sha1] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[device_image_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[device_image_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[device_image_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[device_image_display_name] = true
                        }
                        break
                case "hardware_refs":
                        err = json.Unmarshal(value, &obj.hardware_refs)
                        if err == nil {
                                obj.valid[device_image_hardware_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[device_image_tag_refs] = true
                        }
                        break
                case "physical_router_back_refs":
                        err = json.Unmarshal(value, &obj.physical_router_back_refs)
                        if err == nil {
                                obj.valid[device_image_physical_router_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DeviceImage) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[device_image_device_image_vendor_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_vendor_name)
                if err != nil {
                        return nil, err
                }
                msg["device_image_vendor_name"] = &value
        }

        if obj.modified[device_image_device_image_device_family] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_device_family)
                if err != nil {
                        return nil, err
                }
                msg["device_image_device_family"] = &value
        }

        if obj.modified[device_image_device_image_supported_platforms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_supported_platforms)
                if err != nil {
                        return nil, err
                }
                msg["device_image_supported_platforms"] = &value
        }

        if obj.modified[device_image_device_image_os_version] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_os_version)
                if err != nil {
                        return nil, err
                }
                msg["device_image_os_version"] = &value
        }

        if obj.modified[device_image_device_image_file_uri] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_file_uri)
                if err != nil {
                        return nil, err
                }
                msg["device_image_file_uri"] = &value
        }

        if obj.modified[device_image_device_image_size] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_size)
                if err != nil {
                        return nil, err
                }
                msg["device_image_size"] = &value
        }

        if obj.modified[device_image_device_image_md5] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_md5)
                if err != nil {
                        return nil, err
                }
                msg["device_image_md5"] = &value
        }

        if obj.modified[device_image_device_image_sha1] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.device_image_sha1)
                if err != nil {
                        return nil, err
                }
                msg["device_image_sha1"] = &value
        }

        if obj.modified[device_image_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[device_image_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[device_image_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[device_image_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[device_image_hardware_refs] {
                if len(obj.hardware_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["hardware_refs"] = &value
                } else if !obj.hasReferenceBase("hardware") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.hardware_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["hardware_refs"] = &value
                }
        }


        if obj.modified[device_image_tag_refs] {
                if len(obj.tag_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["tag_refs"] = &value
                } else if !obj.hasReferenceBase("tag") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.tag_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["tag_refs"] = &value
                }
        }


        return json.Marshal(msg)
}

func (obj *DeviceImage) UpdateReferences() error {

        if obj.modified[device_image_hardware_refs] &&
           len(obj.hardware_refs) > 0 &&
           obj.hasReferenceBase("hardware") {
                err := obj.UpdateReference(
                        obj, "hardware",
                        obj.hardware_refs,
                        obj.baseMap["hardware"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[device_image_tag_refs] &&
           len(obj.tag_refs) > 0 &&
           obj.hasReferenceBase("tag") {
                err := obj.UpdateReference(
                        obj, "tag",
                        obj.tag_refs,
                        obj.baseMap["tag"])
                if err != nil {
                        return err
                }
        }

        return nil
}

func DeviceImageByName(c contrail.ApiClient, fqn string) (*DeviceImage, error) {
    obj, err := c.FindByName("device-image", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*DeviceImage), nil
}

func DeviceImageByUuid(c contrail.ApiClient, uuid string) (*DeviceImage, error) {
    obj, err := c.FindByUuid("device-image", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*DeviceImage), nil
}
