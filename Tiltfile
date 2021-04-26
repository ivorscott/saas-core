k8s_yaml([
    'manifests/mic-users-depl.yaml',
    'manifests/mic-projects-depl.yaml',
    'manifests/msg-nats-depl.yaml',
    'manifests/ingress-rules.yaml',
    'manifests/ingress-tls-secrets.yaml',
    'manifests/ingress-traefik-ds.yaml',
    'manifests/secrets.yaml',
])

docker_build('devpies/mic-users:latest', 'core/users',
build_args=read_json('.gitpass'), target='dev')

docker_build('devpies/mic-projects:latest', 'core/projects',
build_args=read_json('.gitpass'), target='dev')
