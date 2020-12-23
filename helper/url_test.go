package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Url_Domain(t *testing.T) {
	d := Domain("com")
	assert.Empty(t, d)

	d = Domain("cultureamp.com")
	assert.Equal(t, "cultureamp.com", d)
	d = Domain("https://username:password@cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "cultureamp.com", d)

	d = Domain("www.cultureamp.com")
	assert.Equal(t, "cultureamp.com", d)
	d = Domain("https://username:password@www.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "cultureamp.com", d)

	d = Domain("customer.cultureamp.com")
	assert.Equal(t, "cultureamp.com", d)
	d = Domain("https://username:password@customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "cultureamp.com", d)

	d = Domain("www.customer.cultureamp.com")
	assert.Equal(t, "cultureamp.com", d)
	d = Domain("https://username:password@www.customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "cultureamp.com", d)
}

func Test_Url_DomainSuffix(t *testing.T) {
	d := DomainSuffix("com")
	assert.Empty(t, d)
	d = DomainSuffix("https://username:password@com:443/path?ip=123.345.456")
	assert.Empty(t, d)

	d = DomainSuffix("cultureamp.com")
	assert.Equal(t, "com", d)
	d = DomainSuffix("https://username:password@cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "com", d)

	d = DomainSuffix("www.cultureamp.com")
	assert.Equal(t, "com", d)
	d = DomainSuffix("https://username:password@www.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "com", d)

	d = DomainSuffix("customer.cultureamp.com")
	assert.Equal(t, "com", d)
	d = DomainSuffix("https://username:password@customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "com", d)

	d = DomainSuffix("www.customer.cultureamp.com")
	assert.Equal(t, "com", d)
	d = DomainSuffix("https://username:password@www.customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "com", d)
}

func Test_Url_HasSubdomain(t *testing.T) {
	ok := HasSubdomain("com")
	assert.False(t, ok)
	ok = HasSubdomain("https://username:password@com:443/path?ip=123.345.456")
	assert.False(t, ok)

	ok = HasSubdomain("cultureamp.com")
	assert.False(t, ok)
	ok = HasSubdomain("https://username:password@cultureamp.com:443/path?ip=123.345.456")
	assert.False(t, ok)

	ok = HasSubdomain("www.cultureamp.com")
	assert.True(t, ok)
	ok = HasSubdomain("https://username:password@www.cultureamp.com:443/path?ip=123.345.456")
	assert.True(t, ok)

	ok = HasSubdomain("customer.cultureamp.com")
	assert.True(t, ok)
	ok = HasSubdomain("https://username:password@customer.cultureamp.com:443/path?ip=123.345.456")
	assert.True(t, ok)

	ok = HasSubdomain("www.customer.cultureamp.com")
	assert.True(t, ok)
	ok = HasSubdomain("https://username:password@www.customer.cultureamp.com:443/path?ip=123.345.456")
	assert.True(t, ok)
}

func Test_Url_SubDomain(t *testing.T) {
	d := Subdomain("com")
	assert.Empty(t, d)
	d = Subdomain("https://username:password@com:443/path?ip=123.345.456")
	assert.Empty(t, d)

	d = Subdomain("cultureamp.com")
	assert.Empty(t, d)
	d = Subdomain("https://username:password@cultureamp.com:443/path?ip=123.345.456")
	assert.Empty(t, d)

	d = Subdomain("www.cultureamp.com")
	assert.Equal(t, "www", d)
	d = Subdomain("https://username:password@www.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "www", d)

	d = Subdomain("customer.cultureamp.com")
	assert.Equal(t, "customer", d)
	d = Subdomain("https://username:password@customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "customer", d)

	d = Subdomain("www.customer.cultureamp.com")
	assert.Equal(t, "www.customer", d)
	d = Subdomain("https://username:password@www.customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Equal(t, "www.customer", d)
}
