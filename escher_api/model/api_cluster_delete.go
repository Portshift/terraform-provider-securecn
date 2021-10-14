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

type DeleteKubernetesClustersKubernetesClusterIDReader struct {
	Formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteKubernetesClustersKubernetesClusterIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteKubernetesClustersKubernetesClusterIDNoContent()
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDeleteKubernetesClustersKubernetesClusterIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteKubernetesClustersKubernetesClusterIDNoContent creates a DeleteKubernetesClustersKubernetesClusterIDNoContent with default headers values
func NewDeleteKubernetesClustersKubernetesClusterIDNoContent() *DeleteKubernetesClustersKubernetesClusterIDNoContent {
	return &DeleteKubernetesClustersKubernetesClusterIDNoContent{}
}

/*DeleteKubernetesClustersKubernetesClusterIDNoContent handles this case with default header values.

Success
*/
type DeleteKubernetesClustersKubernetesClusterIDNoContent struct {
}

func (o *DeleteKubernetesClustersKubernetesClusterIDNoContent) Error() string {
	return fmt.Sprintf("[DELETE /kubernetesClusters/{kubernetesClusterId}][%d] deleteKubernetesClustersKubernetesClusterIdNoContent ", 204)
}

func (o *DeleteKubernetesClustersKubernetesClusterIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteKubernetesClustersKubernetesClusterIDDefault creates a DeleteKubernetesClustersKubernetesClusterIDDefault with default headers values
func NewDeleteKubernetesClustersKubernetesClusterIDDefault(code int) *DeleteKubernetesClustersKubernetesClusterIDDefault {
	return &DeleteKubernetesClustersKubernetesClusterIDDefault{
		_statusCode: code,
	}
}

/*DeleteKubernetesClustersKubernetesClusterIDDefault handles this case with default header values.

unknown error
*/
type DeleteKubernetesClustersKubernetesClusterIDDefault struct {
	_statusCode int

	Payload *APIResponse
}

// Code gets the status code for the delete kubernetes clusters kubernetes cluster ID default response
func (o *DeleteKubernetesClustersKubernetesClusterIDDefault) Code() int {
	return o._statusCode
}

func (o *DeleteKubernetesClustersKubernetesClusterIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /kubernetesClusters/{kubernetesClusterId}][%d] DeleteKubernetesClustersKubernetesClusterID default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteKubernetesClustersKubernetesClusterIDDefault) GetPayload() *APIResponse {
	return o.Payload
}

func (o *DeleteKubernetesClustersKubernetesClusterIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(APIResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type DeleteKubernetesClustersKubernetesClusterIDParams struct {

	/*KubernetesClusterID
	  SecureCN Kubernetes cluster ID

	*/
	KubernetesClusterID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete kubernetes clusters kubernetes cluster ID params
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) WithTimeout(timeout time.Duration) *DeleteKubernetesClustersKubernetesClusterIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete kubernetes clusters kubernetes cluster ID params
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete kubernetes clusters kubernetes cluster ID params
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) WithContext(ctx context.Context) *DeleteKubernetesClustersKubernetesClusterIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete kubernetes clusters kubernetes cluster ID params
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete kubernetes clusters kubernetes cluster ID params
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) WithHTTPClient(client *http.Client) *DeleteKubernetesClustersKubernetesClusterIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete kubernetes clusters kubernetes cluster ID params
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithKubernetesClusterID adds the kubernetesClusterID to the delete kubernetes clusters kubernetes cluster ID params
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) WithKubernetesClusterID(kubernetesClusterID strfmt.UUID) *DeleteKubernetesClustersKubernetesClusterIDParams {
	o.SetKubernetesClusterID(kubernetesClusterID)
	return o
}

// SetKubernetesClusterID adds the kubernetesClusterId to the delete kubernetes clusters kubernetes cluster ID params
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) SetKubernetesClusterID(kubernetesClusterID strfmt.UUID) {
	o.KubernetesClusterID = kubernetesClusterID
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteKubernetesClustersKubernetesClusterIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param kubernetesClusterId
	if err := r.SetPathParam("kubernetesClusterId", o.KubernetesClusterID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
