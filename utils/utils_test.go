package utils

import "testing"

func TestErrorHandler_Error(t *testing.T) {
	type fields struct {
		DevMessage string
		Message    string
		Method     string
		Code       int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &ErrorHandler{
				DevMessage: tt.fields.DevMessage,
				Message:    tt.fields.Message,
				Method:     tt.fields.Method,
				Code:       tt.fields.Code,
			}
			if got := h.Error(); got != tt.want {
				t.Errorf("ErrorHandler.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
