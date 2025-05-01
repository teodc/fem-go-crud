package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

func WriteJSONResponse(w http.ResponseWriter, status int, data Envelope) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	output = append(output, '\n')
	_, err = w.Write(output)
	if err != nil {
		return err
	}

	return nil
}

func ParseIDParamFromURL(r *http.Request, paramName string) (int, error) {
	idParam := chi.URLParam(r, paramName)
	if idParam == "" {
		return 0, errors.New("missing param")
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid param")
	}

	return int(id), nil
}
