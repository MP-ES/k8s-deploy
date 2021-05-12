# k8s-deploy

Action that deploys an application in an On-Premises Kubernetes cluster based in a GitOps repository.

[![Integration](https://github.com/MP-ES/k8s-deploy/workflows/Integration/badge.svg)](https://github.com/MP-ES/k8s-deploy/actions?query=workflow%3AIntegration)

## Requirements

The owner must have a repository named **gitops** with the rules of application deployment. For example, if you are deploying the repository **ORG/application**, then this k8s-deploy will try to get the rules in the repository **ORG/gitops**, once the repository owner is **ORG**.

## Usage

```yaml
- name: Deploy on on-premises K8S
  uses: MP-ES/k8s-deploy@main
  with:
    # Multiline input where each line contains the name of a Kubernetes environment defined in the GitOps repository.
    k8s-envs: |
      env1
      env2

    # Path to the manifest directory, with files to be used for deployment.
    # DEFAULT: kubernetes
    manifest-dir: kubernetes
```

## Developer

```shell
# Copy .env.* example file to .env file
# Simulate a pull request call
cp src/.env.pr src/.env

# Install lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -c 'sh -s -- -b /usr/local/bin'

# Run lint locally
pushd src/ && golangci-lint run && popd
```
