k8s_yaml([
    './manifests/db-admin.yaml',
    './manifests/db-dynamodb.yaml',
    './manifests/db-project.yaml',
    './manifests/nats.yaml',
    './manifests/traefik.yaml',
    './manifests/ingress-rules.yaml',
    './manifests/mic-user.yaml',
    './manifests/mic-tenant.yaml',
    './manifests/mic-registration.yaml',
    './manifests/mic-admin.yaml',
    './manifests/mic-project.yaml',
    './manifests/secrets.yaml'
])

docker_build('user:latest', '.' ,dockerfile = 'deploy/user.dockerfile')
docker_build('tenant:latest', '.', dockerfile = 'deploy/tenant.dockerfile')
docker_build('registration:latest', '.' ,dockerfile = 'deploy/registration.dockerfile')
docker_build('admin:latest', '.', dockerfile = 'deploy/admin.dockerfile')
docker_build('project:latest', '.', dockerfile = 'deploy/project.dockerfile')