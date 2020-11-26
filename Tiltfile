k8s_yaml([
    '__infra__/app-identity-depl.yaml',
    '__infra__/com-identity-depl.yaml',
    '__infra__/ingress-rules.yaml',
    '__infra__/ingress-traefik-ds.yaml',
    '__infra__/msg-nats-depl.yaml',
    '__infra__/secrets.yaml',
    '__infra__/view-db-identity-depl.yaml'
])

docker_build('devpies/client-app-identity', 'identity/application', target='dev')
docker_build('devpies/client-com-identity', 'identity/component', target='dev')
