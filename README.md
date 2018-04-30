# What's this?
A tiny go program to check whether a kubernetes service has one or more available
endpoints. Keeps looping between start time and heat death of universe until
one or more ready addresses are found.

This was created because there's currently no good way of saying "don't stand this
Pod until service X and Y is available". Now there is: run this as an init container.

# Requirements
For local usage, the application requires a kubectl config file in $HOME/.kube/config (overridable using flag -kubeconfig)

To build, `glide` is required.

# Build and usage
Docker:
```
docker build -t k8s-service-availability .
docker run -v ~/.kube:/.kube --rm k8s-service-availability
```

Locally
```
glide up -v
go build -o ./k8stest .
./k8stest -kubeconfig ~/.kube/nonstandardconfiglocation
```

# Usage
Check for endpoint readiness for service `servicename` in namespace `namespacename`

./k8stest -namespace namespacename -service servicename

# End notes
I have no idea what I'm doing.