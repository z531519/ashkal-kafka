#!/bin/sh

docker build . -t  java-kafka-demo
helm upgrade java-kafka-demo helm/java-kafka-demo --namespace ashkal --install -f ./helm/values-local.yaml
