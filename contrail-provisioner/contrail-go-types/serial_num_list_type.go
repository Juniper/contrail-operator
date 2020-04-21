//
// Automatically generated. DO NOT EDIT.
//

package types

type SerialNumListType struct {
	SerialNum []string `json:"serial_num,omitempty"`
}

func (obj *SerialNumListType) AddSerialNum(value string) {
        obj.SerialNum = append(obj.SerialNum, value)
}
