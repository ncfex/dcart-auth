package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Responder interface {
	RespondWithError(w http.ResponseWriter, code int, msg string, err error)
	RespondWithJSON(w http.ResponseWriter, code int, payload interface{})
}

type httpResponder struct {
	logger *log.Logger
}

func NewHTTPResponder(logger *log.Logger) Responder {
	return &httpResponder{
		logger: logger,
	}
}

func (r *httpResponder) RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		r.logger.Println(err)
	}
	if code > 499 {
		r.logger.Printf("Responding with 5XX error: %s", msg)
	}
	r.RespondWithJSON(w, code, ErrorResponse{
		Error: msg,
	})
}

func (r *httpResponder) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		r.logger.Printf("error marshaling json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
