#!/bin/zsh

cd $(dirname $0)

kubectl create namespace observability-sven

kubectl apply -f agent-configmap.yaml

kubectl apply -n observability-sven -f collector-configmap.yaml

kubectl apply -n observability-sven -f collector-service.yaml

kubectl apply -n observability-sven -f collector-deployment.yaml
