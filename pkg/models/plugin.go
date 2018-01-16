package models

type InspectResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Detail  string   `json:"detail"`
	Data    []string `json:"data,omitempty"`
}
