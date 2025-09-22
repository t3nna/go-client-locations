# "Microservices with Go" 

## Project overview

### General concept

Design and develop a system to process user locations and provides the ability to search for clients by location (coordinates) and radius. Except that system should provide the ability to calculate the distance traveled by a person in some time range.

The system provides 3 REST endpoints for the backend clients with the following features:

1. Update current user location by the username.
2. Search for users in some location within the provided radius (with pagination).
3. Returns distance traveled by a person within some date/time range. Time range defaults to 1 day. 

    Examples: 

    - For 35.12314, 27.64532 → 39.12355, 27.64538 distance 445km
    - For 35.12314, 27.64532 → 39.12355, 27.64538 → 35.12314, 27.64532 distance is 890km

REST interface and contracts should be designed. 

The system should validate all input data, and respond with the proper status code and message. 

- username - 4-16 symbols (a-zA-Z0-9 symbols are acceptable)
- coordinates - fractional part of a number should be limited by the 8 signs, latitude and longitude should be validated by the regular rules. For example:
    - 35.12314, 27.64532
    - 39.12355, 27.64538
- dates - use ISO 8601 date format (2021-09-02T11:26:18+00:00)



## Installation
The project requires a couple tools to run, most of which are part of many developer's toolchains.

- Docker
- Go
- Tilt
- A local Kubernetes cluster

### MacOS

1. Install Homebrew from [Homebrew's official website](https://brew.sh/)

2. Install Docker for Desktop from [Docker's official website](https://www.docker.com/products/docker-desktop/)

3. Install Minikube from [Minikube's official website](https://minikube.sigs.k8s.io/docs/)

4. Install Tilt from [Tilt's official website](https://tilt.dev/)

5. Install Go on MacOS using Homebrew:
```bash
brew install go
```

6. Make sure [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/) is installed.

## Run

```bash
tilt up
```

In order to use mongo db please add `secret.yaml` with your personal uri

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mongodb
type: Opaque
stringData:
  uri: "paste Uri here"

```

## Monitor

```bash
kubectl get pods
```

or

```bash
minikube dashboard
```

## Web 
Is not ready yet