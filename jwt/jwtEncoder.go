package jwt

import (
	"crypto/rsa"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Encoder represents a jwt encoder
type Encoder struct {
	pemKey *rsa.PrivateKey
}

// Claims contains the claims to be used to sign JWT's returned by Identity API
type claims struct {
	AccountID       string `json:"accountId"`
	EffectiveUserID string `json:"effectiveUserId"`
	RealUserID      string `json:"realUserId"`
	jwtgo.StandardClaims
}

// NewEncoder creates a new Encoder
func NewEncoder() (Encoder, error) {
	priKey := os.Getenv("AUTH_PRIVATE_KEY")
	return NewEncoderFromBytes([]byte(priKey))
}

// NewEncoderFromPath creates a new Encoder given the private key at 'pemKeyPath'
func NewEncoderFromPath(pemKeyPath string) (Encoder, error) {
	pemBytes, _ := ioutil.ReadFile(filepath.Clean(pemKeyPath))
	return NewEncoderFromBytes(pemBytes)
}

// NewEncoderFromBytes creates a new Encoder given the private key as a []byte
func NewEncoderFromBytes(pemBytes []byte) (Encoder, error) {
	pemKey, err := jwtgo.ParseRSAPrivateKeyFromPEM(pemBytes)
	return Encoder{
		pemKey: pemKey,
	}, err
}

// Encode a Payload
func (encoder Encoder) Encode(payload Payload) (string, error) {
	// Were a little loose on the expiry for now, to avoid possible
	// problems with clock skew, slow requests, background jobs (?) etc.
	expiry := 10 * time.Minute
	claims := encoder.claims(payload, expiry)
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, claims)
	return token.SignedString(encoder.pemKey)
}

// EncodeWithExpiry encodes a Payload with an expiry
func (encoder Encoder) EncodeWithExpiry(payload Payload, duration time.Duration) (string, error) {
	claims := encoder.claims(payload, duration)
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, claims)
	return token.SignedString(encoder.pemKey)
}

func (encoder Encoder) claims(payload Payload, duration time.Duration) claims {
	now := time.Now()
	return claims{
		AccountID:       payload.Customer,
		EffectiveUserID: payload.EffectiveUser,
		RealUserID:      payload.RealUser,
		StandardClaims: jwtgo.StandardClaims{
			IssuedAt: now.Unix(),
			// Were a little loose on the expiry for now, to avoid possible
			// problems with clock skew, slow requests, background jobs (?) etc.
			ExpiresAt: now.Add(duration).Unix(),
		},
	}
}