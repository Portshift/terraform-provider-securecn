package model

import (
	"context"
	"fmt"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"io"
	"net/http"
	"time"
)

type PostKubernetesClustersDefault struct {
	_statusCode int

	Payload *APIResponse
}

// Code gets the status code for the post kubernetes clusters default response
func (o *PostKubernetesClustersDefault) Code() int {
	return o._statusCode
}

func (o *PostKubernetesClustersDefault) Error() string {
	return fmt.Sprintf("[POST /kubernetesClusters][%d] PostKubernetesClusters default  %+v", o._statusCode, o.Payload)
}

func (o *PostKubernetesClustersDefault) GetPayload() *APIResponse {
	return o.Payload
}

func (o *PostKubernetesClustersDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(APIResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type APIResponse struct {

	// message
	Message string `json:"message,omitempty"`
}

type PostKubernetesClustersCreated struct {
	Payload *KubernetesCluster
}

func (o *PostKubernetesClustersCreated) Error() string {
	return fmt.Sprintf("[POST /kubernetesClusters][%d] postKubernetesClustersCreated  %+v", 201, o.Payload)
}

func (o *PostKubernetesClustersCreated) GetPayload() *KubernetesCluster {
	return o.Payload
}

func (o *PostKubernetesClustersCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(KubernetesCluster)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type PostKubernetesClustersParams struct {
	Cluster    *KubernetesCluster
	HTTPClient *http.Client
	Context    context.Context
	Timeout    time.Duration
}

// WithTimeout adds the timeout to the post kubernetes clusters params
func (o *PostKubernetesClustersParams) WithTimeout(timeout time.Duration) *PostKubernetesClustersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post kubernetes clusters params
func (o *PostKubernetesClustersParams) SetTimeout(timeout time.Duration) {
	o.Timeout = timeout
}

// WithContext adds the context to the post kubernetes clusters params
func (o *PostKubernetesClustersParams) WithContext(ctx context.Context) *PostKubernetesClustersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post kubernetes clusters params
func (o *PostKubernetesClustersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post kubernetes clusters params
func (o *PostKubernetesClustersParams) WithHTTPClient(client *http.Client) *PostKubernetesClustersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post kubernetes clusters params
func (o *PostKubernetesClustersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the post kubernetes clusters params
func (o *PostKubernetesClustersParams) WithBody(body *KubernetesCluster) *PostKubernetesClustersParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the post kubernetes clusters params
func (o *PostKubernetesClustersParams) SetBody(body *KubernetesCluster) {
	o.Cluster = body
}

// WriteToRequest writes these params to a swagger request
func (o *PostKubernetesClustersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.Timeout); err != nil {
		return err
	}
	var res []error

	if o.Cluster != nil {
		if err := r.SetBodyParam(o.Cluster); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// PostKubernetesClustersReader is a Reader for the PostKubernetesClusters structure.
type PostKubernetesClustersReader struct {
	Formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostKubernetesClustersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostKubernetesClustersCreated()
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewPostKubernetesClustersDefault(response.Code())
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

func NewPostKubernetesClustersDefault(code int) *PostKubernetesClustersDefault {
	return &PostKubernetesClustersDefault{
		_statusCode: code,
	}
}
func NewPostKubernetesClustersCreated() *PostKubernetesClustersCreated {
	return &PostKubernetesClustersCreated{}
}
