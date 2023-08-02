// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.2 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /v1/status)
	GetStatusV1(w http.ResponseWriter, r *http.Request, params GetStatusV1Params)

	// (POST /v1/status)
	SetStatusV1(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /v1/status)
func (_ Unimplemented) GetStatusV1(w http.ResponseWriter, r *http.Request, params GetStatusV1Params) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /v1/status)
func (_ Unimplemented) SetStatusV1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetStatusV1 operation middleware
func (siw *ServerInterfaceWrapper) GetStatusV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetStatusV1Params

	// ------------- Optional query parameter "q" -------------

	err = runtime.BindQueryParameter("form", true, false, "q", r.URL.Query(), &params.Q)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "q", Err: err})
		return
	}

	// ------------- Optional query parameter "IsFull" -------------

	err = runtime.BindQueryParameter("form", true, false, "IsFull", r.URL.Query(), &params.IsFull)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "IsFull", Err: err})
		return
	}

	headers := r.Header

	// ------------- Optional header parameter "X-Page" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Page")]; found {
		var XPage int32
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "X-Page", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Page", runtime.ParamLocationHeader, valueList[0], &XPage)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-Page", Err: err})
			return
		}

		params.XPage = &XPage

	}

	// ------------- Optional header parameter "X-Page-Size" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Page-Size")]; found {
		var XPageSize int32
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "X-Page-Size", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Page-Size", runtime.ParamLocationHeader, valueList[0], &XPageSize)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-Page-Size", Err: err})
			return
		}

		params.XPageSize = &XPageSize

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetStatusV1(w, r, params)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SetStatusV1 operation middleware
func (siw *ServerInterfaceWrapper) SetStatusV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	ctx = context.WithValue(ctx, ApiKeyAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SetStatusV1(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/status", wrapper.GetStatusV1)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/status", wrapper.SetStatusV1)
	})

	return r
}

type AuthenticationErrorJSONResponse Error

type BadRequestJSONResponse Error

type PermissionDenidJSONResponse Error

type GetStatusV1RequestObject struct {
	Params GetStatusV1Params
}

type GetStatusV1ResponseObject interface {
	VisitGetStatusV1Response(w http.ResponseWriter) error
}

type GetStatusV1200ResponseHeaders struct {
	XPage      int32
	XPageCount int32
	XPageSize  int32
}

type GetStatusV1200JSONResponse struct {
	Body    GetStatusV1
	Headers GetStatusV1200ResponseHeaders
}

func (response GetStatusV1200JSONResponse) VisitGetStatusV1Response(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Page", fmt.Sprint(response.Headers.XPage))
	w.Header().Set("X-Page-Count", fmt.Sprint(response.Headers.XPageCount))
	w.Header().Set("X-Page-Size", fmt.Sprint(response.Headers.XPageSize))
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetStatusV1400JSONResponse struct{ BadRequestJSONResponse }

func (response GetStatusV1400JSONResponse) VisitGetStatusV1Response(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetStatusV1401JSONResponse struct {
	AuthenticationErrorJSONResponse
}

func (response GetStatusV1401JSONResponse) VisitGetStatusV1Response(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetStatusV1403JSONResponse struct{ PermissionDenidJSONResponse }

func (response GetStatusV1403JSONResponse) VisitGetStatusV1Response(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	return json.NewEncoder(w).Encode(response)
}

type SetStatusV1RequestObject struct {
	Body *SetStatusV1JSONRequestBody
}

type SetStatusV1ResponseObject interface {
	VisitSetStatusV1Response(w http.ResponseWriter) error
}

type SetStatusV1200JSONResponse GetStatusV1

func (response SetStatusV1200JSONResponse) VisitSetStatusV1Response(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type SetStatusV1400JSONResponse struct{ BadRequestJSONResponse }

func (response SetStatusV1400JSONResponse) VisitSetStatusV1Response(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type SetStatusV1401JSONResponse struct {
	AuthenticationErrorJSONResponse
}

func (response SetStatusV1401JSONResponse) VisitSetStatusV1Response(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type SetStatusV1403JSONResponse struct{ PermissionDenidJSONResponse }

func (response SetStatusV1403JSONResponse) VisitSetStatusV1Response(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (GET /v1/status)
	GetStatusV1(ctx context.Context, request GetStatusV1RequestObject) (GetStatusV1ResponseObject, error)

	// (POST /v1/status)
	SetStatusV1(ctx context.Context, request SetStatusV1RequestObject) (SetStatusV1ResponseObject, error)
}

type StrictHandlerFunc = runtime.StrictHttpHandlerFunc
type StrictMiddlewareFunc = runtime.StrictHttpMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// GetStatusV1 operation middleware
func (sh *strictHandler) GetStatusV1(w http.ResponseWriter, r *http.Request, params GetStatusV1Params) {
	var request GetStatusV1RequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetStatusV1(ctx, request.(GetStatusV1RequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetStatusV1")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetStatusV1ResponseObject); ok {
		if err := validResponse.VisitGetStatusV1Response(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// SetStatusV1 operation middleware
func (sh *strictHandler) SetStatusV1(w http.ResponseWriter, r *http.Request) {
	var request SetStatusV1RequestObject

	var body SetStatusV1JSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.SetStatusV1(ctx, request.(SetStatusV1RequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SetStatusV1")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(SetStatusV1ResponseObject); ok {
		if err := validResponse.VisitSetStatusV1Response(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RWUW/bNhD+K8Jtj4otxVmA6s2J28HokBrxOgzI/MBI55iFRDLkKY1X6L8PpORIsugq",
	"GxYMexNFfvfxvu945DdIZaGkQEEGkm+g0SgpDLrBvKQdCuIpIy7Fe62ltr9TKQgF2U+mVN5MT78YKew/",
	"k+6wYPbrR41bSOCHacsxrWfNtI5WVVUIGZpUc2WDQHJEGjTrQrhi2S0+lmjo7fdwxbJAN2RVCCvUBTeG",
	"S7FAwbO3528Jg5rRLmlQzpiW0g6Vlgo18dq1jKe0yhltpS7GdrDorrXb6IyXLlF8ZoXKERKIIQT7nxEk",
	"wAXNziEE2iush/iAug5hVM72N6zAPn7Ntkzz4ErLrwZ1izWkuXiwUN5nvJid9xkvL7yM4gTVkKIKwdrK",
	"NWaQ3Fm+Bj3IfPOClfdfMHVlsDgS9lj2E3nPRaaloxpLOPYIXLBnXpQFJJez2U+XIRRc1OPodVLMbxa3",
	"n5aLv6GFL/WXs9/POZWZo+tXr1scuLmw42YUv6p+8EDlC9r959GzQGPYw8ktHabDrj6+djMm1iGQV6pn",
	"0uxFL5bnn7aQ3L2qEYTH+mop6ZqV5mRGuWy6gG/HR1vbVCH8jLQmRqX5Lfa4qcpe7bx796rjZ1zAftV9",
	"mC9/+Xz7flTIBjvU0cbFtNSc9msrUnMhKf4R99YyO+JWhx2yzHWTuvTh97P5ann2EfctN3Oo+g4xPD3A",
	"nfh2/t7+bZfviFTdlLnYyqHycxF8FgJTWwN6H8xXS4vl5BIfzjyhNjUwnpxPZnYbUqFgikMCs0k0iSAE",
	"xWjnMpw+xdNW0AekIf8tUqmFCVhQL5z8Ye23RrpSsG2757ONrlmBhNq4SnS6PZao961sjxB2Lq2BaX7Q",
	"0nwo89yHvJcyRyY60KFPq/oottgMt6zM6UQfPPS9eFiDIyRna/7nKaboey03jqIusafjVpuw/2Q6j6J/",
	"7XnQNdHzSLgutUZBTRVA2GTvttGoO6zdlEqWB+pY+tG+XIUHMa9lKTxV+aukJrD5p5GdTd/bcmCOjRyP",
	"buNf1Kb4tH4xb9p5XzpIPA7xvY8ddjaOPX5Tuo0qaTzSrpHawz446uveUW9erVcy279dHbYdnHSJ1X93",
	"BG7wa1DPw//G5vZic824cyXdbWwj695xd5tqU9UQ/XRo36XOIYGpvUCqTfVXAAAA//+8DDjRwQ0AAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
