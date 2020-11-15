package authz

import "encoding/json"

type PolicyResponse struct {
	Result []EvaluationResponse `json:"result"`
}

type EvaluationResponse struct {
	Policy string `json:"policy"`
	Allow bool `json:"allow"`
	Status string `json:"status"`
	StatusCode int `json:"statusCode"`
}

type PolicyResponseParser struct {
}

func newPolicyEvalRequestParser() *PolicyResponseParser {
	return &PolicyResponseParser{}
}

func (p PolicyResponseParser) ParsePolicyEvalRequest(body string) (*PolicyResponse, error) {

	record := PolicyResponse{}
	err := json.Unmarshal([]byte(body), &record)
	return &record, err
}
