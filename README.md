# go_docker_api
Go Docker Kubernetes Neo4j MongoDB API
## RUN 
cmd/api/main.go

## Authors
- Jose Manuel Broto - jmbroto@cimne.upc.edu
- Gerard Mor - gmor@cimne.upc.edu
- Andreu Marí - amari@cimne.upc.edu

Copyright (c) 2024 Jose Manuel Broto, Gerard Mor, Andreu Marí

[//]: # (docker buildx build --platform linux/amd64,linux/arm64 --push -t docker.tech.beegroup-cimne.com/crbcn_api/go-api . --provenance=false)
[//]: # (kubectl apply -f 'C:\Proyectos\BEE\go-api\go-api-service.yaml'
[//]: # (kubectl apply -f 'C:\Proyectos\BEE\go-api\go-api-deployment.yaml')
[//]: # (kubectl create secret generic crbcn-api-secret --from-file=.env --namespace data-api)

