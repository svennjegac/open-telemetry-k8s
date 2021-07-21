#!/bin/zsh

cd $(dirname $0)

eval $(minikube docker-env)

docker build -t user-events-kafka:v1.0.7 ./../src/

kubectl apply -f deployment.yaml
