#!/bin/sh

docker build . -t  go-kafka-demo
helm upgrade go-kafka-demo helm/go-kafka-demo --create-namespace --namespace ashkal --install -f ./helm/values-local.yaml

