package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ellofae/deanery-gateway/pkg/logger"
)

func RequestDecode(r *http.Request, req interface{}) error {
	logger := logger.GetLogger()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		logger.Printf("Unable to decode the request data. Error: %v.\n", err)
		return err
	}

	return nil
}
