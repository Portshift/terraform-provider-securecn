package model

import (
	"context"
	"fmt"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"

	//"github.com/go-openapi/validate"
	"io"
	"net/http"
)

// PutKubernetesClustersKubernetesClusterIDHandlerFunc turns a function with the right signature into a put kubernetes clusters kubernetes cluster ID handler
type PutKubernetesClustersKubernetesClusterIDHandlerFunc func(PutKubernetesClustersKubernetesClusterIDParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn PutKubernetesClustersKubernetesClusterIDHandlerFunc) Handle(params PutKubernetesClustersKubernetesClusterIDParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// PutKubernetesClustersKubernetesClusterIDHandler interface for that can handle valid put kubernetes clusters kubernetes cluster ID params
type PutKubernetesClustersKubernetesClusterIDHandler interface {
	Handle(PutKubernetesClustersKubernetesClusterIDParams, interface{}) middleware.Responder
}

// NewPutKubernetesClustersKubernetesClusterID creates a new http.Handler for the put kubernetes clusters kubernetes cluster ID operation
func NewPutKubernetesClustersKubernetesClusterID(ctx *middleware.Context, handler PutKubernetesClustersKubernetesClusterIDHandler) *PutKubernetesClustersKubernetesClusterID {
	return &PutKubernetesClustersKubernetesClusterID{Context: ctx, Handler: handler}
}

/*PutKubernetesClustersKubernetesClusterID swagger:route PUT /kubernetesClusters/{kubernetesClusterId} kubernetes putKubernetesClustersKubernetesClusterId

Update the Kubernetes cluster

*/
type PutKubernetesClustersKubernetesClusterID struct {
	Context *middleware.Context
	Handler PutKubernetesClustersKubernetesClusterIDHandler
}

func (o *PutKubernetesClustersKubernetesClusterID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutKubernetesClustersKubernetesClusterIDParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

func NewPutKubernetesClustersKubernetesClusterIDParams() PutKubernetesClustersKubernetesClusterIDParams {

	return PutKubernetesClustersKubernetesClusterIDParams{}
}

// PutKubernetesClustersKubernetesClusterIDParams contains all the bound params for the put kubernetes clusters kubernetes cluster ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutKubernetesClustersKubernetesClusterID

/******************************/

type PutKubernetesClustersKubernetesClusterIDParamsWriter struct {

	/*Body*/
	Body *KubernetesCluster
	/*KubernetesClusterID
	  SecureCN Kubernetes cluster ID
	*/
	KubernetesClusterID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) WithTimeout(timeout time.Duration) *PutKubernetesClustersKubernetesClusterIDParamsWriter {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) WithContext(ctx context.Context) *PutKubernetesClustersKubernetesClusterIDParamsWriter {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) WithHTTPClient(client *http.Client) *PutKubernetesClustersKubernetesClusterIDParamsWriter {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) WithBody(body *KubernetesCluster) *PutKubernetesClustersKubernetesClusterIDParamsWriter {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) SetBody(body *KubernetesCluster) {
	o.Body = body
}

// WithKubernetesClusterID adds the kubernetesClusterID to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) WithKubernetesClusterID(kubernetesClusterID strfmt.UUID) *PutKubernetesClustersKubernetesClusterIDParamsWriter {
	o.SetKubernetesClusterID(kubernetesClusterID)
	return o
}

// SetKubernetesClusterID adds the kubernetesClusterId to the put kubernetes clusters kubernetes cluster ID params
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) SetKubernetesClusterID(kubernetesClusterID strfmt.UUID) {
	o.KubernetesClusterID = kubernetesClusterID
}

// WriteToRequest writes these params to a swagger request
func (o *PutKubernetesClustersKubernetesClusterIDParamsWriter) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param kubernetesClusterId
	if err := r.SetPathParam("kubernetesClusterId", o.KubernetesClusterID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

type PutKubernetesClustersKubernetesClusterIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Body *KubernetesCluster
	/*SecureCN Kubernetes cluster ID
	  Required: true
	  In: path
	*/
	KubernetesClusterID strfmt.UUID
} // BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutKubernetesClustersKubernetesClusterIDParams() beforehand.
func (o *PutKubernetesClustersKubernetesClusterIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body KubernetesCluster
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body"))
	}
	rKubernetesClusterID, rhkKubernetesClusterID, _ := route.Params.GetOK("kubernetesClusterId")
	if err := o.bindKubernetesClusterID(rKubernetesClusterID, rhkKubernetesClusterID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindKubernetesClusterID binds and validates parameter KubernetesClusterID from path.
func (o *PutKubernetesClustersKubernetesClusterIDParams) bindKubernetesClusterID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("kubernetesClusterId", "path", "strfmt.UUID", raw)
	}
	o.KubernetesClusterID = *(value.(*strfmt.UUID))

	if err := o.validateKubernetesClusterID(formats); err != nil {
		return err
	}

	return nil
}

// validateKubernetesClusterID carries on validations for parameter KubernetesClusterID
func (o *PutKubernetesClustersKubernetesClusterIDParams) validateKubernetesClusterID(formats strfmt.Registry) error {

	if err := validate.FormatOf("kubernetesClusterId", "path", "uuid", o.KubernetesClusterID.String(), formats); err != nil {
		return err
	}
	return nil
}

type PutKubernetesClustersKubernetesClusterIDReader struct {
	Formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PutKubernetesClustersKubernetesClusterIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPutKubernetesClustersKubernetesClusterIDOK()
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewPutKubernetesClustersKubernetesClusterIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPutKubernetesClustersKubernetesClusterIDOK creates a PutKubernetesClustersKubernetesClusterIDOK with default headers values
func NewPutKubernetesClustersKubernetesClusterIDOK() *PutKubernetesClustersKubernetesClusterIDOK {
	return &PutKubernetesClustersKubernetesClusterIDOK{}
}

/*PutKubernetesClustersKubernetesClusterIDOK handles this case with default header values.

OK
*/
type PutKubernetesClustersKubernetesClusterIDOK struct {
	Payload *KubernetesCluster
}

func (o *PutKubernetesClustersKubernetesClusterIDOK) Error() string {
	return fmt.Sprintf("[PUT /kubernetesClusters/{kubernetesClusterId}][%d] putKubernetesClustersKubernetesClusterIdOK  %+v", 200, o.Payload)
}

func (o *PutKubernetesClustersKubernetesClusterIDOK) GetPayload() *KubernetesCluster {
	return o.Payload
}

func (o *PutKubernetesClustersKubernetesClusterIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(KubernetesCluster)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutKubernetesClustersKubernetesClusterIDDefault creates a PutKubernetesClustersKubernetesClusterIDDefault with default headers values
func NewPutKubernetesClustersKubernetesClusterIDDefault(code int) *PutKubernetesClustersKubernetesClusterIDDefault {
	return &PutKubernetesClustersKubernetesClusterIDDefault{
		_statusCode: code,
	}
}

/*PutKubernetesClustersKubernetesClusterIDDefault handles this case with default header values.

unknown error
*/
type PutKubernetesClustersKubernetesClusterIDDefault struct {
	_statusCode int

	Payload *APIResponse
}

// Code gets the status code for the put kubernetes clusters kubernetes cluster ID default response
func (o *PutKubernetesClustersKubernetesClusterIDDefault) Code() int {
	return o._statusCode
}

func (o *PutKubernetesClustersKubernetesClusterIDDefault) Error() string {
	return fmt.Sprintf("[PUT /kubernetesClusters/{kubernetesClusterId}][%d] PutKubernetesClustersKubernetesClusterID default  %+v", o._statusCode, o.Payload)
}

func (o *PutKubernetesClustersKubernetesClusterIDDefault) GetPayload() *APIResponse {
	return o.Payload
}

func (o *PutKubernetesClustersKubernetesClusterIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(APIResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
