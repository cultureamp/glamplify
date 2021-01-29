package authz

import (
	"bytes"
	"github.com/cultureamp/glamplify/env"
	"github.com/go-errors/errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	cachecontrol "github.com/pquerna/cachecontrol/cacheobject"
	"golang.org/x/net/context"
)

// Config contains the authz client configuration values
type Config struct {
	Timeout       time.Duration
	CacheDuration time.Duration
}

// Client represents the authz client
type Client struct {
	config           Config
	authzAPIEndpoint string
	http             Transport
	cache            *cache.Cache
}

// NewClient creates a new authz Client
func NewClient(authzAPIEndpoint string, http Transport, configure ...func(*Config)) *Client {
	c := env.GetInt(env.AuthzCacheDurationEnv, 60)
	cacheDuration := time.Duration(c) * time.Second

	t := env.GetInt(env.AuthzClientTimeoutEnv, 10000) // 1- secs
	timeOutDuration := time.Duration(t) * time.Millisecond

	conf := Config{
		Timeout: timeOutDuration,
		CacheDuration:  cacheDuration,
	}

	for _, config := range configure {
		config(&conf)
	}

	return &Client{
		config:           conf,
		authzAPIEndpoint: authzAPIEndpoint,
		http:             http,
		cache:            cache.New(conf.CacheDuration, conf.CacheDuration*5),
	}
}

// EvaluateBooleanPolicy calls authz-api asking it to evaluate the policy, and then returns the result
func (client Client) EvaluateBooleanPolicy(ctx context.Context, policy string, identity IdentityRequest, input InputRequest) (*EvaluationResponse, error) {
	if item, found := client.cache.Get(policy); found {
		result, ok := item.(*EvaluationResponse)
		if ok {
			return result, nil
		}
	}

	response, result, err := client.evaluateBooleanPolicy(policy, identity, input)
	if err != nil {
		return nil, err // if there is a compile error, etc. assume the kill switch is OFF
	}

	controlDirective, err := client.parseResponseCacheControl(response)
	if err == nil && controlDirective.MaxAge > 0 {
		cacheExpiry := time.Duration(controlDirective.MaxAge) * time.Second
		client.cache.Set(policy, result, cacheExpiry)
	}

	return result, nil
}

func (client Client) parseResponseCacheControl(response *http.Response) (*cachecontrol.ResponseCacheDirectives, error) {
	if response == nil || response.Header == nil {
		return nil, nil
	}

	controlDirective, err := cachecontrol.ParseResponseCacheControl(response.Header.Get("Cache-Control"))
	if err != nil {
		return nil, err
	}

	return controlDirective, nil
}

func (client Client) evaluateBooleanPolicy(policy string, identity IdentityRequest, input InputRequest) (*http.Response, *EvaluationResponse, error) {
	postBody, err := client.createRequestPostBody(policy, identity, input)
	if err != nil {
		return nil, nil, err
	}

	httpCtx, _ := context.WithTimeout(context.Background(), client.config.Timeout)
	response, err := client.http.Post(httpCtx, client.authzAPIEndpoint, "application/json", bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		return response, nil, err
	}

	policyResponse, err := client.readResponse(response)
	if err != nil {
		return response, nil, err
	}

	return response, &policyResponse.Result[0], nil
}

func (client Client) createRequestPostBody(policy string, identity IdentityRequest, input InputRequest) (string, error) {
	data := PolicyEvalRequest{}
	policyRequest := PolicyRequest{
		Policy: policy,
		Context: ContextRequest{
			Identity: identity,
			Input:    input,
		},
	}
	data.Policies = append(data.Policies, policyRequest)

	postBody, err := data.GenerateJSON()
	if err != nil {
		return "", err
	}

	return postBody, nil
}

func (client Client) readResponse(response *http.Response) (*PolicyResponse, error) {
	if response == nil {
		return nil, errors.New("response is nil")
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	p := newPolicyEvalRequestParser()
	policyResponse, err := p.ParsePolicyEvalRequest(string(bodyBytes))
	if err != nil {
		return nil, err
	}

	return policyResponse, nil
}
