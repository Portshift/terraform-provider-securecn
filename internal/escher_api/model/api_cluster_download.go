package model

import (
	"fmt"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"io"
	"net/http"
	"net/url"
	golangswaggerpaths "path"
	"strings"
)

// GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzHandlerFunc turns a function with the right signature into a get kubernetes clusters kubernetes cluster ID SecureCN bundle tar gz handler
type GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzHandlerFunc func(GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzHandlerFunc) Handle(params GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzHandler interface for that can handle valid get kubernetes clusters kubernetes cluster ID SecureCN bundle tar gz params
type GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzHandler interface {
	Handle(GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams, interface{}) middleware.Responder
}

// NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz creates a new http.Handler for the get kubernetes clusters kubernetes cluster ID SecureCN bundle tar gz operation
func NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz(ctx *middleware.Context, handler GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzHandler) *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz {
	return &GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz{Context: ctx, Handler: handler}
}

/*GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz swagger:route GET /kubernetesClusters/{kubernetesClusterId}/download_bundle kubernetes getKubernetesClustersKubernetesClusterIdSecureCNBundleTarGz
Get SecureCN installation script
In order to install, you need to extract and run "./install_bundle.sh"
*/
type GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz struct {
	Context *middleware.Context
	Handler GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzHandler
}

func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams()

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

// NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams creates a new GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams object
// no default values defined in spec.
func NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams() GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams {

	return GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams{}
}

// GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams contains all the bound params for the get kubernetes clusters kubernetes cluster ID SecureCN bundle tar gz operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz
type GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*SecureCN Kubernetes cluster ID
	  Required: true
	  In: path
	*/
	KubernetesClusterID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams() beforehand.
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

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
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams) bindKubernetesClusterID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzParams) validateKubernetesClusterID(formats strfmt.Registry) error {

	if err := validate.FormatOf("kubernetesClusterId", "path", "uuid", o.KubernetesClusterID.String(), formats); err != nil {
		return err
	}
	return nil
}

type GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzReader struct {
	Formats strfmt.Registry
	Writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK(o.Writer)
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault(response.Code())
		if err := result.readResponse(response, consumer, o.Formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK creates a GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK with default headers values
func NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK(writer io.Writer) *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK {
	return &GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK{
		Payload: writer,
	}
}

/*GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK handles this case with default header values.

OK
*/
type GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK struct {
	Payload io.Writer
}

func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK) Error() string {
	return fmt.Sprintf("[GET /kubernetesClusters/{kubernetesClusterId}/download_bundle][%d] getKubernetesClustersKubernetesClusterIdSecureCNBundleTarGzOK  %+v", 200, o.Payload)
}

func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault creates a GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault with default headers values
func NewGetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault(code int) *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault {
	return &GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault{
		_statusCode: code,
	}
}

/*GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault handles this case with default header values.

unknown error
*/
type GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault struct {
	_statusCode int

	Payload *APIResponse
}

// Code gets the status code for the get kubernetes clusters kubernetes cluster ID SecureCN bundle tar gz default response
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault) Code() int {
	return o._statusCode
}

func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault) Error() string {
	return fmt.Sprintf("[GET /kubernetesClusters/{kubernetesClusterId}/download_bundle][%d] GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGz default  %+v", o._statusCode, o.Payload)
}

func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault) GetPayload() *APIResponse {
	return o.Payload
}

func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(APIResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL struct {
	KubernetesClusterID strfmt.UUID

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL) WithBasePath(bp string) *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/kubernetesClusters/{kubernetesClusterId}/download_bundle"

	kubernetesClusterID := o.KubernetesClusterID.String()
	if kubernetesClusterID != "" {
		_path = strings.Replace(_path, "{kubernetesClusterId}", kubernetesClusterID, -1)
	} else {
		return nil, errors.New(1, "kubernetesClusterId is required on GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New(1, "scheme is required for a full url on GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL")
	}
	if host == "" {
		return nil, errors.New(1, "host is required for a full url on GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetKubernetesClustersKubernetesClusterIDSecureCNBundleTarGzURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
