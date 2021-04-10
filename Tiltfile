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

docker_build('devpies/app-identity:latest', 'identity/application', target='dev')

docker_build('devpies/mic-identity:latest', 'identity/microservice', target='dev')

docker_build('devpies/mic-projects:latest', 'projects', target='dev')

docker_build('devpies/agg-identity:latest', 'identity/aggregator', target='dev')

docker_build('devpies/app-accounting:latest', 'accounting/application', target='dev')
