k8s_yaml([
    'manifests/app-identity-depl.yaml',
    'manifests/app-accounting-depl.yaml',
    'manifests/agg-identity-depl.yaml',
    'manifests/mic-identity-depl.yaml',
    'manifests/mic-projects-depl.yaml',
    'manifests/ingress-rules.yaml',
    'manifests/ingress-tls-secrets.yaml',
    'manifests/ingress-traefik-ds.yaml',
    'manifests/msg-nats-depl.yaml',
    'manifests/secrets.yaml',
])

docker_build('devpies/app-identity:latest', 'core/identity/application', target='dev')

docker_build('devpies/mic-identity:latest', 'core/identity/microservice', target='dev')

docker_build('devpies/mic-projects:latest', 'core/projects', target='dev')

docker_build('devpies/agg-identity:latest', 'core/identity/aggregator', target='dev')

docker_build('devpies/app-accounting:latest', 'integrations/accounting/application', target='dev')
