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

// PodNameWorkloadRuleType pod name workload rule type
// swagger:model PodNameWorkloadRuleType
type PodNameWorkloadRuleType struct {

	// pods that match one of the given names
	Names []string `json:"names"`

	// pod validation
	PodValidation *PodValidation `json:"podValidation,omitempty"`
}

// WorkloadRuleType gets the workload rule type of this subtype
func (m *PodNameWorkloadRuleType) WorkloadRuleType() string {
	return "PodNameWorkloadRuleType"
}

// SetWorkloadRuleType sets the workload rule type of this subtype
func (m *PodNameWorkloadRuleType) SetWorkloadRuleType(val string) {

}

// Names gets the names of this subtype

// PodValidation gets the pod validation of this subtype

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *PodNameWorkloadRuleType) UnmarshalJSON(raw []byte) error {
	var data struct {

		// pods that match one of the given names
		Names []string `json:"names"`

		// pod validation
		PodValidation *PodValidation `json:"podValidation,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		WorkloadRuleType string `json:"workloadRuleType"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result PodNameWorkloadRuleType

	if base.WorkloadRuleType != result.WorkloadRuleType() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid workloadRuleType value: %q", base.WorkloadRuleType)
	}

	result.Names = data.Names

	result.PodValidation = data.PodValidation

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m PodNameWorkloadRuleType) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// pods that match one of the given names
		Names []string `json:"names"`

		// pod validation
		PodValidation *PodValidation `json:"podValidation,omitempty"`
	}{

		Names: m.Names,

		PodValidation: m.PodValidation,
	},
	)
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		WorkloadRuleType string `json:"workloadRuleType"`
	}{

		WorkloadRuleType: m.WorkloadRuleType(),
	},
	)
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this pod name workload rule type
func (m *PodNameWorkloadRuleType) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePodValidation(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PodNameWorkloadRuleType) validatePodValidation(formats strfmt.Registry) error {

	if swag.IsZero(m.PodValidation) { // not required
		return nil
	}

	if m.PodValidation != nil {
		if err := m.PodValidation.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("podValidation")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PodNameWorkloadRuleType) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PodNameWorkloadRuleType) UnmarshalBinary(b []byte) error {
	var res PodNameWorkloadRuleType
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}