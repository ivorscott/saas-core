# SaaS-Core

This project is a part of "SaaS app in 30 days" - _Proof of Concept_

## Getting Started
You can print the purpose of each makefile target using `make`.

```bash
> make
admin-client      Run admin frontend with live reload
admin-api         Run admin backend with live reload
```

>__TIP__
> 
>Enable bash-completion of the makefile targets. Add this in your `~/.bash_profile` file or `~/.bashrc` file.
> ```bash
> complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make
> ```

## Environment Variables

The `.env` file contains variables for all programs. Using `make` automatically references these values.
Program requirements are also documented in help text.
```bash
> go run ./cmd/admin-client -h
Usage: admin-client [options] [arguments]

OPTIONS
  --web-debug/$ADMIN_WEB_DEBUG                                      <string>    (default: localhost:6060)
  --web-production/$ADMIN_WEB_PRODUCTION                            <bool>      (default: false)
  --web-read-timeout/$ADMIN_WEB_READ_TIMEOUT                        <duration>  (default: 5s)
  --web-write-timeout/$ADMIN_WEB_WRITE_TIMEOUT                      <duration>  (default: 5s)
  --web-shutdown-timeout/$ADMIN_WEB_SHUTDOWN_TIMEOUT                <duration>  (default: 5s)
  --web-app-frontend/$ADMIN_WEB_APP_FRONTEND                        <string>    (default: localhost:4000)
  --web-app-backend/$ADMIN_WEB_APP_BACKEND                          <string>    (default: localhost:4001)
  --cognito-app-client-id/$ADMIN_COGNITO_APP_CLIENT_ID              <string>    (default: none)
  --cognito-user-pool-client-id/$ADMIN_COGNITO_USER_POOL_CLIENT_ID  <string>    (default: none)
  --stripe-key/$ADMIN_STRIPE_KEY                                    <string>    (default: none)
  --stripe-secret/$ADMIN_STRIPE_SECRET                              <string>    (default: none)
  --help/-h                                                         
  display this help message
```

> __TIP__  
> 
> Export `.env` file variables before running go binaries to avoid using CLI flags:  
> ```bash
> export $(grep -v '^#' .env | xargs)
> go run ./cmd/{PROGRAM}
>```

## Conventions
Every program has its own README under `cmd`. Common packages are kept in `pkg`. Most logic for apps and services live under `internal`.
Apps and services __MUST__ use a 3 layered architecture: _handler, service, and repository layers._ Here's an [example](https://github.com/devpies/employee-service/blob/bbeebf3d6a301750c2bdef999bb4e8c9bae9467b/cmd/employee/main.go#L81-L83).
