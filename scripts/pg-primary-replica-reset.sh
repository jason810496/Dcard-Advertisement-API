#!/bin/bash

# ==== Color ====
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "$GREEN Stop all containers$NC"
docker compose -f docker-compose-pg-primary-replica.yaml down

if [ -d "./stateful_volumes/dev/primary" ]; then
    echo -e "$RED Remove$NC old data from$RED primary-replica/primary$NC"
    rm -r ./stateful_volumes/dev/primary
fi

if [ -d "./stateful_volumes/dev/replica" ]; then
    echo -e "$RED Remove$NC old data from$RED primary-replica/replica$NC"
    rm -r ./stateful_volumes/dev/replica
fi


if [ -d "./stateful_volumes/dev/primary_copy" ]; then
    echo -e "$RED Remove$NC old data from$RED primary-replica/primary_copy$NC"
    rm -r ./stateful_volumes/dev/primary_copy
fi


# echo -e "$GREEN Remove old data$NC"
# if [ -d "./stateful_volumes/dev/primary-replica" ]; then
#     echo -e "$RED Remove$NC old data from$RED primary-replica$NC"
#     rm -r ./stateful_volumes/dev/primary-replica
# fi

echo -e "$GREEN Start all containers$NC"
docker compose -f docker-compose-pg-primary-replica.yaml up -d
