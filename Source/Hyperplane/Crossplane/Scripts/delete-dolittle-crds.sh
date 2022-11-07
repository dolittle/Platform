#!/bin/bash
kubectl get crd -o name | grep dolittle.io | xargs kubectl delete