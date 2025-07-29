package server

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := ErrorResponse{
		Error:   message,
		Message: getErrorMessage(statusCode),
		Code:    getErrorCode(statusCode),
	}

	json.NewEncoder(w).Encode(errorResp)
}

func writeJSONSuccess(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	successResp := SuccessResponse{
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(successResp)
}

func getErrorMessage(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusConflict:
		return "Conflict"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Error"
	}
}

func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "ERROR"
	}
}
