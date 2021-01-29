package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cultureamp/glamplify/cache"
	"github.com/cultureamp/glamplify/env"
)

// ParameterStoreConfig configures how the parameter store should work
type ParameterStoreConfig struct {
	Profile string
	// todo other aws config?
	CacheErrorsAsEmpty bool
	CacheDuration      time.Duration
}

// ParameterStore allows easy access to retrieve configuration
type ParameterStore struct {
	conf    *ParameterStoreConfig
	session *session.Session
	ssm     *ssm.SSM
	cache   *cache.Cache
}

// NewParameterStore creates a new ParameterStore
func NewParameterStore(configure ...func(*ParameterStoreConfig)) *ParameterStore {
	conf := &ParameterStoreConfig{
		Profile:            env.GetString(env.AwsProfileEnv, "default"),
		CacheErrorsAsEmpty: false,
		CacheDuration:      1 * time.Minute,
	}
	for _, config := range configure {
		config(conf)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           conf.Profile, // eg. "default", or "dev-admin" etc
	}))
	ssm := ssm.New(sess)
	c := cache.New()
	return &ParameterStore{
		conf:    conf,
		session: sess,
		ssm:     ssm,
		cache:   c,
	}
}

// Get a secret from the parameter store for 'key'
func (ps ParameterStore) Get(key string) (string, error) {
	if x, found := ps.cache.Get(key); found {
		if val, ok := x.(string); ok {
			return val, nil
		}
	}

	// This makes a network call - can be slow...
	result, err := ps.ssm.GetParameter(&ssm.GetParameterInput{
		Name: &key,
	})
	val := ""
	if err == nil {
		val = *result.Parameter.Value
	} else if !ps.conf.CacheErrorsAsEmpty {
		return "", err
	}

	// cache this for a minute, in case multiple calls request the same key in a short duration
	ps.cache.Set(key, val, ps.conf.CacheDuration)
	return val, err
}
