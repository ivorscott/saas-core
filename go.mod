module github.com/devpies/core

go 1.18

require (
	github.com/Masterminds/squirrel v1.5.3
	github.com/alexedwards/scs/pgxstore v0.0.0-20220528130143-d93ace5be94b
	github.com/alexedwards/scs/v2 v2.5.0
	github.com/ardanlabs/conf v1.5.0
	github.com/aws/aws-sdk-go v1.44.17
	github.com/aws/aws-sdk-go-v2/config v1.15.7
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.15.5
	github.com/devpies/core/pkg/log v0.0.0-20220519154201-382d963b2ca0
	github.com/devpies/core/pkg/web v0.0.0-20220519154201-382d963b2ca0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.6
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.21.0
)

require (
	github.com/alexedwards/scs/postgresstore v0.0.0-20220528130143-d93ace5be94b // indirect
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
	github.com/golang-migrate/migrate/v4 v4.15.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	github.com/jackc/pgx/v4 v4.16.1 // indirect
	github.com/jackc/puddle v1.2.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/devpies/core/pkg/web v0.0.0-20220519154201-382d963b2ca0 => ./pkg/web

replace github.com/devpies/core/pkg/log v0.0.0-20220519154201-382d963b2ca0 => ./pkg/log
