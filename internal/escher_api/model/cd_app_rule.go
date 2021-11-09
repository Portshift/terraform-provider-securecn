// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"
	"io"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CdAppRule A rule that states what Apps are allowed to run on what environments.
// swagger:model CdAppRule
type CdAppRule struct {

	// action
	// Required: true
	Action AppRuleType `json:"action"`

	appField WorkloadRuleType

	// group name
	GroupName string `json:"groupName,omitempty"`

	// id
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// name
	// Required: true
	Name *string `json:"name"`

	// A way to identify the rule scope. Only one of the below should be not null, and used.
	Scope WorkloadRuleScopeType `json:"scope,omitempty"`

	// status
	// Required: true
	Status AppRuleStatus `json:"status"`
}

// App gets the app of this base type
func (m *CdAppRule) App() WorkloadRuleType {
	return m.appField
}

// SetApp sets the app of this base type
func (m *CdAppRule) SetApp(val WorkloadRuleType) {
	m.appField = val
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *CdAppRule) UnmarshalJSON(raw []byte) error {
	var data struct {
		Action AppRuleType `json:"action"`

		App json.RawMessage `json:"app,omitempty"`

		GroupName string `json:"groupName,omitempty"`

		ID strfmt.UUID `json:"id,omitempty"`

		Name *string `json:"name"`

		Scope WorkloadRuleScopeType `json:"scope,omitempty"`

		Status AppRuleStatus `json:"status"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var propApp WorkloadRuleType
	if string(data.App) != "null" {
		app, err := UnmarshalWorkloadRuleType(bytes.NewBuffer(data.App), runtime.JSONConsumer())
		if err != nil && err != io.EOF {
			return err
		}
		propApp = app
	}

	var result CdAppRule

	// action
	result.Action = data.Action

	// app
	result.appField = propApp

	// groupName
	result.GroupName = data.GroupName

	// id
	result.ID = data.ID

	// name
	result.Name = data.Name

	// scope
	result.Scope = data.Scope

	// status
	result.Status = data.Status

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m CdAppRule) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {
		Action AppRuleType `json:"action"`

		GroupName string `json:"groupName,omitempty"`

		ID strfmt.UUID `json:"id,omitempty"`

		Name *string `json:"name"`

		Scope WorkloadRuleScopeType `json:"scope,omitempty"`

		Status AppRuleStatus `json:"status"`
	}{

		Action: m.Action,

		GroupName: m.GroupName,

		ID: m.ID,

		Name: m.Name,

		Scope: m.Scope,

		Status: m.Status,
	},
	)
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		App WorkloadRuleType `json:"app,omitempty"`
	}{

		App: m.appField,
	},
	)
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this cd app rule
func (m *CdAppRule) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAction(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateApp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateScope(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CdAppRule) validateAction(formats strfmt.Registry) error {

	if err := m.Action.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("action")
		}
		return err
	}

	return nil
}

func (m *CdAppRule) validateApp(formats strfmt.Registry) error {

	if swag.IsZero(m.App()) { // not required
		return nil
	}

	if err := m.App().Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("app")
		}
		return err
	}

	return nil
}

func (m *CdAppRule) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *CdAppRule) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *CdAppRule) validateScope(formats strfmt.Registry) error {

	if swag.IsZero(m.Scope) { // not required
		return nil
	}

	if err := m.Scope.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("scope")
		}
		return err
	}

	return nil
}

func (m *CdAppRule) validateStatus(formats strfmt.Registry) error {

	if err := m.Status.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("status")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CdAppRule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CdAppRule) UnmarshalBinary(b []byte) error {
	var res CdAppRule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}