package dto

type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}
