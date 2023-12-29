package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/OYE0303/expense-tracker-go/pkg/logger"
	"github.com/gorilla/mux"
)

// WriteJSON writes the provided data to the client in JSON format.
func WriteJSON(w http.ResponseWriter, status int, data map[string]interface{}, headers http.Header) error {
	// json.Marshal converts the map to JSON
	js, err := json.Marshal(data)
	if err != nil {
		logger.Error("json.Marshal failed", "package", "jsutil", "err", err)
		return err
	}

	// EX: w.Header().Set("Content-Type", "application/json")
	//  => w.Header()["Content-Type"] = []string{"application/json"}
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// ReadJson decodes the JSON request body into the input dst.
func ReadJson(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	desc := json.NewDecoder(r.Body)
	desc.DisallowUnknownFields()

	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = desc.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// ReadID reads the ID from the URL path.
func ReadID(r *http.Request) (int64, error) {
	rawID := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		logger.Error("strconv.ParseInt failed", "package", "jsutil", "err", err)
		return 0, err
	}

	return id, nil
}
