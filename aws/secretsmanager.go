package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/cultureamp/glamplify/cache"
	"github.com/cultureamp/glamplify/env"
)

// SecretsManagerConfig configures how the secrets manager should work
type SecretsManagerConfig struct {
	Profile string
	// todo other aws config?
	CacheErrorsAsEmpty bool
	CacheDuration      time.Duration
}

// SecretsManager allows easy access to retrieve secrets
type SecretsManager struct {
	conf           *SecretsManagerConfig
	session        *session.Session
	secretsManager *secretsmanager.SecretsManager
	cache          *cache.Cache
}

// NewSecretsManager creates a new SecretsManager
func NewSecretsManager(configure ...func(*SecretsManagerConfig)) *SecretsManager {
	conf := &SecretsManagerConfig{
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

	sm := secretsmanager.New(sess)
	c := cache.New()

	return &SecretsManager{
		conf:           conf,
		session:        sess,
		secretsManager: sm,
		cache:          c,
	}
}

// Get a secret by 'key'
func (sm SecretsManager) Get(key string) (string, error) {
	if x, found := sm.cache.Get(key); found {
		if val, ok := x.(string); ok {
			return val, nil
		}
	}

	// This makes a network call - can be slow...
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	result, err := sm.secretsManager.GetSecretValue(input)

	val := ""
	if err == nil {
		val = result.String()
	} else if !sm.conf.CacheErrorsAsEmpty {
		return "", err
	}

	// cache this for a minute, in case multiple calls request the same key in a short duration
	sm.cache.Set(key, val, sm.conf.CacheDuration)
	return val, err
}
