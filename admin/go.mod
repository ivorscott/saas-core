module github.com/devpies/core/admin

go 1.18

require (
	github.com/alexedwards/scs/v2 v2.5.0
	github.com/ardanlabs/conf v1.5.0
	github.com/devpies/core/pkg/log v0.0.0-20220509162407-c76b90e617d1
	github.com/devpies/core/pkg/web v0.0.0
	github.com/go-chi/chi/v5 v5.0.7
	go.uber.org/zap v1.21.0
)

require (
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/devpies/core/pkg/web v0.0.0 => ../pkg/web
