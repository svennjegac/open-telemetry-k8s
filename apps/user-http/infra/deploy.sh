#!/bin/zsh

cd $(dirname $0)

eval $(minikube docker-env)

docker build -t user-http:v1.0.23 ./../src/

kubectl apply -f deployment.yaml
