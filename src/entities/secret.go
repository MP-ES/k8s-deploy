package entities

import (
	"fmt"
	"k8s-deploy/infra"

	"github.com/hashicorp/go-multierror"
)

type Secret struct {
	Name string
}

func (s *Secret) String() string {
	return s.Name
}

func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return err
	}
	s.Name = output
	return nil
}

func ValidateSecretsFromAppDeploy(appDeployPath string, repoRules *RepositoryRules) error {
	var globalErr *multierror.Error
	var err error

	secrets, err := infra.YqSearchQueryInFileWithStringSliceReturn(appDeployPath,
		".spec.jobTemplate.spec.template.spec.containers[].env[].name,.spec.template.spec.containers[].env[].name")
	if err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		for _, secret := range secrets {
			if !repoRules.IsSecretEnabled(secret) {
				globalErr = multierror.Append(globalErr,
					fmt.Errorf("secret '%s' is not enabled in repository '%s'. Check the GitOps repository",
						secret, repoRules.Name))
			}
		}
	}

	return globalErr.ErrorOrNil()
}
