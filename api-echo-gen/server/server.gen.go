// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
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
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/v1/status)
	GetStatusV1(ctx echo.Context, params GetStatusV1Params) error

	// (POST /api/v1/status)
	SetStatusV1(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetStatusV1 converts echo context to params.
func (w *ServerInterfaceWrapper) GetStatusV1(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetStatusV1Params
	// ------------- Optional query parameter "q" -------------

	err = runtime.BindQueryParameter("form", true, false, "q", ctx.QueryParams(), &params.Q)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter q: %s", err))
	}

	// ------------- Optional query parameter "IsFull" -------------

	err = runtime.BindQueryParameter("form", true, false, "IsFull", ctx.QueryParams(), &params.IsFull)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter IsFull: %s", err))
	}

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Page" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Page")]; found {
		var XPage int32
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Page, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Page", runtime.ParamLocationHeader, valueList[0], &XPage)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Page: %s", err))
		}

		params.XPage = &XPage
	}
	// ------------- Optional header parameter "X-Page-Size" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Page-Size")]; found {
		var XPageSize int32
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Page-Size, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Page-Size", runtime.ParamLocationHeader, valueList[0], &XPageSize)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Page-Size: %s", err))
		}

		params.XPageSize = &XPageSize
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStatusV1(ctx, params)
	return err
}

// SetStatusV1 converts echo context to params.
func (w *ServerInterfaceWrapper) SetStatusV1(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{})

	ctx.Set(ApiKeyAuthScopes, []string{})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SetStatusV1(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/v1/status", wrapper.GetStatusV1)
	router.POST(baseURL+"/api/v1/status", wrapper.SetStatusV1)

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
	w.Header().Set("X-Page", fmt.Sprint(response.Headers.XPage))
	w.Header().Set("X-Page-Count", fmt.Sprint(response.Headers.XPageCount))
	w.Header().Set("X-Page-Size", fmt.Sprint(response.Headers.XPageSize))
	w.Header().Set("Content-Type", "application/json")
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

	// (GET /api/v1/status)
	GetStatusV1(ctx context.Context, request GetStatusV1RequestObject) (GetStatusV1ResponseObject, error)

	// (POST /api/v1/status)
	SetStatusV1(ctx context.Context, request SetStatusV1RequestObject) (SetStatusV1ResponseObject, error)
}

type StrictHandlerFunc = runtime.StrictEchoHandlerFunc
type StrictMiddlewareFunc = runtime.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// GetStatusV1 operation middleware
func (sh *strictHandler) GetStatusV1(ctx echo.Context, params GetStatusV1Params) error {
	var request GetStatusV1RequestObject

	request.Params = params

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetStatusV1(ctx.Request().Context(), request.(GetStatusV1RequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetStatusV1")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetStatusV1ResponseObject); ok {
		return validResponse.VisitGetStatusV1Response(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// SetStatusV1 operation middleware
func (sh *strictHandler) SetStatusV1(ctx echo.Context) error {
	var request SetStatusV1RequestObject

	var body SetStatusV1JSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.SetStatusV1(ctx.Request().Context(), request.(SetStatusV1RequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SetStatusV1")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(SetStatusV1ResponseObject); ok {
		return validResponse.VisitSetStatusV1Response(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RWUW/bNhD+K8Jtj4otxVmA6s2J28HokBrxOgzI/MBI55iFRDLkKY1X6L8PpORIsugq",
	"GxYMexNFfvfdfXc83jdIZaGkQEEGkm+g0SgpDLrFvKQdCuIpIy7Fe62ltr9TKQgF2U+mVN5sT78YKew/",
	"k+6wYPbrR41bSOCHacsxrXfNtLZWVVUIGZpUc2WNQHJEGjTnQrhi2S0+lmjo7X24YlmgG7IqhBXqghvD",
	"pVig4Nnb87eEQc1ojzQol5iW0i6Vlgo18TprGU9plTPaSl2MebDonrVudNZLFyg+s0LlCAnEEIL9zwgS",
	"4IJm5xAC7RXWS3xAXZswKmf7G1ZgH79mW6Z5cKXlV4O6xRrSXDxYKO8zXszO+4yXF15GcYJqSFGFYNPK",
	"NWaQ3Fm+Bj2IfPOClfdfMHVlsDgS9lj2E3HPRaaloxoLOPYIXLBnXpQFJJez2U+XIRRc1OvodVLMbxa3",
	"n5aLv6GFL/SXu9+POZWZo+tXrzscuL2wk80oflX94IHKZ7T7z6Nngcawh5MuHbbDrj6+djMm1sGQV6pn",
	"0uxFL5bnn7aQ3L2qEYTH+mop6ZqV5mREuWy6gM/jI9c2VQg/I62JUWl+iz3ZVGWvdt69e9X1M85gv+o+",
	"zJe/fL59Pypkgx3qaO1iWmpO+7UVqXmQFP+Ie5syu+JWhx2yzHWTuvTh97P5ann2EfctN3Oo+g0xPD3A",
	"nfh2/97+bY/viFTdlLnYyqHycxF8FgJTWwN6H8xXS4vl5AIf7jyhNjUwnpxPZtYNqVAwxSGB2SSaRBCC",
	"YrRzEU6Z4tOneNqK+oA09OEWqdTCBCyoD07+sCVgk+nKwbbuXq4tg2YFEmrjqtFp91ii3rfSPULYebgG",
	"ifODluZDmec+5L2UOTLRgQ5ztaqvY4vNcMvKnE70wkPvi4d1OEJytuZ/nmKKvtd24yjqEnu6brUJ+2PT",
	"eRT9ayNCN4meQeG61BoFNVUAYRO9c6NRd1i/KZUsD9Sx9KO9uQoPYl7LUniq8ldJjWHzTy27NH3P5cAc",
	"J3LcurV/USfFp/VL8qadGdNB4nGIb0Z22Nk49niudI4qaTzSrpHayz646uveVW8m1yuZ7d+uDtsuTrrE",
	"6r+7Ajf4Naj34X+T5vZxc8248yzdbWwj675zd5tqU9UQ/XRo36XOIYGpfdT/CgAA///j9PvEwg0AAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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
