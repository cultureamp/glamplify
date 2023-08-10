package jwt

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup_decoder(t *testing.T) Decoder {
	defaultKeyBytes, _ := ioutil.ReadFile(filepath.Clean("jwt.rs256.key.development.pub"))
	extra_1, _ := ioutil.ReadFile(filepath.Clean("jwt.rs256.key.development.extra_1.pub"))
	extra_2, _ := ioutil.ReadFile(filepath.Clean("jwt.rs256.key.development.extra_2.pub"))
	additionalKeys := map[string][]byte{
		"extra_1": extra_1,
		"extra_2": extra_2,
	}
	jwt, err := NewMultiKeyDecoderFromBytes(defaultKeyBytes, additionalKeys)
	assert.Nil(t, err)
	return jwt
}
func Test_JWT_MultiKey_Decode_Default(t *testing.T) {
	jwt := setup_decoder(t)

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxOTAzOTMwNzA0LCJpYXQiOjE1ODg1NzA3MDR9.XGm34FDIgtBFvx5yC2HTUu-cf3DaQI4TmIBVLx0H7y89oNVNWJaKA3dLvWS0oOZoYIuGhj6GzPREBEmou2f9JsUerqnc-_Tf8oekFZWU7kEfzu9ECBiSWPk7ljPJeZLbau62sSqD7rYb-m3v1mohqz4tKJ_7leWu9L1uHHliC7YGlSRl1ptVDllJjKXKjOg9ifeGSXDEMeU35KgCFwIwKdu8WmCTd8ztLSKEnLT1OSaRZ7MSpmHQ4wUZtS6qvhLBiquvHub9KdQmc4mYWLmfKdDiR5DH-aswJFGLVu3yisFRY8uSfeTPQRhQXd_UfdgifCTXdWTnCvNZT-BxULYG-5mlvAFu-JInTga_9-r-wHRzFD1SrcKjuECF7vUG8czxGNE4sPjFrGVyBxE6fzzcFsdrhdqS-LB_shVoG940fD-ecAhXQZ9VKgr-rmCvmxuv5vYI2HoMfg9j_-zeXkucKxvPYvDQZYMdeW4wFsUORliGplThoHEeRQxTX8d_gvZFCy_gGg0H57FmJwCRymWk9v29s6uyHUMor_r-e7e6ZlShFBrCPAghXL04S9IFJUxUv30wNie8aaSyvPuiTqCgGiEwF_20ZaHCgYX0zupdGm4pHTyJrx2wv31yZ4VZYt8tKjEW6-BlB0nxzLGk5OUN83vq-RzH-92WmY5kMndF6Jo"
	payload, err := jwt.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, "abc123", payload.Customer)
	assert.Equal(t, "xyz234", payload.RealUser)
	assert.Equal(t, "xyz345", payload.EffectiveUser)
}

func Test_JWT_MultiKey_Decode_Extra_1(t *testing.T) {
	jwt := setup_decoder(t)

	// Signed with kid extra_1
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImV4dHJhXzEifQ.eyJyZWFsVXNlcklkIjoiYXNkZjEyMzQiLCJlZmZlY3RpdmVVc2VySWQiOiJhc2RmMTIzNCIsImFjY291bnRJZCI6ImRlZjQ1NiIsImlhdCI6MTY5MTY0Njc2MSwiZXhwIjoxNjkxNzMzMTYxfQ.l1cp6d85M2cMtdjy5PldD8ZUUYgBy9cvw2Bsnr5IEKMRqHf6RflkziwiEB9Bdx4GiYIwpu37BTNv7OC7KC_QPnM3oYWadzk_too17O7QMun_ni5JOCp5jbjpC10lJHmJH1npKXqswJz9XNVVvvfrakwAHmGAQhYCyio3N0eJYxwyADiSfZh-UME3vVewPMuVL8wJyCptQjeAog_br10nukMzbvZITXoZDCghqRwAUmTc7vIswWD0fSgk6cpX0JF3VbdOgPxSYiCsKH4ioOQ9tW3ATnwhyH94B0wazvcaPfO4yZHgkdy1ewftlGUozVRShP5aos9W0A92sqobcBGToHBVOq7KAaC76ratRlhLVkcHUquGReoCpJxky4-kDlnyJEfTL5bmAF4ao3H1vUMzsJgtI7AUjoAeB-3EcgfvZKxK40V3v59IMk0ARzxvk5S2rCcMBxKvjU9JNpIKbsuYIRh4pFvnxXNIFD-wM1ZHE0TXRIgGAE05jEwujjmhdNRY12qa38lLEvf32bftrlTNj-Hm1nCrfjVgQwZOdXVGpQ5ysV-dh4LZzgYIU6e9sFdCwnOtbXvjPP5FPITEaSTcuKS0lX6BoYel-ivR3nFVPcIOBFxKik2zdw0NL4ZnF4riP2yJxOdEDyfKAdjuZO7lwcHmVATJJGRydJXhCtdBLAg"
	payload, err := jwt.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, "def456", payload.Customer)
	assert.Equal(t, "asdf1234", payload.RealUser)
	assert.Equal(t, "asdf1234", payload.EffectiveUser)
}

func Test_JWT_MultiKey_Decode_Extra_2(t *testing.T) {
	jwt := setup_decoder(t)

	// Signed with kid extra_2
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImV4dHJhXzIifQ.eyJyZWFsVXNlcklkIjoiM2MiLCJlZmZlY3RpdmVVc2VySWQiOiIyYiIsImFjY291bnRJZCI6IjFhIiwiaWF0IjoxNjkxNjQ2ODg2LCJleHAiOjE2OTE3MzMyODZ9.ioz6ZCkVQQT4RTEDroclGTaq4ujRc353XU06ZnCaQNLLYdb1-X-62Sb3nnv1yzT8WKdAeNeBkVi1FW8_mhTh78txqRlwJZzCzl_6JUrqljW3TuKx6MC0uvultq4Uc-sM1iqqgWTJ0pr2ehH6qB7wx49GEu49HUm3P8GDonlnLK57T2x3b7t8n14Jc_XdXmak1WUvC9DHuVQVHHqnOtinLOw3GohfPzPQ3d5iBL2BV2P5LU5Hbv4S3AdnVO5vY-idahT5GkDEDKZKCO7vpeBMahiW92rMP-wCy4NJAABCwPGH-Bc_XtTUyAWhIiXR4Wr4USuwL_CmcfKWhL0IpHryDziGINyCRnt2B4NohPbMgT_VRlWQncaiggMRKFUHF-7QYvy8dE6LY6OuYfw9bNhTS7cEEEc0f8yuHIMAL6tQuxI3Szqfctb8TnDJgNK2wtfUBfHaiBqOnFal3TIlS9jYGBNWyM0f2zAsSeduzd8LV20gtGHxosCU8eZOhSLNWo-SpF08BYv4Uc3AYGsdIxr6dAMd0xI_Fp8th0ps-OySWiHFHM6kdn1JGB5fiPTDeA8CNXLxpeEMMjiIrt0jv3p6N-ZEAv4V6t87NW8vDMMMlmiSJbnTFDPkxfnVo_y7jQ1vln8zQmOFuTz35UeWYQVdIdcrp-cAQIxBq_JjJEWLyeo"
	payload, err := jwt.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, "1a", payload.Customer)
	assert.Equal(t, "2b", payload.RealUser)
	assert.Equal(t, "3c", payload.EffectiveUser)
}

func Test_JWT_MultiKey_Decode_Missing_Extra_3(t *testing.T) {
	jwt := setup_decoder(t)

	// Signed with kid extra_3
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImV4dHJhXzMifQ.eyJyZWFsVXNlcklkIjoiYXNkZjEyMzQiLCJlZmZlY3RpdmVVc2VySWQiOiJhc2RmMTIzNCIsImFjY291bnRJZCI6ImRlZjQ1NiIsImlhdCI6MTY5MTY0NjkzMiwiZXhwIjoxNjkxNzMzMzMyfQ.ZTMbyzC2BuHTsseZeuFYlHqJWUxK__Kn403wEfiqB1Ok0VSCgR8Krvtm3toHT18uw2cQ0TUI1xlB61x8Nz8kM_MAS2Y2IJ-iWjC8OzPcqHNNFnD7h50zZAsg9QxMKwi581dPGf7enk4BqWNQ59KxUvM8jB8XabE6cJ9yLXXp3rpSzLBZQmGpkX__WDtzmdgHQjEo3pdmRtVVl38iVBSPXUEMCGRf7eTUThXmbZnDKmvkS3RbExe2DxRNdZ9H7CoSJyONb3Xrg7eaBJ24mdkWzsJhRbbS8Qp6Df1d_5VORWq2Y8wEnkCU2LieL3EhJI84Z4siXWe0a4EOkt60o6X0BAS-nUNWEx1xnQWBse6vkGbnLpDJSdY7mmozrj2rgIID-QjAFwcWbvQPF5P_962jPCXInQt9J-3qMqVxIGfGIN0Uyyf1WPZtm3nwz6UtO-OD3uOCpeiVaxM2tVVvaSaJDeAJ_8OuH8IqOxMeqMh_3dbICvMjOLMdmQm5VWbYoPzbD75A_iWgwWz7T94Ip5jyh7b1wc1D-wNpemlfYGb8s4ybyq4VZstOiNVI3q5afYclpzFWNCrnZVUs5lmoNUar7WCfHxCEugJ3AvDQyl8y65eejigXmcg5nNcDsHQ_POufJLx97N1ljfjyKOi_mw9h_NGfERX3HDLDv4D1xM8HJcA"
	_, err := jwt.Decode(token)
	assert.NotNil(t, err)
}
func Test_JWT_MultiKey_Decode_Missing_Extra_3_Encoded_With_Default(t *testing.T) {
	jwt := setup_decoder(t)

	// Signed with kid extra_3
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImV4dHJhXzMifQ.eyJyZWFsVXNlcklkIjoiYXNkZjEyMzQiLCJlZmZlY3RpdmVVc2VySWQiOiJhc2RmMTIzNCIsImFjY291bnRJZCI6ImRlZjQ1NiIsImlhdCI6MTY5MTY0NzU5NiwiZXhwIjoxNjkxNzMzOTk2fQ.DDl2aHY0fPYfrBbvTTsnXz1kU-4fqDLsFEs2nPhhKI4Kyg2Na7gQ4Fr4fJHHfmRm4wcsd8R3FPKO3oi5hp9uGyS9t-57MxIVbju76tdx5XZPZUkIcGuTCsu4EAsXGOkjhhfAqwLCWdA-qeP-KSSkQ-wCK6ApRji_glJogtPZKFiZ7ti0VnqbgK9Cjlr-aHFxrjOPFfHiraDnYpg5Jknq1iLkNtqQE5rIvH_J17tMi4XBVwqP8BJZgcD1CeYmyNzpCcV3H7RW3yCNfziILEGd8Xb0-Efh-_ghf98lkhMLUxJi4ToK0VzPcyVrCAvQJ37IaUro_z8-plrbJOfspUrpVrQE-Y73QSQ2dE5rvfp3PK9UNkrgvngmongEXbI44rUH9BU7Sh3O_C_cIxnY0nYEL1hNsDelzP1dAcbHbbU6OkiwZIKThJJzYlSfzoNLUG6-0Q8wQJ7sP-c1Gsi0UcvvGXgR6SYqmoyQBVx7uy4Wzc8X9wVvkJBvs2fy4C8ONqMDTs5-gqHkLfVE99ZUXYsaY9HaFOwgGFROmGAT_v9Og_uh5S0hbEZ4XcD26X8r3b9aFO3Pg9HgYk3a2sZew-MOHCt-elRzHkFCTS4KpiFmbF7lB-SIIBMAA2FUivh7lP6YezEzYc12mAPG4XUgCNCON48zwlfsPKEwojhM9_En--M"
	payload, err := jwt.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, "1a", payload.Customer)
	assert.Equal(t, "2b", payload.RealUser)
	assert.Equal(t, "3c", payload.EffectiveUser)
}