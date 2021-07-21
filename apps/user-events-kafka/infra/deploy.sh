#!/bin/zsh

cd $(dirname $0)

eval $(minikube docker-env)

docker build -t user-events-kafka:v1.0.8 ./../src/

kubectl apply -f deployment.yaml
