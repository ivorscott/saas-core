k8s_yaml([
    '__infra__/app-identity-depl.yaml',
    '__infra__/agg-identity-depl.yaml',
    '__infra__/mic-identity-depl.yaml',
    '__infra__/ingress-rules.yaml',
    '__infra__/ingress-traefik-ds.yaml',
    '__infra__/msg-nats-depl.yaml',
    '__infra__/secrets.yaml',
])

docker_build('devpies/client-app-identity', 'identity/application', target='dev')

docker_build('devpies/client-mic-identity', 'identity/microservice', target='dev')

docker_build('devpies/client-agg-identity', 'identity/aggregator', target='dev')