package authz

import "encoding/json"

// PolicyEvalRequest contains the entire policy request
type PolicyEvalRequest struct {
	Policies []PolicyRequest `json:"policies"`
}

// PolicyRequest contains a single policy request
type PolicyRequest struct {
	Policy string `json:"policy"`
	Context ContextRequest `json:"context"`
}

// ContextRequest contains the context in the request
type ContextRequest struct {
	Identity IdentityRequest `json:"identity"`
	Input    InputRequest    `json:"input"`
}

// IdentityRequest contains the identity request
type IdentityRequest struct {
	RealUserID string `json:"real_user_id"` // The real_user_id is the UUID of the user who logged in
	UserID     string `json:"user_id"`      // // The user_id is the UUID of the current user (either the user logged in, or if masquerading, the ID of the user we are masquerading as).
	AccountID  string `json:"account_id"`
}

// InputRequest type for OPA
type InputRequest map[string]interface{}

// GenerateJSON returns JSON for the given PolicyEvalRequest
func (pr PolicyEvalRequest) GenerateJSON() (string, error) {
	b, err := json.Marshal(pr)
	return string(b), err
}
