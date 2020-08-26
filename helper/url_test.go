package helper

import (
	"testing"

	"gotest.tools/assert"
)

func Test_Url_Domain(t *testing.T) {
	d := Domain("com")
	assert.Assert(t, d == "", d)

	d = Domain("cultureamp.com")
	assert.Assert(t, d == "cultureamp.com", d)
	d = Domain("https://username:password@cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "cultureamp.com", d)

	d = Domain("www.cultureamp.com")
	assert.Assert(t, d == "cultureamp.com", d)
	d = Domain("https://username:password@www.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "cultureamp.com", d)

	d = Domain("customer.cultureamp.com")
	assert.Assert(t, d == "cultureamp.com", d)
	d = Domain("https://username:password@customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "cultureamp.com", d)

	d = Domain("www.customer.cultureamp.com")
	assert.Assert(t, d == "cultureamp.com", d)
	d = Domain("https://username:password@www.customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "cultureamp.com", d)
}

func Test_Url_DomainSuffix(t *testing.T) {
	d := DomainSuffix("com")
	assert.Assert(t, d == "", d)
	d = DomainSuffix("https://username:password@com:443/path?ip=123.345.456")
	assert.Assert(t, d == "", d)

	d = DomainSuffix("cultureamp.com")
	assert.Assert(t, d == "com", d)
	d = DomainSuffix("https://username:password@cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "com", d)

	d = DomainSuffix("www.cultureamp.com")
	assert.Assert(t, d == "com", d)
	d = DomainSuffix("https://username:password@www.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "com", d)

	d = DomainSuffix("customer.cultureamp.com")
	assert.Assert(t, d == "com", d)
	d = DomainSuffix("https://username:password@customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "com", d)

	d = DomainSuffix("www.customer.cultureamp.com")
	assert.Assert(t, d == "com", d)
	d = DomainSuffix("https://username:password@www.customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "com", d)
}

func Test_Url_HasSubdomain(t *testing.T) {
	ok := HasSubdomain("com")
	assert.Assert(t, !ok, ok)
	ok = HasSubdomain("https://username:password@com:443/path?ip=123.345.456")
	assert.Assert(t, !ok, ok)

	ok = HasSubdomain("cultureamp.com")
	assert.Assert(t, !ok, ok)
	ok = HasSubdomain("https://username:password@cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, !ok, ok)

	ok = HasSubdomain("www.cultureamp.com")
	assert.Assert(t, ok, ok)
	ok = HasSubdomain("https://username:password@www.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, ok, ok)

	ok = HasSubdomain("customer.cultureamp.com")
	assert.Assert(t, ok, ok)
	ok = HasSubdomain("https://username:password@customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, ok, ok)

	ok = HasSubdomain("www.customer.cultureamp.com")
	assert.Assert(t, ok, ok)
	ok = HasSubdomain("https://username:password@www.customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, ok, ok)
}

func Test_Url_SubDomain(t *testing.T) {
	d := Subdomain("com")
	assert.Assert(t, d == "", d)
	d = Subdomain("https://username:password@com:443/path?ip=123.345.456")
	assert.Assert(t, d == "", d)

	d = Subdomain("cultureamp.com")
	assert.Assert(t, d == "", d)
	d = Subdomain("https://username:password@cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "", d)

	d = Subdomain("www.cultureamp.com")
	assert.Assert(t, d == "www", d)
	d = Subdomain("https://username:password@www.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "www", d)

	d = Subdomain("customer.cultureamp.com")
	assert.Assert(t, d == "customer", d)
	d = Subdomain("https://username:password@customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "customer", d)

	d = Subdomain("www.customer.cultureamp.com")
	assert.Assert(t, d == "www.customer", d)
	d = Subdomain("https://username:password@www.customer.cultureamp.com:443/path?ip=123.345.456")
	assert.Assert(t, d == "www.customer", d)
}
