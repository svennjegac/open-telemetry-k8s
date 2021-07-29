#!/bin/zsh

cd $(dirname $0)

eval $(minikube docker-env)

docker build -t user-http:v2.0.10 ./../src/

kubectl apply -f deployment.yaml
