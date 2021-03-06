// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// ConnectionRuleAction connection rule action
// swagger:model ConnectionRuleAction
type ConnectionRuleAction string

const (

	// ConnectionRuleActionDETECT captures enum value "DETECT"
	ConnectionRuleActionDETECT ConnectionRuleAction = "DETECT"

	// ConnectionRuleActionBLOCK captures enum value "BLOCK"
	ConnectionRuleActionBLOCK ConnectionRuleAction = "BLOCK"

	// ConnectionRuleActionALLOW captures enum value "ALLOW"
	ConnectionRuleActionALLOW ConnectionRuleAction = "ALLOW"

	// ConnectionRuleActionENCRYPT captures enum value "ENCRYPT"
	ConnectionRuleActionENCRYPT ConnectionRuleAction = "ENCRYPT"

	// ConnectionRuleActionENCRYPTDIRECT captures enum value "ENCRYPT_DIRECT"
	ConnectionRuleActionENCRYPTDIRECT ConnectionRuleAction = "ENCRYPT_DIRECT"
)

// for schema
var connectionRuleActionEnum []interface{}

func init() {
	var res []ConnectionRuleAction
	if err := json.Unmarshal([]byte(`["DETECT","BLOCK","ALLOW","ENCRYPT","ENCRYPT_DIRECT"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		connectionRuleActionEnum = append(connectionRuleActionEnum, v)
	}
}

func (m ConnectionRuleAction) validateConnectionRuleActionEnum(path, location string, value ConnectionRuleAction) error {
	if err := validate.Enum(path, location, value, connectionRuleActionEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this connection rule action
func (m ConnectionRuleAction) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateConnectionRuleActionEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
