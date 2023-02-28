# App for uploading artifacts

WORK IN PROGRESS

Current status: MVP

[![Release published](https://github.com/privatepages/privatepages/actions/workflows/release.yml/badge.svg)](https://github.com/privatepages/privatepages/actions/workflows/release.yml)

[![Code checks](https://github.com/privatepages/privatepages/actions/workflows/audit.yml/badge.svg)](https://github.com/privatepages/privatepages/actions/workflows/audit.yml)


If you want to deploy privatepages at your k8s - use that helm chart [https://github.com/privatepages/privatepages-chart](https://github.com/privatepages/privatepages-chart)

If you want to use privatepages in your Actions - use that Action [https://github.com/privatepages/upload-action](https://github.com/privatepages/upload-action)

## Develop

### In minikube with skaffold (linux) (without oauth proxy)

    minikube start
    eval $(minikube docker-env)
    kubectl config use-context minikube
    kubectl create ns privatepages-api
    kubectl config set-context --current --namespace=privatepages-api

    # APP
    cp skaffold/secrets-example.yaml skaffold/secrets.yaml
    # and edit secrets
    kubectl apply -n privatepages-api -f skaffold/secrets.yaml
    kubectl apply -n privatepages-api -f skaffold/pvc.yaml
    kubectl apply -n privatepages-api -f skaffold/configmap.yaml
    skaffold dev

### Local (linux) (without nginx and oauth proxy)

    mkdir /tmp/tests
    cd ./src
    go mod tidy
    export API_SECRET=secret
    export HTTP_LISTEN=:8080
    export LOG_LEVEL=debug
    export ARTIFACT_STORAGE_PATH=/tmp/tests
    go run cmd/app/main.go

### Run on docker (without nginx and oauth proxy)

    docker build -t pp .
    docker run -it --rm \
        -e ARTIFACT_STORAGE_PATH=/data \
        -e LOG_LEVEL=debug \
        -e HTTP_LISTEN=:8080 \
        -e API_SECRET=secret \
        -v /tmp/tests:/data \
        -p 8080:8080 \
        pp

## Working with artifacts

### Upload artifact

Try use upload-test.html

OR

    curl 127.0.0.1:8080/upload -F 'file=@/path/to/local/file.txt' -F 'artifactname=my-folder-name' -F 'token=secret'

### Remove artifact

    curl -X POST http://localhost:8080/remove -H "Content-Type: multipart/form-data" -F 'artifactname=my-project33' -F 'token=secsret'


## To Do

* promhttp not collecting metrics for all routes (github.com/zsais/go-gin-prometheus?)
* gin access log to json, use LOG_LEVEL, GIN_MODE=release | gin.SetMode(gin.ReleaseMode)
* swagger
* extensive readme documents in each repo
* chart "how to" with any solutions like cert-manager, nginx-ingress, etc.
