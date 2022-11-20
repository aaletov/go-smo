// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// APIBuffer defines model for APIBuffer.
type APIBuffer struct {
	BufNum    int      `json:"bufNum"`
	Current   *ReqWT   `json:"current,omitempty"`
	Processed []ReqWSE `json:"processed"`
}

// APIDevice defines model for APIDevice.
type APIDevice struct {
	Current *ReqWT   `json:"current,omitempty"`
	DevNum  int      `json:"devNum"`
	Done    []ReqWSE `json:"done"`
}

// APISource defines model for APISource.
type APISource struct {
	Generated []ReqWT `json:"generated"`
	SourceNum int     `json:"sourceNum"`
}

// ReqWSE defines model for ReqWSE.
type ReqWSE struct {
	End     time.Time `json:"end"`
	Request Request   `json:"request"`
	Start   time.Time `json:"start"`
}

// ReqWT defines model for ReqWT.
type ReqWT struct {
	Request Request   `json:"request"`
	Time    time.Time `json:"time"`
}

// Request defines model for Request.
type Request struct {
	RequestNumber int `json:"requestNumber"`
	SourceNumber  int `json:"sourceNumber"`
}

// WaveInfo defines model for WaveInfo.
type WaveInfo struct {
	Buffers   []APIBuffer `json:"buffers"`
	Devices   []APIDevice `json:"devices"`
	Done      []ReqWSE    `json:"done"`
	EndTime   time.Time   `json:"endTime"`
	Rejected  []ReqWSE    `json:"rejected"`
	Sources   []APISource `json:"sources"`
	StartTime time.Time   `json:"startTime"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all info for waveform
	// (GET /getWaveInfo)
	GetWaveNumber(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetWaveNumber operation middleware
func (siw *ServerInterfaceWrapper) GetWaveNumber(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetWaveNumber(w, r)
	})

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

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
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
		r.Get(options.BaseURL+"/getWaveInfo", wrapper.GetWaveNumber)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xWT2/aThD9Ktb+fkfHpumNU2kbVUhVmgakHKIcNmZsNvL+yezaCEV892rGNqZgq9Dm",
	"Zrwz897MvH3mTWRWO2vABC+mb8Jna9CSH2d3889VngPSD4fWAQYFfPRc5beVpqewdSCmQpkABaDYxSKr",
	"EMEEOvwfIRdT8V/aY6QtQHoPrw9LindoM/AeVpShAmh/TurihnJbdIkot2K3iwXCa6WQaj12JA8RnvYp",
	"9vkFskA1Znfzr1CrDE7bvLSVFdSjY1lZA+/dYYvXFh9pbmErHGquAAMow4VzX56SioVniJHWjzj3sfEB",
	"hSHu7RBOiINhyrlFLYOYipUMcBWUBrEv4gMqU4gWGvw5K+QwaiZIDOcCHDXXoXVVYiY71tzytLfL6TKv",
	"f2XL4SM0Oz6DRG8r/dwYxKni95seiRgTBoXHRwBD3B5kDXOT20F7ygH92cLunW5A3Cs2h4uKtX4yVOw9",
	"XIBVtbxg7zRNmto7eGy31ovm0VrQUDG6Jsu/l3Cf3w+lpxjvldCv8WAWo75JKKpVVmZNkBnfANBSlcTe",
	"qQBSf/IbWRSAibIiFkZSE2LRvItmd/NoCZJcrkJKWofg/DRND5JYXD5D5YKyRkzFzHBebjH6WUEFyhTR",
	"YusDUJlSZWA8D6rFmjmZrSG6TiYnKJvNJpF8nFgs0jbXp9/nX25uFzdX18kkWQdd8k4Atf+RLwCbr+AA",
	"1ZRjUnaKUFLMNztAsQb0TScfkkkyoeLWgZFOian4yK9i4WRYs2jSAsLhHS6Ap0wXWdJA5ivCaWIObME7",
	"S61Q6PVk0q2o/UxL50qVcXr64olK95fmTzLdM+Ht/76XewgVmsgwicjmUWtOnvXoK60lbhuykSzLiLTD",
	"W9zIGkjVTVEPWLMtPb4drauwV17bRMoSgq0TrFKa2e5p9ysAAP//yzOSBZ8JAAA=",
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