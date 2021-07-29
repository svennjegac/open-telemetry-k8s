#!/bin/zsh

cd $(dirname $0)

eval $(minikube docker-env)

docker build -t ip-grpc:v1.0.10 ./../src/

kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
