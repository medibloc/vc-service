package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func RegisterHandlers(r *mux.Router) {
	// health check
	r.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// credentials
	r.HandleFunc("/credentials/issue", issueCredentialHandler).Methods("POST")
}

func healthCheckHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	io.WriteString(writer, `{"alive": true}`)
}

func issueCredentialHandler(writer http.ResponseWriter, request *http.Request) {
	var req issueCredentialRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		e := fmt.Errorf("failed to decode request: %w", err)
		log.Error(e)
		http.Error(writer, e.Error(), http.StatusBadRequest)
		return
	}

	vc, err := issueCredential(&req)
	if err != nil {
		e := fmt.Errorf("failed to issue credential: %w", err)
		log.Error(e)
		http.Error(writer, e.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write(vc)
}
