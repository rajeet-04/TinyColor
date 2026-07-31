package compat

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	ID        string `json:"id"`
	Operation string `json:"operation"`
	Input     any    `json:"input"`
	Args      any    `json:"args,omitempty"`
}

type Response struct {
	ID     string `json:"id"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Success(id string, result any) (Response, error) {
	response := Response{ID: id, Result: result}
	return response, response.Validate()
}

func Failure(id, message string) (Response, error) {
	response := Response{ID: id, Error: message}
	return response, response.Validate()
}

func (r Response) Validate() error {
	if (r.Result == nil) == (r.Error == "") {
		return fmt.Errorf("response must contain exactly one of result or error")
	}
	return nil
}

func Decode(line []byte) (Request, error) {
	var request Request
	if err := json.Unmarshal(line, &request); err != nil {
		return Request{}, fmt.Errorf("malformed JSON")
	}
	if request.ID == "" || request.Operation == "" {
		return Request{}, fmt.Errorf("id and operation are required")
	}
	return request, nil
}

func Encode(response Response) ([]byte, error) {
	if err := response.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(response)
}
