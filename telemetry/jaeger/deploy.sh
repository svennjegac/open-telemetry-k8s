#!/bin/zsh

cd $(dirname $0)

kubectl create namespace observability-sven

kubectl create -f operator/custom-resource-definition.yaml

kubectl create -n observability-sven -f operator/service-account.yaml

kubectl create -n observability-sven -f operator/role.yaml

kubectl create -n observability-sven -f operator/role-binding.yaml

kubectl create -n observability-sven -f operator/operator.yaml

kubectl create -n observability-sven -f operator/cluster-role.yaml

kubectl create -n observability-sven -f operator/cluster-role-binding.yaml

kubectl apply -n observability-sven -f deploy-simple/all-in-one.yaml
