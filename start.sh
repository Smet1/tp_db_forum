#!/usr/bin/env bash

docker stop db-kek
docker container prune << EOF
Y
EOF
docker image prune << EOF
Y
EOF

docker build -t tp-db-forum .
docker run -p 5000:5000 --name db-kek tp-db-forum:latest