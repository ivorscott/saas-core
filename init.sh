#!/bin/bash
# This file initializes the project.
pwd

sed -i 's/path: .*/path: $(pwd)/data/admin' ./manifests/db-admin.yaml