package plugin

import "encoding/json"

type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Detail  string   `json:"detail"`
	Data    []string `json:"data,omitempty"`
}

func NewResponse(data []byte) (*Response, error) {
	r := new(Response)
	err := json.Unmarshal(data, r)
	return r, err
}
