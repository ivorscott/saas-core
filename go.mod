module github.com/devpies/core

go 1.18

require (
	github.com/alexedwards/scs/v2 v2.5.0
	github.com/ardanlabs/conf v1.5.0
	github.com/aws/aws-sdk-go v1.44.17
	github.com/aws/aws-sdk-go-v2/config v1.15.7
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.15.5
	github.com/devpies/core/pkg/log v0.0.0-20220519154201-382d963b2ca0
	github.com/devpies/core/pkg/web v0.0.0-20220519154201-382d963b2ca0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
	go.uber.org/zap v1.21.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.16.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.6 // indirect
	github.com/aws/smithy-go v1.11.2 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/devpies/core/pkg/web v0.0.0-20220519154201-382d963b2ca0 => ./pkg/web

replace github.com/devpies/core/pkg/log v0.0.0-20220519154201-382d963b2ca0 => ./pkg/log
