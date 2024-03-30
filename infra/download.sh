#!/bin/bash

BASE_URL=https://raw.githubusercontent.com/terraform-google-modules/terraform-google-kubernetes-engine/v30.2.0/examples/simple_autopilot_public/
FILES=(main.tf network.tf outputs.tf variables.tf versions.tf README.md)

for FILE in ${FILES[@]}; do
    echo "Downloading $FILE"
    curl -s -O $BASE_URL$FILE
done