#!/bin/bash
# This file updates database manifest files that depend on the host path.
# It ensures the path used is defined correctly. By default, this project places
# database files in the 'data' folder at the project root.

pwd=`pwd`
# escape path string
path=${pwd//\//\\\/}

sed -i.bakup 's/path:.*/path: '$path'\/data\/admin/' manifests/db-admin.yaml
sed -i.bakup 's/path:.*/path: '$path'\/data\/nats/' manifests/db-nats.yaml
sed -i.bakup 's/path:.*/path: '$path'\/data\/project/' manifests/db-project.yaml
sed -i.bakup 's/path:.*/path: '$path'\/data\/subscription/' manifests/db-subscription.yaml
sed -i.bakup 's/path:.*/path: '$path'\/data\/user/' manifests/db-user.yaml

rm manifests/db-*.yaml.bakup