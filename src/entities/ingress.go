package entities

import (
	"fmt"
	"k8s-deploy/infra"

	"github.com/hashicorp/go-multierror"
)

type Ingress struct {
	Name string
}

func (i *Ingress) String() string {
	return i.Name
}

func (i *Ingress) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return err
	}
	i.Name = output
	return nil
}

func ValidateIngressesFromAppDeploy(appDeployPath string, kEnv *K8sEnv, repoRules *RepositoryRules) error {
	var globalErr *multierror.Error

	ingresses, err := infra.YqSearchQueryInFileWithStringSliceReturn(appDeployPath,
		".spec.rules[].host,.spec.tls[].hosts[]")
	if err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		for _, ingress := range ingresses {
			if !repoRules.IsIngressEnabled(ingress, *kEnv) {
				globalErr = multierror.Append(globalErr,
					fmt.Errorf("ingress '%s' is not enabled in repository '%s' for K8S environment '%s'. Check the GitOps repository",
						ingress, repoRules.Name, kEnv.Name))
			}
		}
	}

	return globalErr.ErrorOrNil()
}
