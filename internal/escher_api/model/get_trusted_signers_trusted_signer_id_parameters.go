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

// NewGetTrustedSignersTrustedSignerIDParams creates a new GetTrustedSignersTrustedSignerIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetTrustedSignersTrustedSignerIDParams() *GetTrustedSignersTrustedSignerIDParams {
	return &GetTrustedSignersTrustedSignerIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetTrustedSignersTrustedSignerIDParamsWithTimeout creates a new GetTrustedSignersTrustedSignerIDParams object
// with the ability to set a timeout on a request.
func NewGetTrustedSignersTrustedSignerIDParamsWithTimeout(timeout time.Duration) *GetTrustedSignersTrustedSignerIDParams {
	return &GetTrustedSignersTrustedSignerIDParams{
		timeout: timeout,
	}
}

// NewGetTrustedSignersTrustedSignerIDParamsWithContext creates a new GetTrustedSignersTrustedSignerIDParams object
// with the ability to set a context for a request.
func NewGetTrustedSignersTrustedSignerIDParamsWithContext(ctx context.Context) *GetTrustedSignersTrustedSignerIDParams {
	return &GetTrustedSignersTrustedSignerIDParams{
		Context: ctx,
	}
}

// NewGetTrustedSignersTrustedSignerIDParamsWithHTTPClient creates a new GetTrustedSignersTrustedSignerIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetTrustedSignersTrustedSignerIDParamsWithHTTPClient(client *http.Client) *GetTrustedSignersTrustedSignerIDParams {
	return &GetTrustedSignersTrustedSignerIDParams{
		HTTPClient: client,
	}
}

/* GetTrustedSignersTrustedSignerIDParams contains all the parameters to send to the API endpoint
   for the get trusted signers trusted signer ID operation.

   Typically these are written to a http.Request.
*/
type GetTrustedSignersTrustedSignerIDParams struct {

	// TrustedSignerID.
	//
	// Format: uuid
	TrustedSignerID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get trusted signers trusted signer ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetTrustedSignersTrustedSignerIDParams) WithDefaults() *GetTrustedSignersTrustedSignerIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get trusted signers trusted signer ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetTrustedSignersTrustedSignerIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get trusted signers trusted signer ID params
func (o *GetTrustedSignersTrustedSignerIDParams) WithTimeout(timeout time.Duration) *GetTrustedSignersTrustedSignerIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get trusted signers trusted signer ID params
func (o *GetTrustedSignersTrustedSignerIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get trusted signers trusted signer ID params
func (o *GetTrustedSignersTrustedSignerIDParams) WithContext(ctx context.Context) *GetTrustedSignersTrustedSignerIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get trusted signers trusted signer ID params
func (o *GetTrustedSignersTrustedSignerIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get trusted signers trusted signer ID params
func (o *GetTrustedSignersTrustedSignerIDParams) WithHTTPClient(client *http.Client) *GetTrustedSignersTrustedSignerIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get trusted signers trusted signer ID params
func (o *GetTrustedSignersTrustedSignerIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithTrustedSignerID adds the trustedSignerID to the get trusted signers trusted signer ID params
func (o *GetTrustedSignersTrustedSignerIDParams) WithTrustedSignerID(trustedSignerID strfmt.UUID) *GetTrustedSignersTrustedSignerIDParams {
	o.SetTrustedSignerID(trustedSignerID)
	return o
}

// SetTrustedSignerID adds the trustedSignerId to the get trusted signers trusted signer ID params
func (o *GetTrustedSignersTrustedSignerIDParams) SetTrustedSignerID(trustedSignerID strfmt.UUID) {
	o.TrustedSignerID = trustedSignerID
}

// WriteToRequest writes these params to a swagger request
func (o *GetTrustedSignersTrustedSignerIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param trustedSignerId
	if err := r.SetPathParam("trustedSignerId", o.TrustedSignerID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
