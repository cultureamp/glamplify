package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/cultureamp/glamplify/cache"
	"time"
)

// SecretsManager allows easy access to retrieve secrets
type SecretsManager struct {
	session        *session.Session
	secretsManager *secretsmanager.SecretsManager
	cache          *cache.Cache
}

// NewSecretsManager creates a new SecretsManager
func NewSecretsManager(profile string) *SecretsManager {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile, // eg. "default", or "dev-admin" etc
	}))

	sm := secretsmanager.New(sess)
	c := cache.New()

	return &SecretsManager{
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
		SecretId:     aws.String(key),
	}

	result, err := sm.secretsManager.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	val :=  result.String()

	// cache this for a minute, in case multiple calls request the same key in a short duration
	sm.cache.Set(key, val, 1 * time.Minute)
	return val, nil
}

