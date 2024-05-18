# stress

A simple tool for generating CPU workload and memory allocation. Useful for
creating stress test and benchmarking.

## Compiling

```bash
go build -o stress
./stress -time 30

# Or compile for other platforms. Copy the output '.exe' file to run on Windows
GOOS=windows GOARCH=amd64 go build -o stress.exe
```

If you don't have golang installed locally, use the docker image to compile the executable.

```bash
# assuming you're in the project directory
docker run --rm -it -v "./:/go" -e "GOOS=windows" -e "GOARCH=amd64" golang:1.22.3-alpine go build -o stress.exe
```

## Usage

```
$ stress -time 30

Performing CPU work for 30 seconds
CPU score: 9801
Workload generation complete.
```

See available flags below.

#### Generate max CPU workload for 30 seconds

```bash
$ stress -time 30
```

#### Allocate 20GB of memory, or error if system doesn't have enough available memory

```bash
$ stress -mem 20G
```

#### Allocate 20GB of memory, and generate CPU workload for 60 seconds

```bash
$ stress -time 60 -mem 20G
```

#### Locally

```bash
git clone https://github.com/karunsiri/stress
cd stress
go run main.go -time 30
```

#### Docker

```bash
docker run --rm -it karunsiri/stress [flags]
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
