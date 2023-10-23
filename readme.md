# kube-ipam

## ⚠️ Disclaimer

This project is a prototype and is intended for educational purposes only. It is not designed to be used in production environments.

## Overview

kube-ipam is an IP address manager designed to allocate external IPs to Kubernetes services of type LoadBalancer.

## Installation

### Docker

To build the Docker image, navigate to the `config` directory and execute:

```bash
docker build -t kube-ipam .
```

### Kubernetes

Apply the Kubernetes manifest from the `config` directory:

```bash
kubectl apply -f kubernetes_manifest.yaml
```

## Prerequisites

- Go (version 1.21 or higher)
- Docker (if building a Docker image)
- Kubernetes cluster (if deploying to Kubernetes)
- `kubectl` configured to interact with your Kubernetes cluster

## Usage

### Command-Line Flags

- `--kubeconfig` : Absolute path to the kubeconfig file. If not provided, will take config from home directory (`~/.kube/config`). If run in a cluster, will use service account token.
- `--cidr` : CIDR block for LoadBalancer IPs allocation.

### Examples

#### Running locally

```bash
./kube-ipam --cidr=192.168.0.0/24
```

#### Running in Kubernetes

Change the `--cidr` flag value in the `kubernetes_manifest.yaml` if necessary.

```yaml
args:
  - --cidr=192.168.10.0/24
```

Then re-apply the manifest.

```bash
kubectl apply -f kubernetes_manifest.yaml
```

## Project Structure

- `cmd/`: Contains the main function entry for the project.
- `internal/`: Houses utility functions, IPPool implementation, and Kubernetes client code.
- `config/`: Contains Dockerfile and Kubernetes manifest files for deployment.

## Development

To work on this project locally, follow these steps:

1. Clone the repo:

    ```bash
    git clone https://github.com/blakelead/kube-ipam.git
    ```

2. Navigate to the root directory and initialize the Go module:

    ```bash
    go mod init kube-ipam
    ```

3. Download dependencies:

    ```bash
    go mod tidy
    ```

4. Build the project:

    ```bash
    go build -o kube-ipam ./cmd
    ```

Refer to the "Installation" section for details on how to build a Docker image or deploy to a Kubernetes cluster.

## License

This project is licensed under the MIT License.
