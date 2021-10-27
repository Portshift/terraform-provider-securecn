package model

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type SecureCNConnectionRuleSource struct {

	// Required: true
	Type *string `json:"type"` //ENUM

	IP *string `json:"ip"`
	//OR
	//not added fields for external
	//OR
	Pod *SecureCNConnectionRulePod `json:"pod"`
}

// Validate validates this kubernetes cluster
func (m *SecureCNConnectionRuleSource) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SecureCNConnectionRuleSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SecureCNConnectionRuleSource) UnmarshalBinary(b []byte) error {
	var res SecureCNConnectionRuleSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
