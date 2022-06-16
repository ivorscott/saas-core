k8s_yaml([
    './manifests/db-admin.yaml',
    './manifests/db-dynamodb.yaml',
    './manifests/db-project.yaml',
    './manifests/db-user.yaml',
    './manifests/nats.yaml',
    './manifests/traefik.yaml',
    './manifests/ingress-rules.yaml',
    './manifests/mic-identity.yaml',
    './manifests/mic-tenant.yaml',
    './manifests/mic-registration.yaml',
    './manifests/mic-admin.yaml',
    './manifests/mic-project.yaml',
    './manifests/mic-user.yaml',
    './manifests/secrets.yaml'
])

docker_build('identity:latest', '.' ,dockerfile = 'deploy/identity.dockerfile')
docker_build('tenant:latest', '.', dockerfile = 'deploy/tenant.dockerfile')
docker_build('registration:latest', '.' ,dockerfile = 'deploy/registration.dockerfile')
docker_build('admin:latest', '.', dockerfile = 'deploy/admin.dockerfile')
docker_build('project:latest', '.', dockerfile = 'deploy/project.dockerfile')
docker_build('user:latest', '.', dockerfile = 'deploy/user.dockerfile')