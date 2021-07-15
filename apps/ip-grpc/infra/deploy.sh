#!/bin/zsh

cd $(dirname $0)

eval $(minikube docker-env)

docker build -t ip-grpc:v1.0.2 ./../src/

kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
