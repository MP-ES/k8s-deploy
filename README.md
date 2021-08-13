# k8s-deploy

Action that deploys an application in an On-Premises Kubernetes cluster based in a GitOps repository.

[![Coverage](https://codecov.io/gh/MP-ES/k8s-deploy/branch/main/graphs/badge.svg?branch=main)](https://codecov.io/gh/MP-ES/k8s-deploy)
[![Integration](https://github.com/MP-ES/k8s-deploy/workflows/Integration/badge.svg)](https://github.com/MP-ES/k8s-deploy/actions?query=workflow%3AIntegration)

## Requirements

The owner must have a repository named **gitops** with the rules of application deployment. For example, if you are deploying the repository **ORG/application**, then this k8s-deploy will try to get the rules in the repository **ORG/gitops**, once the repository owner is **ORG**.

## Usage

```yaml
- name: Deploy on on-premises K8S
  uses: MP-ES/k8s-deploy@main
  with:
    # Multiline input where each line contains the name of a Kubernetes environment defined in the GitOps repository
    k8s-envs: |
      env1
      env2

    # Path to the manifest directory, with files to be used for deployment
    # DEFAULT: kubernetes
    manifest-dir: kubernetes

    # GitHub PAT with read permission on gitOps repository, if gitOps is private
    gitops-token: ${{ secrets.SECRET_NAME }}
```

## Outputs

Following outputs are available:

| Name     | Type        | Description                                   |
| -------- | ----------- | --------------------------------------------- |
| `status` | JSON object | Array of deployment status by K8S environment |

Output example:

```json
[
   {
      "K8sEnv":"dev",
      "Deployed":true,
      "ErrMsg":"",
      "Ingresses":[
         "inova-pr6.dev.mpes.mp.br"
      ]
   },
   {
      "K8sEnv":"app",
      "Deployed":false,
      "ErrMsg":"Server unavailable",
      "Ingresses":[]
   }
]
```

## Developer

```shell
# Copy .env.* example file to .env file
# Simulate a pull request call
cp src/.env.pr src/.env

# Install lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -c 'sh -s -- -b /usr/local/bin'

# Run lint locally
# From src directory
golangci-lint run

# Run tests
# From src directory
go test -race -v -covermode=atomic -coverprofile=coverage.out ./...

# See cover report
# From src directory
go tool cover -html=coverage.out
```
