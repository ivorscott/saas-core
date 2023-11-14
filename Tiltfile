k8s_yaml([
    './manifests/configmap.yaml',
    './manifests/secrets.yaml',
    './manifests/db-nats.yaml',
    './manifests/db-admin.yaml',
    './manifests/db-project.yaml',
    './manifests/db-project-test.yaml',
    './manifests/db-user.yaml',
    './manifests/db-subscription.yaml',
    './manifests/traefik-crd.yaml',
    './manifests/traefik-depl.yaml',
    './manifests/traefik-svc.yaml',
    './manifests/traefik-headers.yaml',
    './manifests/mic-tenant.yaml',
    './manifests/mic-registration.yaml',
    './manifests/mic-admin.yaml',
    './manifests/mic-project.yaml',
    './manifests/mic-user.yaml',
    './manifests/mic-subscription.yaml'
])

docker_build('subscription:latest', '.', dockerfile = 'deploy/subscription.dockerfile')
docker_build('tenant:latest', '.', dockerfile = 'deploy/tenant.dockerfile')
docker_build('registration:latest', '.' ,dockerfile = 'deploy/registration.dockerfile')
docker_build('admin:latest', '.', dockerfile = 'deploy/admin.dockerfile')
docker_build('project:latest', '.', dockerfile = 'deploy/project.dockerfile')
docker_build('user:latest', '.', dockerfile = 'deploy/user.dockerfile')