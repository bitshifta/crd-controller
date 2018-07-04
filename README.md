# Pod Counter: Controller for Kubernetes


## What
PodCounter is a controller that tracks all the pods that are created in a cluster.
PodCounter uses a Kubernetes Custom Resource Definition(CRD) to store the number of running pods, as well as a historical count of pods that have been scheduled.

## Why
Why not?
This is an exercise to see how the Kubernetes API can be extended and how you store information in the etcd masters

## Compilation

```bash
go build *.go
```

## Running

Running the controller is relatively straightforward.

The `manifests` folder contains the k8s manifests that are needed to run the controller, ranging from the Custom Resource Definition itself, the RBAC permissions for the PodCounter and the deployment of the controller code.

To run the code locally, using a configuration file, run it like:

```bash
kubectl apply -f manifests/crd.yaml
./controller --kubeconfig ~/my-config-location/kube.conf
```

To run the controller in a cluster via deployment, build the docker image like so:

```bash
# Replace `user` with your own docker handle or registry name
docker build -t user/controller:latest .

# Push the image to the registry
docker push user/controller:latest

# Update the deployment yaml with the registry name and image tag
-> manifests/deployment.yaml

# Apply the manifests using kubectl
kubectl apply -f manifests/
```
## Hang on a second, how do I see the numbers?

Using your favourite browser, or curl if you have it handy.

```bash
# First make sure your cluster's api is reachable
$ kubectl proxy
Starting to serve on 127.0.0.1:8001

# Then curl this endpoint, or open it in your browser
curl http://127.0.0.1:8001/apis/khalilt.com/v1/namespaces/default/podcounters/cluster-counter
```
# TODO

- [x] write basic controller
- [x] store the metrics in CRD
- [ ] write end-to-end tests
- [ ] make the controller HA
