package authz

import "encoding/json"

// PolicyResponse contains the entire response
type PolicyResponse struct {
	Result []EvaluationResponse `json:"result"`
}

// EvaluationResponse contains a single policy response
type EvaluationResponse struct {
	Policy string `json:"policy"`
	Allow bool `json:"allow"`
	Status string `json:"status"`
	StatusCode int `json:"statusCode"`
}

// PolicyResponseParser allows to parse a response body to a PolicyResponse
type PolicyResponseParser struct {
}

func newPolicyEvalRequestParser() *PolicyResponseParser {
	return &PolicyResponseParser{}
}

// ParsePolicyEvalRequest parses a body response to a PolicyResponse
func (p PolicyResponseParser) ParsePolicyEvalRequest(body string) (*PolicyResponse, error) {

	record := PolicyResponse{}
	err := json.Unmarshal([]byte(body), &record)
	return &record, err
}
