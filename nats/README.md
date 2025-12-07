# Nats

Learning more about NATS.

- Using Kind to create a k8s cluster.
- Using Helm charts to install NATS.

Useful commands

```bash 
kind create cluster --config kind-config.yaml --name nats-cluster

kubectl create namespace nats

helm install nats nats/nats --version 2.12.2 -f nats-values.yaml --namespace nats

helm upgrade nats nats/nats -f nats-values.yaml

# To list the pods and all nodes that each one belongs to:
kubectl get pods -n nats -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,NODE:.spec.nodeName,IP:.status.podIP,RESTARTS:.status.containerStatuses[0].restartCount
```

