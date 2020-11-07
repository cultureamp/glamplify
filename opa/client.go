package opa

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
)

type HttpClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type OPAHttpClient struct {}

func NewOPAHttpClient() HttpClient {
	return &OPAHttpClient{}
}

func (client OPAHttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(url, contentType, body)
}

type OPAClient struct {
	authzAPIEndpoint string
	http             HttpClient
	cache            *cache.Cache
}

func NewOPAClient(authzAPIEndpoint string, http HttpClient) *OPAClient {
	return &OPAClient{
		authzAPIEndpoint: authzAPIEndpoint,
		http:             http,
		cache:            cache.New(1*time.Minute, 1*time.Minute), // TODO pass args and/or read from ENV
	}
}

func (client OPAClient) EvaluateBooleanPolicy(ctx context.Context, policy string, identity IdentityRequest, input InputRequest) (*EvaluationResponse, error) {

	if item, found := client.cache.Get(policy); found {
		result, ok := item.(*EvaluationResponse)
		if ok {
			return result, nil
		}
	}

	result, err := client.evaluateBooleanPolicy(ctx, policy, identity, input)
	if err != nil {
		return nil, err // if there is a compile error, etc. assume the kill switch is OFF
	}

	client.cache.Set(policy, result, cache.DefaultExpiration)  // todo should be Cache-Control header value
	return result, nil
}

func (client OPAClient) evaluateBooleanPolicy(ctx context.Context, policy string, identity IdentityRequest, input InputRequest) (*EvaluationResponse, error) {
	postBody, err := client.createRequestPostBody(ctx, policy, identity, input)

	response, err := client.http.Post(client.authzAPIEndpoint, "application/json", bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		return nil, err
	}

	policyResponse, err := client.readResponse(ctx, response)
	if err != nil {
		return nil, err
	}

	return &policyResponse.Result[0], nil
}

func (client OPAClient) createRequestPostBody(ctx context.Context, policy string, identity IdentityRequest, input InputRequest) (string, error) {
	data := PolicyEvalRequest{}
	policyRequest := PolicyRequest{
		Policy: policy,
		Context: ContextRequest{
			Identity: identity,
			Input: input,
		},
	}
	data.Policies = append(data.Policies, policyRequest)

	postBody, err := data.GenerateJSON()
	if err != nil {
		return "", err
	}

	return postBody, nil
}

func (client OPAClient) readResponse(ctx context.Context, response *http.Response) (*PolicyResponse, error) {

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