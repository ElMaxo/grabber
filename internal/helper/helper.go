package helper

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/loads"
)

func StringPtr(in string) *string {
	return &in
}

func Int64Ptr(in int64) *int64 {
	return &in
}

func SafeStringGet(in *string) string {
	if in != nil {
		return *in
	}
	return ""
}

func SafeInt64Get(in *int64) int64 {
	if in != nil {
		return *in
	}
	return 0
}

func NewTransport(h http.Handler) http.RoundTripper {
	return &handlerTransport{h: h}
}

type handlerTransport struct {
	h http.Handler
}

func (s *handlerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	s.h.ServeHTTP(resp, req)
	return resp.Result(), nil
}

func ValidateSpec(orig, flat json.RawMessage) *loads.Document {
	swaggerSpec, err := loads.Embedded(orig, flat)
	if err != nil {
		log.Fatalln(err)
	}
	return swaggerSpec
}
