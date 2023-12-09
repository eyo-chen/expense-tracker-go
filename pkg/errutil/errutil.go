package errutil

import (
	"net/http"

	"github.com/OYE0303/expense-tracker-go/pkg/jsutil"
	"github.com/OYE0303/expense-tracker-go/pkg/logger"
)

// ServerErrorResponse is a helper function for returning a 500 Internal Server Error response
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusInternalServerError, err.Error())
}

// BadRequestResponse is a helper function for returning a 400 Bad Request response
func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// AuthenticationErrorResponse is a helper function for returning a 401 Unauthorized response
func AuthenticationErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusUnauthorized, err.Error())
}

// VildateErrorResponse is a helper function for returning a 400 Bad Request response,
// and sending the validation error message in JSON format
func VildateErrorResponse(w http.ResponseWriter, r *http.Request, err map[string]string) {
	var errMap = make(map[string]interface{})
	for k, v := range err {
		errMap[k] = v
	}

	if err := jsutil.WriteJSON(w, http.StatusBadRequest, errMap, nil); err != nil {
		logger.Error("jsutil.WriteJSON failed", "package", "errutil", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	err := map[string]interface{}{"error": message}

	if err := jsutil.WriteJSON(w, status, err, nil); err != nil {
		logger.Error("jsutil.WriteJSON failed", "package", "errutil", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
