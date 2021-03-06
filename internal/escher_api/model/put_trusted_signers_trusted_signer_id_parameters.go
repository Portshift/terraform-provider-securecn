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

// NewPutTrustedSignersTrustedSignerIDParams creates a new PutTrustedSignersTrustedSignerIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPutTrustedSignersTrustedSignerIDParams() *PutTrustedSignersTrustedSignerIDParams {
	return &PutTrustedSignersTrustedSignerIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPutTrustedSignersTrustedSignerIDParamsWithTimeout creates a new PutTrustedSignersTrustedSignerIDParams object
// with the ability to set a timeout on a request.
func NewPutTrustedSignersTrustedSignerIDParamsWithTimeout(timeout time.Duration) *PutTrustedSignersTrustedSignerIDParams {
	return &PutTrustedSignersTrustedSignerIDParams{
		timeout: timeout,
	}
}

// NewPutTrustedSignersTrustedSignerIDParamsWithContext creates a new PutTrustedSignersTrustedSignerIDParams object
// with the ability to set a context for a request.
func NewPutTrustedSignersTrustedSignerIDParamsWithContext(ctx context.Context) *PutTrustedSignersTrustedSignerIDParams {
	return &PutTrustedSignersTrustedSignerIDParams{
		Context: ctx,
	}
}

// NewPutTrustedSignersTrustedSignerIDParamsWithHTTPClient creates a new PutTrustedSignersTrustedSignerIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewPutTrustedSignersTrustedSignerIDParamsWithHTTPClient(client *http.Client) *PutTrustedSignersTrustedSignerIDParams {
	return &PutTrustedSignersTrustedSignerIDParams{
		HTTPClient: client,
	}
}

/* PutTrustedSignersTrustedSignerIDParams contains all the parameters to send to the API endpoint
   for the put trusted signers trusted signer ID operation.

   Typically these are written to a http.Request.
*/
type PutTrustedSignersTrustedSignerIDParams struct {

	// Body.
	Body *TrustedSigner

	// TrustedSignerID.
	//
	// Format: uuid
	TrustedSignerID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the put trusted signers trusted signer ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutTrustedSignersTrustedSignerIDParams) WithDefaults() *PutTrustedSignersTrustedSignerIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the put trusted signers trusted signer ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PutTrustedSignersTrustedSignerIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) WithTimeout(timeout time.Duration) *PutTrustedSignersTrustedSignerIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) WithContext(ctx context.Context) *PutTrustedSignersTrustedSignerIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) WithHTTPClient(client *http.Client) *PutTrustedSignersTrustedSignerIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) WithBody(body *TrustedSigner) *PutTrustedSignersTrustedSignerIDParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) SetBody(body *TrustedSigner) {
	o.Body = body
}

// WithTrustedSignerID adds the trustedSignerID to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) WithTrustedSignerID(trustedSignerID strfmt.UUID) *PutTrustedSignersTrustedSignerIDParams {
	o.SetTrustedSignerID(trustedSignerID)
	return o
}

// SetTrustedSignerID adds the trustedSignerId to the put trusted signers trusted signer ID params
func (o *PutTrustedSignersTrustedSignerIDParams) SetTrustedSignerID(trustedSignerID strfmt.UUID) {
	o.TrustedSignerID = trustedSignerID
}

// WriteToRequest writes these params to a swagger request
func (o *PutTrustedSignersTrustedSignerIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param trustedSignerId
	if err := r.SetPathParam("trustedSignerId", o.TrustedSignerID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
