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

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetDeployersParams creates a new GetDeployersParams object
// with the default values initialized.
func NewGetDeployersParams() *GetDeployersParams {

	return &GetDeployersParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetDeployersParamsWithTimeout creates a new GetDeployersParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetDeployersParamsWithTimeout(timeout time.Duration) *GetDeployersParams {

	return &GetDeployersParams{

		timeout: timeout,
	}
}

// NewGetDeployersParamsWithContext creates a new GetDeployersParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetDeployersParamsWithContext(ctx context.Context) *GetDeployersParams {

	return &GetDeployersParams{

		Context: ctx,
	}
}

// NewGetDeployersParamsWithHTTPClient creates a new GetDeployersParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetDeployersParamsWithHTTPClient(client *http.Client) *GetDeployersParams {

	return &GetDeployersParams{
		HTTPClient: client,
	}
}

/*GetDeployersParams contains all the parameters to send to the API endpoint
for the get deployers operation typically these are written to a http.Request
*/
type GetDeployersParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get deployers params
func (o *GetDeployersParams) WithTimeout(timeout time.Duration) *GetDeployersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get deployers params
func (o *GetDeployersParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get deployers params
func (o *GetDeployersParams) WithContext(ctx context.Context) *GetDeployersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get deployers params
func (o *GetDeployersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get deployers params
func (o *GetDeployersParams) WithHTTPClient(client *http.Client) *GetDeployersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get deployers params
func (o *GetDeployersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetDeployersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
