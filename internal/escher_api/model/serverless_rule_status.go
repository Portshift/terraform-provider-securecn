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

// ServerlessRuleStatus serverless rule status
//
// swagger:model ServerlessRuleStatus
type ServerlessRuleStatus string

func NewServerlessRuleStatus(value ServerlessRuleStatus) *ServerlessRuleStatus {
	v := value
	return &v
}

const (

	// ServerlessRuleStatusENABLED captures enum value "ENABLED"
	ServerlessRuleStatusENABLED ServerlessRuleStatus = "ENABLED"

	// ServerlessRuleStatusDISABLED captures enum value "DISABLED"
	ServerlessRuleStatusDISABLED ServerlessRuleStatus = "DISABLED"

	// ServerlessRuleStatusDELETED captures enum value "DELETED"
	ServerlessRuleStatusDELETED ServerlessRuleStatus = "DELETED"
)

// for schema
var serverlessRuleStatusEnum []interface{}

func init() {
	var res []ServerlessRuleStatus
	if err := json.Unmarshal([]byte(`["ENABLED","DISABLED","DELETED"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		serverlessRuleStatusEnum = append(serverlessRuleStatusEnum, v)
	}
}

func (m ServerlessRuleStatus) validateServerlessRuleStatusEnum(path, location string, value ServerlessRuleStatus) error {
	if err := validate.EnumCase(path, location, value, serverlessRuleStatusEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this serverless rule status
func (m ServerlessRuleStatus) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateServerlessRuleStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this serverless rule status based on context it is used
func (m ServerlessRuleStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
