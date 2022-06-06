k8s_yaml([
    './manifests/db-admin.yaml',
    './manifests/db-dynamodb.yaml',
    './manifests/nats.yaml',
    './manifests/traefik.yaml',
    './manifests/ingress-rules.yaml',
    './manifests/mic-registration.yaml',
    './manifests/mic-admin.yaml',
    './manifests/secrets.yaml'
])

docker_build('registration:latest', '.' ,dockerfile = 'deploy/registration.dockerfile')
docker_build('admin:latest', '.', dockerfile = 'deploy/admin.dockerfile')