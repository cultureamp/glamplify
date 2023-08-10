package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	jwtgo "github.com/golang-jwt/jwt/v4"
)

// Decoder represents how to decode a JWT
type Decoder struct {
	defaultVerifyKey *rsa.PublicKey            // Web Gateway does not provide a kid header
	additionalKeys   map[string]*rsa.PublicKey // Optional to accept jwts signed by other services
}

// NewDecoder creates a new Decoder
func NewDecoder() (Decoder, error) {
	pubKey := os.Getenv("AUTH_PUBLIC_KEY")
	return NewDecoderFromBytes([]byte(pubKey))
}

// NewDecoderFromPath creates a new Decoder with the public key in 'pubKeyPath'
func NewDecoderFromPath(pubKeyPath string) (Decoder, error) {
	verifyBytes, _ := ioutil.ReadFile(filepath.Clean(pubKeyPath))
	return NewDecoderFromBytes(verifyBytes)
}

// NewDecoderFromBytes creates a new Decoder given the public key as a []byte
func NewDecoderFromBytes(verifyBytes []byte) (Decoder, error) {
	return NewMultiKeyDecoderFromBytes(verifyBytes, make(map[string][]byte))
}

// NewDecoderFromBytes creates a new Decoder given the public key as a []byte
func NewMultiKeyDecoderFromBytes(verifyBytes []byte, additionalVerifyBytes map[string][]byte) (Decoder, error) {
	verifyKey, err := jwtgo.ParseRSAPublicKeyFromPEM(verifyBytes)
	additionalKeys := make(map[string]*rsa.PublicKey)
	for key, bytes := range additionalVerifyBytes {
		additionalKey, err := jwtgo.ParseRSAPublicKeyFromPEM(bytes)
		if err != nil {
			return Decoder{}, err
		}
		additionalKeys[key] = additionalKey
	}
	return Decoder{
		defaultVerifyKey: verifyKey,
		additionalKeys:   additionalKeys,
	}, err
}

// Decode a jwt token and return the Payload
func (jwt Decoder) Decode(tokenString string) (Payload, error) {
	// sample token string in the form "header.payload.signature"
	//eg. "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	data := Payload{}

	token, err := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		return fetchPublicKey(token, jwt), nil
	})

	if err != nil {
		return data, err
	}

	if claims, ok := token.Claims.(jwtgo.MapClaims); ok && token.Valid {
		data.Customer, err = jwt.extractData(claims, "accountId")
		if err != nil {
			return data, err
		}
		data.RealUser, err = jwt.extractData(claims, "realUserId")
		if err != nil {
			return data, err
		}
		data.EffectiveUser, err = jwt.extractData(claims, "effectiveUserId")
		if err != nil {
			return data, err
		}
		return data, nil
	}

	return data, errors.New("invalid claim token in jwt")
}

func fetchPublicKey(token *jwtgo.Token, jwt Decoder) *rsa.PublicKey {
	kid, found := token.Header["kid"]
	if !found {
		return jwt.defaultVerifyKey
	}

	key, found := jwt.additionalKeys[kid.(string)]
	if !found {
		return jwt.defaultVerifyKey
	}

	return key
}

func (jwt Decoder) extractData(claims jwtgo.MapClaims, key string) (string, error) {
	val, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("missing %s in jwt token", key)
	}

	return val, nil
}
