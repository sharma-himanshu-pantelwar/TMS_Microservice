package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	ResponseWriter http.ResponseWriter `json:"-"`
	StatusCode     int                 `json:"statusCode"`
	Headers        map[string]string   `json:"-"`
	Message        string              `json:"message,omitempty"`
	Error          string              `json:"error,omitempty"`
	Data           interface{}         `json:"data,omitempty"`
}

func (r *Response) Set() {
	r.ResponseWriter.WriteHeader(r.StatusCode)
	for key, value := range r.Headers {
		r.ResponseWriter.Header().Set(key, value)

	}
	err := json.NewEncoder(r.ResponseWriter).Encode(r)
	if err != nil {
		r.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		r.ResponseWriter.Write([]byte("Error Encoding data"))
	}
}
