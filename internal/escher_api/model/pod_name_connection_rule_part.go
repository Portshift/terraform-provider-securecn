// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// PodNameConnectionRulePart pod name connection rule part
// swagger:model PodNameConnectionRulePart
type PodNameConnectionRulePart struct {

	// environments
	Environments []string `json:"environments"`

	// names
	Names []string `json:"names"`

	// vulnerability
	VulnerabilitySeverityLevel string `json:"vulnerabilitySeverityLevel,omitempty"`
}

// ConnectionRulePartType gets the connection rule part type of this subtype
func (m *PodNameConnectionRulePart) ConnectionRulePartType() string {
	return "PodNameConnectionRulePart"
}

// SetConnectionRulePartType sets the connection rule part type of this subtype
func (m *PodNameConnectionRulePart) SetConnectionRulePartType(val string) {

}

// Environments gets the environments of this subtype

// Names gets the names of this subtype

// Vulnerability gets the vulnerability of this subtype

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *PodNameConnectionRulePart) UnmarshalJSON(raw []byte) error {
	var data struct {

		// environments
		Environments []string `json:"environments"`

		// names
		Names []string `json:"names"`

		// vulnerability
		Vulnerability string `json:"vulnerabilitySeverityLevel,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		ConnectionRulePartType string `json:"connectionRulePartType"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result PodNameConnectionRulePart

	if base.ConnectionRulePartType != result.ConnectionRulePartType() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid connectionRulePartType value: %q", base.ConnectionRulePartType)
	}

	result.Environments = data.Environments

	result.Names = data.Names

	result.VulnerabilitySeverityLevel = data.Vulnerability

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m PodNameConnectionRulePart) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// environments
		Environments []string `json:"environments"`

		// names
		Names []string `json:"names"`

		// vulnerability
		Vulnerability string `json:"vulnerabilitySeverityLevel,omitempty"`
	}{

		Environments: m.Environments,

		Names: m.Names,

		Vulnerability: m.VulnerabilitySeverityLevel,
	},
	)
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		ConnectionRulePartType string `json:"connectionRulePartType"`
	}{

		ConnectionRulePartType: m.ConnectionRulePartType(),
	},
	)
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this pod name connection rule part
func (m *PodNameConnectionRulePart) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateVulnerability(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PodNameConnectionRulePart) validateVulnerability(formats strfmt.Registry) error {

	if swag.IsZero(m.VulnerabilitySeverityLevel) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PodNameConnectionRulePart) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PodNameConnectionRulePart) UnmarshalBinary(b []byte) error {
	var res PodNameConnectionRulePart
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

//func (s *PodNameConnectionRulePart) Equal(raw interface{}) bool {
//	return false
//	other, ok := raw.(*PodNameConnectionRulePart)
//	if !ok {
//		return false
//	}
//
//	return reflect.DeepEqual(s.Names, other.Names) && reflect.DeepEqual(s.Names, other.Names) && reflect.DeepEqual(s.Environments, other.Environments) && reflect.DeepEqual(s.VulnerabilitySeverityLevel, other.VulnerabilitySeverityLevel)
//}
