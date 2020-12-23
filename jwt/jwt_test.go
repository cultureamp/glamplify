package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

const (
	day  = 24 * time.Hour
	year = 365 * day // approx
)

func Test_JWT_Encode(t *testing.T) {

	jwt, err := NewEncoderFromPath("jwt.rs256.key.development.pem")
	assert.Nil(t, err)

	expiry := 25 * year
	token, err := jwt.EncodeWithExpiry(Payload{
		Customer:      "abc123",
		RealUser:      "xyz234",
		EffectiveUser: "xyz345",
	}, expiry)

	assert.Nil(t, err)
	splitToken := strings.Split(token, ".")
	assert.Equal(t, 3, len(splitToken))

	header := splitToken[0]
	assert.Equal(t, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9", header)
	assert.NotEmpty(t, token)
	// fmt.Print(token)
	// Note: Hard to test payload/signature because there is a elapsed time field that always changes...
}

func Test_JWT_Decode(t *testing.T) {
	jwt, err := NewDecoderFromPath("jwt.rs256.key.development.pub")
	assert.Nil(t, err)

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxOTAzOTMwNzA0LCJpYXQiOjE1ODg1NzA3MDR9.XGm34FDIgtBFvx5yC2HTUu-cf3DaQI4TmIBVLx0H7y89oNVNWJaKA3dLvWS0oOZoYIuGhj6GzPREBEmou2f9JsUerqnc-_Tf8oekFZWU7kEfzu9ECBiSWPk7ljPJeZLbau62sSqD7rYb-m3v1mohqz4tKJ_7leWu9L1uHHliC7YGlSRl1ptVDllJjKXKjOg9ifeGSXDEMeU35KgCFwIwKdu8WmCTd8ztLSKEnLT1OSaRZ7MSpmHQ4wUZtS6qvhLBiquvHub9KdQmc4mYWLmfKdDiR5DH-aswJFGLVu3yisFRY8uSfeTPQRhQXd_UfdgifCTXdWTnCvNZT-BxULYG-5mlvAFu-JInTga_9-r-wHRzFD1SrcKjuECF7vUG8czxGNE4sPjFrGVyBxE6fzzcFsdrhdqS-LB_shVoG940fD-ecAhXQZ9VKgr-rmCvmxuv5vYI2HoMfg9j_-zeXkucKxvPYvDQZYMdeW4wFsUORliGplThoHEeRQxTX8d_gvZFCy_gGg0H57FmJwCRymWk9v29s6uyHUMor_r-e7e6ZlShFBrCPAghXL04S9IFJUxUv30wNie8aaSyvPuiTqCgGiEwF_20ZaHCgYX0zupdGm4pHTyJrx2wv31yZ4VZYt8tKjEW6-BlB0nxzLGk5OUN83vq-RzH-92WmY5kMndF6Jo"
	payload, err := jwt.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, "abc123", payload.Customer)
	assert.Equal(t, "xyz234", payload.RealUser)
	assert.Equal(t, "xyz345", payload.EffectiveUser)
}

func Test_JWT_Encode_Decode(t *testing.T) {

	jwtEncoder, err := NewEncoderFromPath("jwt.rs256.key.development.pem")
	assert.Nil(t, err)

	token, err := jwtEncoder.Encode(Payload{
		Customer:      "abc123",
		RealUser:      "xyz234",
		EffectiveUser: "xyz345",
	})

	jwtDecoder, err := NewDecoderFromPath("jwt.rs256.key.development.pub")
	assert.Nil(t, err)

	payload, err := jwtDecoder.Decode(token)

	assert.Nil(t, err)
	assert.Equal(t, "abc123", payload.Customer)
	assert.Equal(t, "xyz234", payload.RealUser)
	assert.Equal(t, "xyz345", payload.EffectiveUser)
}

func Test_PayloadFromRequest_NoAuthorizationHeader(t *testing.T) {
	jwt, err := NewDecoderFromPath("jwt.rs256.key.development.pub")
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	_, err = PayloadFromRequest(req, jwt)
	assert.NotNil(t, err)
	assert.Equal(t, "missing authorization header", err.Error())
}

func Test_PayloadFromRequest_NoBearer(t *testing.T) {
	jwt, err := NewDecoderFromPath("jwt.rs256.key.development.pub")
	assert.Nil(t, err)

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxNTg4NTY4OTY4LCJpYXQiOjE1ODg1NjgzNjh9.bQKvgB8ZfGSZDQnyfM1remW_cB_sF95iIS-QfJmdn3jOCK60xiwMB7cNkXBLeVsCyMScHiTyENePvCOudsMruNhWO8YvnBpv6712O4n2sRckTKaNdYXAwidILDFXzRZMzrFAOJu1zKjPSaXiXdEv6zgq3OMruXcBF9RzsOlKPoOlgBI1Q9ctGurgI-p4WovCA0YmV9w7I2c1t3WQMMlapoJPKW1-bM37sgEgJpJrmjCavYswQ_mWY0yk9h8ftXGvQRPvLXM_K-kkhmUJ1cLT-H4iXIZkCk-Y-ONAej9lPOgBGiCmOq5DHHcggOKzzqcT0YNKrZHfCrigd7ZbT-zRSw9ukzYafduabCSj9MAq_oKzYYbYpqu6yNtzHXFBZ7izWjGVMUxpQX5gaFh6W0aezWwmBL6drO1NzYDSMX2lJ-FwVCVfKbvqFPxS5mqYQCAQueGrTlrIndWqVdDbJFw2LHTFxVLAFQGgnM292WJYp6KYVKm07mRpzHdozb8ER7lfB_hlloudEBh14WxnV4iKZabjsGZmpzXldSdVKceXTBSY4jobE-vD_U2YfBcNU7y_A2qJtTnhdIWMq_UZrmi5ycV_Nq5MUSjLat-J8iFSkCeOEsyxQ3ybMiaxFEbpBZoZAIrRUJxx3KWtdzywyS4dN9frv36O0UuWomskaR1il6U"

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Authorization", "Beare "+token)

	_, err = PayloadFromRequest(req, jwt)
	assert.NotNil(t, err)
	assert.Equal(t, "missing 'Bearer' token in authorization header", err.Error())
}

func Test_PayloadFromRequest_InvalidToken(t *testing.T) {
	jwt, err := NewDecoderFromPath("jwt.rs256.key.development.pub")
	assert.Nil(t, err)

	token := "INVALID.TOKEN."

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	_, err = PayloadFromRequest(req, jwt)
	assert.NotNil(t, err)
}
