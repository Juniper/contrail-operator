//
// Automatically generated. DO NOT EDIT.
//

package types

type AlarmOperand2 struct {
	UveAttribute string `json:"uve_attribute,omitempty"`
	JsonValue string `json:"json_value,omitempty"`
}

type AlarmExpression struct {
	Operation string `json:"operation,omitempty"`
	Operand1 string `json:"operand1,omitempty"`
	Operand2 *AlarmOperand2 `json:"operand2,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

func (obj *AlarmExpression) AddVariables(value string) {
        obj.Variables = append(obj.Variables, value)
}

type AlarmAndList struct {
	AndList []AlarmExpression `json:"and_list,omitempty"`
}

func (obj *AlarmAndList) AddAndList(value *AlarmExpression) {
        obj.AndList = append(obj.AndList, *value)
}

type AlarmOrList struct {
	OrList []AlarmAndList `json:"or_list,omitempty"`
}

func (obj *AlarmOrList) AddOrList(value *AlarmAndList) {
        obj.OrList = append(obj.OrList, *value)
}
