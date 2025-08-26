package model

type APIError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type APIResponse struct {
	Error    *APIError   `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}
