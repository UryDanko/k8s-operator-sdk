# k8s-operator-go

K8s operator was created using popular framework [Operator SDK](https://sdk.operatorframework.io/docs/overview/)

## Pre-requisite
* Working K8s cluster. 
* Installed Docker
* Installed Operator SDK [Installation guide](https://sdk.operatorframework.io/docs/installation/)

## Initialization
* Initialize operator with `operator-sdk init --domain example.com --repo gitlab.com/sandbox/simple-server-operator`
* create API with `operator-sdk create api --group sandbox --version v1alpha1 --kind SimpleServer --resource --controller`

## Run the Operator
build the simpleton-server first `make build-server`
Docker images will be pushed to a docker repository so you might want to change the image names if you want to deploy it.
Image names are set in both `server/Makefile` and `server-operator/Makefile`

There are two ways to run the operator:
execute `cd server-operator`
* As Go program outside a cluster. `make install run`
* As a Deployment inside a Kubernetes cluster `make docker-build docker-push deploy`

Operator will take care of Deployment, ConfigMap and Service objects. 
ConfigMap contains just a few dummy variables and mounted to a pod of the deployment. 
Server will just print these variables into the response. 

## Cleanup
Depends on a chosen way of use - do either `make undeploy` or `make uninstall`

