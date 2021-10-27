package model

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

type SecureCNConnectionRule struct {
	// Id of the rule.
	// Read Only: true
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// Required: true
	Source *SecureCNConnectionRuleSource `json:"source"`

	// Required: true
	Destination *SecureCNConnectionRuleDestination `json:"destination"`
}

// Validate validates this kubernetes cluster
func (m *SecureCNConnectionRule) Validate(formats strfmt.Registry) error {
	var res []error
	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SecureCNConnectionRule) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SecureCNConnectionRule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SecureCNConnectionRule) UnmarshalBinary(b []byte) error {
	var res SecureCNConnectionRule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
