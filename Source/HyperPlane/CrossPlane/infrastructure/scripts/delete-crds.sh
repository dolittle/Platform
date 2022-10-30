#!/bin/bash
kubectl get crd -o name | grep crossplane.io | xargs kubectl delete