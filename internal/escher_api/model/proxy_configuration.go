// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// ProxyConfiguration proxy configuration
// swagger:model ProxyConfiguration
type ProxyConfiguration struct {

	// Specifies if the proxy configuration should be used
	EnableProxy *bool `json:"enableProxy,omitempty"`

	// https proxy
	HTTPSProxy string `json:"httpsProxy,omitempty"`
}

// Validate validates this proxy configuration
func (m *ProxyConfiguration) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ProxyConfiguration) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProxyConfiguration) UnmarshalBinary(b []byte) error {
	var res ProxyConfiguration
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
