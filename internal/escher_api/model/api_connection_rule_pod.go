package model

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type SecureCNConnectionRulePod struct {

	// Required: true
	Type *string `json:"type"` //ENUM

	// AND

	Name *string `json:"name"`
	// OR
	LabelKey   *string `json:"labelKey"`
	LabelValue *string `json:"labelValue"`
	// OR
	// no added fields for type 'any'

	// AND (optional)

	VulnerabilitySeverityLevel *string `json:"vulnerabilitySeverityLevel"` //ENUM

}

// Validate validates this kubernetes cluster
func (m *SecureCNConnectionRulePod) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SecureCNConnectionRulePod) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SecureCNConnectionRulePod) UnmarshalBinary(b []byte) error {
	var res SecureCNConnectionRulePod
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
