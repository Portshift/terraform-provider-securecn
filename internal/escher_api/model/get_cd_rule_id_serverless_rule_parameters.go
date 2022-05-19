// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetCdRuleIDServerlessRuleParams creates a new GetCdRuleIDServerlessRuleParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetCdRuleIDServerlessRuleParams() *GetCdRuleIDServerlessRuleParams {
	return &GetCdRuleIDServerlessRuleParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetCdRuleIDServerlessRuleParamsWithTimeout creates a new GetCdRuleIDServerlessRuleParams object
// with the ability to set a timeout on a request.
func NewGetCdRuleIDServerlessRuleParamsWithTimeout(timeout time.Duration) *GetCdRuleIDServerlessRuleParams {
	return &GetCdRuleIDServerlessRuleParams{
		timeout: timeout,
	}
}

// NewGetCdRuleIDServerlessRuleParamsWithContext creates a new GetCdRuleIDServerlessRuleParams object
// with the ability to set a context for a request.
func NewGetCdRuleIDServerlessRuleParamsWithContext(ctx context.Context) *GetCdRuleIDServerlessRuleParams {
	return &GetCdRuleIDServerlessRuleParams{
		Context: ctx,
	}
}

// NewGetCdRuleIDServerlessRuleParamsWithHTTPClient creates a new GetCdRuleIDServerlessRuleParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetCdRuleIDServerlessRuleParamsWithHTTPClient(client *http.Client) *GetCdRuleIDServerlessRuleParams {
	return &GetCdRuleIDServerlessRuleParams{
		HTTPClient: client,
	}
}

/* GetCdRuleIDServerlessRuleParams contains all the parameters to send to the API endpoint
   for the get cd rule ID serverless rule operation.

   Typically these are written to a http.Request.
*/
type GetCdRuleIDServerlessRuleParams struct {

	/* RuleID.

	   ruleId (uid)

	   Format: uuid
	*/
	RuleID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get cd rule ID serverless rule params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetCdRuleIDServerlessRuleParams) WithDefaults() *GetCdRuleIDServerlessRuleParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get cd rule ID serverless rule params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetCdRuleIDServerlessRuleParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get cd rule ID serverless rule params
func (o *GetCdRuleIDServerlessRuleParams) WithTimeout(timeout time.Duration) *GetCdRuleIDServerlessRuleParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get cd rule ID serverless rule params
func (o *GetCdRuleIDServerlessRuleParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get cd rule ID serverless rule params
func (o *GetCdRuleIDServerlessRuleParams) WithContext(ctx context.Context) *GetCdRuleIDServerlessRuleParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get cd rule ID serverless rule params
func (o *GetCdRuleIDServerlessRuleParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get cd rule ID serverless rule params
func (o *GetCdRuleIDServerlessRuleParams) WithHTTPClient(client *http.Client) *GetCdRuleIDServerlessRuleParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get cd rule ID serverless rule params
func (o *GetCdRuleIDServerlessRuleParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRuleID adds the ruleID to the get cd rule ID serverless rule params
func (o *GetCdRuleIDServerlessRuleParams) WithRuleID(ruleID strfmt.UUID) *GetCdRuleIDServerlessRuleParams {
	o.SetRuleID(ruleID)
	return o
}

// SetRuleID adds the ruleId to the get cd rule ID serverless rule params
func (o *GetCdRuleIDServerlessRuleParams) SetRuleID(ruleID strfmt.UUID) {
	o.RuleID = ruleID
}

// WriteToRequest writes these params to a swagger request
func (o *GetCdRuleIDServerlessRuleParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param ruleId
	if err := r.SetPathParam("ruleId", o.RuleID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
