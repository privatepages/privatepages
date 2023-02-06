# REST API for uploading artifacts
![code checks](https://github.com/privatepages/privatepages/actions/workflows/audit.yml/badge.svg)

## Develop

### In minikube with skaffold (linux)

    minikube start
    eval $(minikube docker-env)
    kubectl config use-context minikube
    kubectl create ns privatepages-api
    kubectl config set-context --current --namespace=privatepages-api

    # APP
    cp skaffold/secrets-example.yaml skaffold/secrets.yaml
    # and edit secrets
    kubectl apply -n privatepages-api -f skaffold/secrets.yaml
    skaffold dev

### Local (linux)

    cd ./src
    export API_SECRET=test
    export HTTP_LISTEN=:8080
    export LOG_LEVEL=debug
    go run cmd/app/main.go

## To Do

* promhttp not collecting metrics for all routes (github.com/zsais/go-gin-prometheus?)
* gin access log to json, use LOG_LEVEL, GIN_MODE=release | gin.SetMode(gin.ReleaseMode)

