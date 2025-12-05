package http

import (
	"encoding/json"
	"net/http"
)

// Response represents a standardized API response format.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// JSON sends a JSON response with the standardized format.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: status >= 200 && status < 300,
		Data:    data,
	}

	if !response.Success {
		response.Error = http.StatusText(status)
	}

	json.NewEncoder(w).Encode(response)
}

// Success sends a success response (200 OK) with data.
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// Created sends a created response (201 Created) with data.
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// Error sends an error response with the standardized format.
func Error(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: false,
		Error:   message,
	}

	if err != nil {
		response.Message = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

// BadRequest sends a bad request response (400 Bad Request).
func BadRequest(w http.ResponseWriter, message string, err error) {
	Error(w, http.StatusBadRequest, message, err)
}

// NotFound sends a not found response (404 Not Found).
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message, nil)
}

// InternalServerError sends an internal server error response (500 Internal Server Error).
func InternalServerError(w http.ResponseWriter, message string, err error) {
	Error(w, http.StatusInternalServerError, message, err)
}

// Unauthorized sends an unauthorized response (401 Unauthorized).
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a forbidden response (403 Forbidden).
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message, nil)
}
