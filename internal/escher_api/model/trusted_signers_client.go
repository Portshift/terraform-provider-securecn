// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new trusted signers API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for trusted signers API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	DeleteTrustedSignersTrustedSignerID(params *DeleteTrustedSignersTrustedSignerIDParams, opts ...ClientOption) (*DeleteTrustedSignersTrustedSignerIDNoContent, error)

	GetTrustedSignersTrustedSignerID(params *GetTrustedSignersTrustedSignerIDParams, opts ...ClientOption) (*GetTrustedSignersTrustedSignerIDOK, error)

	PostTrustedSigners(params *PostTrustedSignersParams, opts ...ClientOption) (*PostTrustedSignersCreated, error)

	PutTrustedSignersTrustedSignerID(params *PutTrustedSignersTrustedSignerIDParams, opts ...ClientOption) (*PutTrustedSignersTrustedSignerIDCreated, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  DeleteTrustedSignersTrustedSignerID deletes a trusted signer
*/
func (a *Client) DeleteTrustedSignersTrustedSignerID(params *DeleteTrustedSignersTrustedSignerIDParams, opts ...ClientOption) (*DeleteTrustedSignersTrustedSignerIDNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteTrustedSignersTrustedSignerIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteTrustedSignersTrustedSignerID",
		Method:             "DELETE",
		PathPattern:        "/trustedSigners/{trustedSignerId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteTrustedSignersTrustedSignerIDReader{Formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteTrustedSignersTrustedSignerIDNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteTrustedSignersTrustedSignerID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetTrustedSignersTrustedSignerID gets existing trusted signer
*/
func (a *Client) GetTrustedSignersTrustedSignerID(params *GetTrustedSignersTrustedSignerIDParams, opts ...ClientOption) (*GetTrustedSignersTrustedSignerIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetTrustedSignersTrustedSignerIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetTrustedSignersTrustedSignerID",
		Method:             "GET",
		PathPattern:        "/trustedSigners/{trustedSignerId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetTrustedSignersTrustedSignerIDReader{Formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetTrustedSignersTrustedSignerIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetTrustedSignersTrustedSignerID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  PostTrustedSigners adds new trusted signer
*/
func (a *Client) PostTrustedSigners(params *PostTrustedSignersParams, opts ...ClientOption) (*PostTrustedSignersCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostTrustedSignersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostTrustedSigners",
		Method:             "POST",
		PathPattern:        "/trustedSigners",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostTrustedSignersReader{Formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostTrustedSignersCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostTrustedSigners: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  PutTrustedSignersTrustedSignerID edits existing trusted signer
*/
func (a *Client) PutTrustedSignersTrustedSignerID(params *PutTrustedSignersTrustedSignerIDParams, opts ...ClientOption) (*PutTrustedSignersTrustedSignerIDCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPutTrustedSignersTrustedSignerIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PutTrustedSignersTrustedSignerID",
		Method:             "PUT",
		PathPattern:        "/trustedSigners/{trustedSignerId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PutTrustedSignersTrustedSignerIDReader{Formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PutTrustedSignersTrustedSignerIDCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PutTrustedSignersTrustedSignerID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
