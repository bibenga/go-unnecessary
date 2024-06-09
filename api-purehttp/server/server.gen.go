//go:build go1.22

// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /v1/status)
	GetStatusV1(w http.ResponseWriter, r *http.Request, params GetStatusV1Params)

	// (POST /v1/status)
	SetStatusV1(w http.ResponseWriter, r *http.Request)
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

		err = runtime.BindStyledParameterWithOptions("simple", "X-Page", valueList[0], &XPage, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
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

		err = runtime.BindStyledParameterWithOptions("simple", "X-Page-Size", valueList[0], &XPageSize, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-Page-Size", Err: err})
			return
		}

		params.XPageSize = &XPageSize

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetStatusV1(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SetStatusV1 operation middleware
func (siw *ServerInterfaceWrapper) SetStatusV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	ctx = context.WithValue(ctx, ApiKeyAuthScopes, []string{})

	ctx = context.WithValue(ctx, FirebaseAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SetStatusV1(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
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
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
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

	m.HandleFunc("GET "+options.BaseURL+"/v1/status", wrapper.GetStatusV1)
	m.HandleFunc("POST "+options.BaseURL+"/v1/status", wrapper.SetStatusV1)

	return m
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

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

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
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
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
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RW32/bNhD+VwRuj4otxVmA6s2Jm8HrkBpxuw3I/EBL55iFRDJHKo1X6H8fSEnRLzry",
	"hgVD3ySS33133x3v+I3EIpOCA9eKRN8IgpKCK7A/81zvgWsWU80Ef48o0CzHgmvg2nxSKdNqe/pFCW7W",
	"VLyHjJqvHxF2JCI/TBuOabmrpqW1oih8koCKkUljhEQ9Uq8655MrmtzBYw5Kv70PVzTxsCIrfLICzJhS",
	"TPAFcJa8PX9D6JWM5kiFsolpKM2vRCEBNSuzlrBYr1KqdwKzMQ8W7bPGjdb/0gYKzzSTKZCIhMQnZp1q",
	"EhHG9eyc+EQfJJS/8ABYmlAypYdbmkEXv6Y7isy7QvFVATZYpZHxBwNlXcaL2XmX8fLCyciPUA0pCp+Y",
	"tDKEhET3hq9CDyLfvGDF9gvEtgwWPWH7sh+Je84TFJZqLODQIXBGn1mWZyS6nM1+uvRJxnj5H5wmxfx2",
	"cfdxufgHWrhCf7n73ZhjkVi6bvXaw57d81vZDMKT6gdqKpfR9ppDzwyUog9HXaq3/bY+rnYzJlZtyCnV",
	"s0b6ohdN0487Et2f1Aj8vr4ohL6muToaUSqqLuDyuOfapvDJz6DXmupc/RY6sinzTu28e3fS9VPWYLfq",
	"bubLXz/fvR8VssIOdTR2Ic6R6cPaiFQNJMk+wMGkzPwxo8MeaGK7SVn65I+z+Wp59gEODTe1qHKGKBbX",
	"cCu+2d+a1eb4XmtpDt8whC1VUJ/fAkXAm1qQX37/RPyWEbvbt2LiYHwnhvmbc+8z5xCbSsKDN18tDZZp",
	"K99w5wlQlcBwcj6ZGf+EBE4lIxGZTYJJQHwiqd5bnaZP4bRJywPoIf8d6By58qhXHpz8aYrIlIMtKNP8",
	"O9VirCPNQAMqW89W/ccc8NCI/1grYkffIPVu0FLd5GnqQm6FSIHyFnSY7VV5oRtsAjuap/pIN627Zzis",
	"5BGSszX76xhT8FrjDoOgTezo28XG7z68zoPgP3tktJPoeGpc54jAdVUFxK+it25U6g5rN9Y5TT3Zl360",
	"uxd+Lea1yLmjKj8JXRlW/9ayTdNrLnuqn8hx68b+RZkUl9YvyZu2XqkWEo5DXK9si52NY/svU+uoFMoh",
	"7Rp0c9kHV33duerV2/dKJIe3q8NmDmjMofj/rsAtfPXKffLdpLkZj7YZtwbb/cY0svakLFe64+x+U2yK",
	"0gw+1S09x5REZGqGSrEp/g4AAP//cD2ecxsOAAA=",
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