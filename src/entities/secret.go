package entities

import (
	"encoding/base64"
	"fmt"
	"k8s-deploy/infra"
	"os"

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

	// checking secrets name: must be the repository name
	secretsName, err := infra.YqSearchQueryInFileWithStringSliceReturn(appDeployPath,
		".spec.jobTemplate.spec.template.spec.containers[].env[].valueFrom.secretKeyRef.name,.spec.template.spec.containers[].env[].valueFrom.secretKeyRef.name,.spec.template.spec.volumes[].secret.secretName")
	if err != nil {
		globalErr = multierror.Append(globalErr, err)
	}
	if len(secretsName) > 1 {
		globalErr = multierror.Append(globalErr,
			fmt.Errorf("more than one k8s secret by repository is not allowed. current k8s-secrets: %v", secretsName))
	}
	if len(secretsName) > 0 && secretsName[0] != repoRules.Name {
		globalErr = multierror.Append(globalErr,
			fmt.Errorf("the k8s-secret name must be the same as the repository name. Current name: %s; expected: %s", secretsName[0], repoRules.Name))
	}

	// checking if all secrets was declared and was setted as env
	secrets, err := infra.YqSearchQueryInFileWithStringSliceReturn(appDeployPath,
		".spec.jobTemplate.spec.template.spec.containers[].env[].valueFrom.secretKeyRef.key,.spec.template.spec.containers[].env[].valueFrom.secretKeyRef.key,.spec.template.spec.volumes[].secret.items[].key")
	if err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		for _, secret := range secrets {
			if !repoRules.IsSecretEnabled(secret) {
				globalErr = multierror.Append(globalErr,
					fmt.Errorf("secret '%s' is not enabled in repository '%s'. Check the GitOps repository",
						secret, repoRules.Name))
			}

			if _, ok := os.LookupEnv(secret); !ok {
				globalErr = multierror.Append(globalErr,
					fmt.Errorf("secret '%s' is not setted as environment variable",
						secret))
			}
		}
	}

	return globalErr.ErrorOrNil()
}

func GetSecretsDeploy(secrets []*Secret) map[string]string {
	secretsList := map[string]string{}

	for _, secret := range secrets {
		secretsList[secret.Name] = base64.StdEncoding.EncodeToString([]byte(os.Getenv(secret.Name)))
	}

	return secretsList
}
