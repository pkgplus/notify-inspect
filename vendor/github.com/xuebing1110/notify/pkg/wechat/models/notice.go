package models

import (
	"encoding/json"
)

type Notice struct {
	UserID   string   `json:"touser"`
	Template string   `json:"template_id"`
	Emphasis string   `json:"emphasis"`
	Page     string   `json:"page"`
	Values   []string `json:"values"`
}

func (n *Notice) ToJson() []byte {
	data, _ := json.Marshal(n)
	return data
}
