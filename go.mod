module github.com/cultureamp/glamplify

go 1.15

replace github.com/aws/aws-xray-sdk-go v1.6.0 => github.com/aws/aws-xray-sdk-go v1.6.1-0.20211110224843-1f272e4024a5

require (
	github.com/DataDog/datadog-go v4.8.3+incompatible // indirect
	github.com/DataDog/datadog-lambda-go v1.3.0
	github.com/DataDog/sketches-go v1.2.1 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.42.12
	github.com/aws/aws-xray-sdk-go v1.6.1-0.20211110224843-1f272e4024a5
	github.com/bobesa/go-domain-util v0.0.0-20190911083921-4033b5f7dd89
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/getsentry/sentry-go v0.11.0
	github.com/go-errors/errors v1.1.1
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/google/uuid v1.3.0
	github.com/gookit/color v1.3.6
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.0.0-20200921180117-858c6e7e6b7e
	github.com/sony/gobreaker v0.5.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tinylib/msgp v1.1.6 // indirect
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1 // indirect
	google.golang.org/grpc v1.42.0 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.34.0
)

exclude (
	github.com/kataras/iris/v12 v12.1.8
	github.com/labstack/echo/v4 v4.1.11
)
