# Setup Guide

## Requirements

Tested on a m1 mac . It should work on linux as well.

- aws account
- install kubernetes
- install [terraform](https://www.terraform.io/)
- install [go v1.21 or higher](https://go.dev/doc/install)
- install [tilt](https://tilt.dev/)
- install [mkcert](https://github.com/FiloSottile/mkcert)
- install [mockery](https://github.com/vektra/mockery)
- install [pgcli](https://www.pgcli.com/)
- install [golangci-lint](https://github.com/golangci/golangci-lint)
- install [go-migrate](https://github.com/golang-migrate/migrate)
- [saas-infra resources](https://github.com/devpies/saas-infra/tree/main/local/saas)

## Instructions 

1. Checkout `saas-infra` and deploy the `local` infrastructure.
   - You will need to supply a valid email for the _SaaS provider admin user_. This user is used to
   login to the admin web app.
   ![](img/admin-webapp.png)
2. Use terraform output values for this repository's `.env` file.
3. Copy `.env.sample` in the project root and create your own `.env` file.
4. Copy `./manifests/secrets.sample.yaml` and create your own `./manifests/secrets.yaml` file.
5. Generate valid tls self-signed certificates: `mkcert devpie.local "*.devpie.local" localhost 127.0.0.1 ::1`
6. Generate the `tls-secret` yaml for traefik with the certificate values: 
   ```
   kubectl create secret generic tls-secret --from-file=tls.crt=./devpie.local.pem --from-file=tls.key=./devpie.local-key.pem -o yaml 
   ```
   Then add the contents to the bottom of your secrets.yaml file.
7. Modify your hosts file:
   ```bash
    ##
    # Host Database
    #
    # localhost is used to configure the loopback interface
    # when the system is booting.  Do not change this entry.
    ##
    127.0.0.1       localhost devpie.local admin.devpie.local api.devpie.local 
    ```
   
8. Start containers: `tilt up`

9. In another terminal, port forward the traefik ports: `make ports`
10. In another terminal, deploy ingress routes: `make routes`
11. Navigate to http://localhost:8080/dashboard/#/http/routers. You should see `4` tls terminated routers.

![](img/traefik.png)

## Getting Help
If you need help or have questions create an issue. Alternatively, you can join our [discord server](https://discord.gg/MeKKvHBKQG) 
and reach out there.

## Contributions
If you have ideas on automating this setup feel free to submit a PR. 

### Gotchas
golangci-lint returns error on new mac versions.
https://github.com/golangci/golangci-lint/discussions/3327

FIX: `brew install diffutils`