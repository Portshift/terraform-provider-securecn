// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// ServerlessRuleOrigin serverless rule origin
//
// swagger:model ServerlessRuleOrigin
type ServerlessRuleOrigin string

func NewServerlessRuleOrigin(value ServerlessRuleOrigin) *ServerlessRuleOrigin {
	v := value
	return &v
}

const (

	// ServerlessRuleOriginUSER captures enum value "USER"
	ServerlessRuleOriginUSER ServerlessRuleOrigin = "USER"

	// ServerlessRuleOriginAUTOMATEDPOLICY captures enum value "AUTOMATED_POLICY"
	ServerlessRuleOriginAUTOMATEDPOLICY ServerlessRuleOrigin = "AUTOMATED_POLICY"

	// ServerlessRuleOriginSYSTEM captures enum value "SYSTEM"
	ServerlessRuleOriginSYSTEM ServerlessRuleOrigin = "SYSTEM"
)

// for schema
var serverlessRuleOriginEnum []interface{}

func init() {
	var res []ServerlessRuleOrigin
	if err := json.Unmarshal([]byte(`["USER","AUTOMATED_POLICY","SYSTEM"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		serverlessRuleOriginEnum = append(serverlessRuleOriginEnum, v)
	}
}

func (m ServerlessRuleOrigin) validateServerlessRuleOriginEnum(path, location string, value ServerlessRuleOrigin) error {
	if err := validate.EnumCase(path, location, value, serverlessRuleOriginEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this serverless rule origin
func (m ServerlessRuleOrigin) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateServerlessRuleOriginEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this serverless rule origin based on context it is used
func (m ServerlessRuleOrigin) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
