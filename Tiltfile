k8s_yaml([
    '__k8s__/app-identity-depl.yaml',
    '__k8s__/app-accounting-depl.yaml',
    '__k8s__/agg-identity-depl.yaml',
    '__k8s__/mic-identity-depl.yaml',
    '__k8s__/mic-projects-depl.yaml',
    '__k8s__/ingress-rules.yaml',
    '__k8s__/ingress-tls-secrets.yaml',
    '__k8s__/ingress-traefik-ds.yaml',
    '__k8s__/msg-nats-depl.yaml',
    '__k8s__/secrets.yaml',
])

docker_build('devpies/client-app-identity', 'identity/application', target='dev')

docker_build('devpies/client-mic-identity', 'identity/microservice', target='dev')

docker_build('devpies/client-mic-projects', 'projects', target='dev')

docker_build('devpies/client-agg-identity', 'identity/aggregator', target='dev')

docker_build('devpies/client-app-accounting', 'accounting/application', target='dev')
