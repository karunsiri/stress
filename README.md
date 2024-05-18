# stress

A simple tool for generating CPU workload and memory allocation. Useful for
creating stress test and benchmarking.

## Usage

### Docker

See available flags below.
```bash
docker run --rm -it karunsiri/stress [flags]
```


Generate max CPU workload for 30 seconds

```bash
docker run --rm -it karunsiri/stress -time 30
```

Allocate 20GB of memory, or error if system doesn't have enough available memory

```bash
docker run --rm -it karunsiri/stress -mem 20G
```

Allocate 20GB of memory, and generate CPU workload for 60 seconds

```bash
docker run --rm -it karunsiri/stress -time 60 -mem 20G
```

### Kubernetes

**Use case:** Testing for cluster alerting rules.
Example Deployment.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: workload
  labels:
    app: workload
spec:
  replicas: 3
  selector:
    matchLabels:
      app: workload
  template:
    metadata:
      labels:
        app: workload
    spec:
      containers:
      - name: stress
        image: karunsiri/stress
        args:
          - '-time=60'
          - '-mem=1G'
```

## Available flags:

- `-help`: Show usage
- `-mem`: Amount of memory to allocate. Support single number, or number with suffixes 'K', 'M', 'G'.
  For example, `400` means 400 bytes. `400M` means 400MB, and so on.
- `-time`: Duration of CPU workload in seconds.
