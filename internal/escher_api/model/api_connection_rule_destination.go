package model

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type SecureCNConnectionRuleDestination struct {

	// Required: true
	Type *string `json:"type"` //ENUM

	IP *string `json:"ip"`
	//OR
	Address *string `json:"address"`
	//not added fields for external
	//OR
	Pod *SecureCNConnectionRulePod `json:"pod"`
}

// Validate validates this kubernetes cluster
func (m *SecureCNConnectionRuleDestination) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SecureCNConnectionRuleDestination) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SecureCNConnectionRuleDestination) UnmarshalBinary(b []byte) error {
	var res SecureCNConnectionRuleDestination
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
