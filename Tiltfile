k8s_yaml([
    'manifests/agg-identity-depl.yaml',
    'manifests/app-identity-depl.yaml',
    'manifests/app-accounting-depl.yaml',
    'manifests/mh-identity-depl.yaml',
    'manifests/mic-projects-depl.yaml',
    'manifests/msg-nats-depl.yaml',
    'manifests/ingress-rules.yaml',
    'manifests/ingress-tls-secrets.yaml',
    'manifests/ingress-traefik-ds.yaml',
    'manifests/secrets.yaml',
])

docker_build('devpies/mh-identity:latest', 'core/identity/handler', target='dev')

docker_build('devpies/agg-identity:latest', 'core/identity/aggregator', target='dev')

docker_build('devpies/app-accounting:latest', 'integrations/freshbooks/application', target='dev')

docker_build('devpies/app-identity:latest', 'core/identity/application', target='dev')

docker_build('devpies/mic-projects:latest', 'core/projects', target='dev')


docker_build('devpies/view-db-identity-migration:latest', 'databases/viewdata')

docker_build('devpies/mic-db-projects-migration:latest', 'databases/projects')

docker_build(' devpies/msg-db-nats-migration:latest', 'databases/nats')
