# k8s-operator-go

## Pre-requisite
* working K8s cluster
* installed docker

## Initialization
* Initialize operator with `operator-sdk init --domain example.com --repo gitlab.com/sandbox/simple-server-operator`
* create API with `operator-sdk create api --group sandbox --version v1alpha1 --kind SimpleServer --resource --controller`

## Run the Operator
build the simpleton-server first `make build-server`

There are two ways to run the operator:
execute `cd server-operator`
* As Go program outside a cluster. `make install run`
* As a Deployment inside a Kubernetes cluster `make docker-build docker-push deploy`

## Cleanup
Depends on a chosen way of use - do either `make undeploy` or `make uninstall`

