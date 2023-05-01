#!/bin/sh

docker build . -t  node-kafka-demo
helm upgrade node-kafka-demo helm/node-kafka-demo --create-namespace --namespace ashkal --install -f ./helm/values-local.yaml
