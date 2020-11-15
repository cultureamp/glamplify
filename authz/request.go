package authz

import "encoding/json"

type PolicyEvalRequest struct {
	Policies []PolicyRequest `json:"policies"`
}

type PolicyRequest struct {
	Policy string `json:"policy"`
	Context ContextRequest `json:"context"`
}

type ContextRequest struct {
	Identity IdentityRequest `json:"identity"`
	Input    InputRequest    `json:"input"`
}

type IdentityRequest struct {
	RealUserID string `json:"real_user_id"` // The real_user_id is the UUID of the user who logged in
	UserID     string `json:"user_id"`      // // The user_id is the UUID of the current user (either the user logged in, or if masquerading, the ID of the user we are masquerading as).
	AccountID  string `json:"account_id"`
}

type InputRequest map[string]interface{}


func (pr PolicyEvalRequest) GenerateJSON() (string, error) {
	b, err := json.Marshal(pr)
	return string(b), err
}
