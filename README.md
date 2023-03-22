# k8s-deploy

Action that deploys an application in an On-Premises Kubernetes cluster based in a GitOps repository.

[![Coverage](https://codecov.io/gh/MP-ES/k8s-deploy/branch/main/graphs/badge.svg?branch=main)](https://codecov.io/gh/MP-ES/k8s-deploy)
[![Integration](https://github.com/MP-ES/k8s-deploy/workflows/Integration/badge.svg)](https://github.com/MP-ES/k8s-deploy/actions?query=workflow%3AIntegration)
[![Release](https://github.com/MP-ES/k8s-deploy/workflows/Release/badge.svg)](https://github.com/MP-ES/k8s-deploy/actions?query=workflow%3ARelease)

## Requirements

The owner must have a repository named **gitops** with the rules of application deployment. For example, if you are deploying the repository **ORG/application**, then this k8s-deploy will try to get the rules in the repository **ORG/gitops**, once the repository owner is **ORG**.

## Usage

```yaml
- name: Deploy on on-premises K8S
  uses: MP-ES/k8s-deploy@v2
  with:
    # Multiline input where each line contains the name of a Kubernetes environment defined in the GitOps repository
    k8s_envs: |
      env1
      env2

    # Path to the manifest directory, with files to be used for deployment
    # DEFAULT: kubernetes
    manifest_dir: kubernetes

    # Personal access token (PAT) used to manage comments on pull request
    # DEFAULT: ${{ github.token }}
    repo_token: ${{ github.token }}

    # GitHub PAT with read permission on gitOps repository, if gitOps is private
    gitops_token: ${{ secrets.SECRET_NAME }}

    # Deployment strategy to be used. Allowed values are none, canary and blue-green
    # More details below
    # DEFAULT: none
    strategy: none

  env:
    # list of app secrets, defined in gitOps repository
    app_secret1: ${{ secrets.app_secret1 }}
    app_secret2: ${{ secrets.app_secret2 }}

    # base64 of kubeconfig file for each Kubernetes environment defined in k8s_envs
    # See below an example of an expected kubeconfig
    base64_kubeconfig_env1: ${{ secrets.base64_kubeconfig_env1 }}
    base64_kubeconfig_env2: ${{ secrets.base64_kubeconfig_env2 }}
```

### Strategy

Deployment strategy to be used while applying manifest files on the cluster. Acceptable values are none, canary and blue-green.

#### none

No deployment strategy is used when deploying. The files are changed on the cluster in force mode. This is sufficient to pull requests deployments or if the application can have short downtime during deployment.

#### canary

*not implemented yet.*

#### blue-green

*not implemented yet.*

### kubeconfig example

The most important part is the **context name**, which **must be** the same as the **Kubernetes environment name** to which the kubeconfig belongs.

```yaml
apiVersion: v1
kind: Config
clusters:
  - cluster:
      certificate-authority-data: base64-encoded of ca-file
      server: https://server.domain.com:6443
    name: k8s-cluster
users:
  - name: kube-admin-user
    user:
      client-certificate-data: base64-encoded of cert-file
      client-key-data: base64-encoded of key-file
contexts:
  - context:
      cluster: k8s-cluster
      user: kube-admin-user
    name: env1
```

You can generate a base64 of the file with `base64 -w 0 kubeconfig_file.yaml`.

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
      "DeploymentLog":"deployment.apps/test created\nservice/test created\ningress.extensions/test created\nnamespace/test unchanged\nresourcequota/test unchanged\nsecret/test unchanged\n",
      "Ingresses":[
         "ingress.env.domain.com"
      ]
   },
   {
      "K8sEnv":"app",
      "Deployed":false,
      "ErrMsg":"1 error occurred:\n\t* exit status 1\n\n",
      "DeploymentLog":"resourcequota/test created\nsecret/test created\nError from server (NotFound): error when creating \"../.deploy/pr/final.yaml\": namespaces \"test\" not found\n",
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
