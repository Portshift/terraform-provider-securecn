package utils

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-securecn/internal/escher_api/model"
)

func ReadNestedListStringFromTF(d *schema.ResourceData, mainField string, subField string, index int) []string {
	ipsData, exists := d.GetOk(mainField)

	if exists {
		return readNestedListString(ipsData, subField, index)
	}

	return []string{}

}

func readNestedListString(data interface{}, subField string, index int) []string {
	interfaces := data.([]interface{})
	i := interfaces[index]
	imap := i.(map[string]interface{})
	values := make([]string, 0, len(imap))
	sub := imap[subField].([]interface{})
	for _, valueData := range sub {
		value := valueData.(string)
		values = append(values, value)
	}
	return FilterEmptyStrings(values)
}

func ReadNestedStringFromTF(d *schema.ResourceData, mainField string, subField string, index int) string {
	data, exists := d.GetOk(mainField)

	if exists {
		return readNestedString(data, subField, index)
	}

	return ""
}

func readNestedString(data interface{}, subField string, index int) string {

	interfaces := data.([]interface{})

	i := interfaces[index]
	if i == nil {
		return ""
	}
	imap := i.(map[string]interface{})

	sub := imap[subField].(interface{})
	value := sub.(string)
	return value
}

func ReadNestedBoolFromTF(d *schema.ResourceData, mainField string, subField string, index int) bool {
	data, exists := d.GetOk(mainField)

	if exists {
		return readNestedBool(data, subField, index)
	}

	return false
}

func readNestedBool(data interface{}, subField string, index int) bool {

	interfaces := data.([]interface{})

	i := interfaces[index]
	if i == nil {
		return false
	}
	imap := i.(map[string]interface{})

	sub := imap[subField].(interface{})
	value := sub.(bool)
	return value
}

func ReadNestedIntFromTF(d *schema.ResourceData, mainField string, subField string, index int) int {
	data, exists := d.GetOk(mainField)

	if exists {
		return readNestedInt(data, subField, index)
	}

	return 0
}

func readNestedInt(data interface{}, subField string, index int) int {

	interfaces := data.([]interface{})

	i := interfaces[index]
	if i == nil {
		return 0
	}
	imap := i.(map[string]interface{})

	sub := imap[subField].(interface{})
	value := sub.(int)
	return value
}

func ReadNestedMapStringFromTF(d *schema.ResourceData, mainField string, subField string, index int) map[string]string {
	data, exists := d.GetOk(mainField)

	if exists {
		return readNestedMapString(data, subField, index)
	}

	return nil
}

func readNestedMapString(data interface{}, subField string, index int) map[string]string {
	interfaces := data.([]interface{})
	inter := interfaces[index]
	imap := inter.(map[string]interface{})
	sub := imap[subField].(map[string]interface{})
	values := make(map[string]string, len(sub))

	for subkey, subvalue := range sub {
		values[subkey] = subvalue.(string)
	}

	return values
}

func GetListStringFromLabels(labelsMap []*model.Label) map[string]string {
	labels := make(map[string]string, 0)
	for _, label := range labelsMap {
		labels[label.Key] = label.Value
	}

	return labels
}

func GetLabelsFromMap(labelsMap map[string]string) []*model.Label {
	labels := make([]*model.Label, 0, len(labelsMap))
	for k, v := range labelsMap {
		label := &model.Label{
			Key:   k,
			Value: v,
		}
		labels = append(labels, label)
	}

	return FilterEmptyLabels(labels)
}

func FilterEmptyLabels(labels []*model.Label) []*model.Label {
	var ans []*model.Label
	for _, label := range labels {
		if label != nil {
			ans = append(ans, label)
		}
	}
	return ans
}
