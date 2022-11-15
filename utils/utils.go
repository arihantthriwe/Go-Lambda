package utils

import "fmt"

type ErrorHandler struct {
	DevMessage string `json:"developer_message,omitempty"`
	Message    string `json:"message"`
	Method     string `json:"-"`
	Code       int    `json:"code,omitempty"`
}

func (h *ErrorHandler) Error() string {
	return fmt.Sprintf("Something went wrong with the %v request. Client Message: %v, Developer Message: %v", h.Method, h.Message, h.DevMessage)
}
