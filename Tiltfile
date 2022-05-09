k8s_yaml([
    'manifests/mic-users-depl.yaml',
    'manifests/mic-projects-depl.yaml',
    'manifests/msg-nats-depl.yaml',
    'manifests/ingress-traefik-ds.yaml',
    'manifests/ingress-rules.yaml',
    'manifests/secrets.yaml',
])

docker_build('devpies/users:latest', 'users',
build_args=read_json('.gitpass'), target='dev')

docker_build('devpies/projects:latest', 'projects',
build_args=read_json('.gitpass'), target='dev')
