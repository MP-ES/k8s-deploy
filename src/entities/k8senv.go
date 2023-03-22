package entities

import (
	"encoding/base64"
	"errors"
	"fmt"
	"k8s-deploy/infra"
	"k8s-deploy/utils"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

type K8sEnv struct {
	Name       string
	Kubeconfig string
}

func (k *K8sEnv) String() string {
	return k.Name
}

func (k *K8sEnv) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return err
	}
	k.Name = output
	return nil
}

func (k *K8sEnv) IsValidToRepository(gitOpsRepo *GitOpsRepository, repoRules *RepositoryRules, event *eventRef) error {
	if !repoRules.IsK8sEnvEnabled(k) {
		return fmt.Errorf("k8s-env '%s' is not enabled in repository '%s'. Check the GitOps repository", k.Name, repoRules.Name)
	}

	// check pull request event
	if event.Type == utils.EventTypePullRequest {
		if _, ok := gitOpsRepo.AvailableK8sEnvsToPR[k.Name]; !ok {
			return fmt.Errorf("k8s-env '%s' is not enabled in repository '%s' on pull request events", k.Name, repoRules.Name)
		}
	}

	return nil
}

func (k *K8sEnv) ValidateKubeconfig() error {
	content, err := base64.StdEncoding.DecodeString(k.Kubeconfig)
	if err != nil {
		return fmt.Errorf("wrong kubeconfig data format: %s", err.Error())
	}

	if res := infra.CreateKubeconfigFile(k.Name, content); !res {
		return errors.New("error when try create kubeconfig file")
	}

	if err := infra.KubectlCheckClusterConnection(k.Name); err != nil {
		return err
	}

	return nil
}

func getK8sEnv(availableK8sEnvs *map[string]struct{}, s string) (*K8sEnv, error) {
	if _, ok := (*availableK8sEnvs)[s]; !ok {
		return nil, fmt.Errorf("kubernetes environment '%s' unknown", s)
	}

	kubeconfig := os.Getenv(fmt.Sprintf("base64_kubeconfig_%s", s))
	if kubeconfig == "" {
		return nil, fmt.Errorf("kubeconfig of k8s-env '%s' not set (expected value in base64_kubeconfig_%s environment variable)", s, s)
	}

	return &K8sEnv{Name: s, Kubeconfig: kubeconfig}, nil
}

func GetK8sDeployEnvironments(availableK8sEnvs *map[string]struct{}) ([]*K8sEnv, error) {
	k8sEnvs := []*K8sEnv{}

	k8sEnvsInput := githubactions.GetInput("k8s_envs")
	if k8sEnvsInput == "" {
		return nil, errors.New("'k8s-env' is required")
	}

	envs := strings.Split(k8sEnvsInput, "\n")
	for _, e := range envs {
		if kEvent, err := getK8sEnv(availableK8sEnvs, e); err == nil {
			k8sEnvs = append(k8sEnvs, kEvent)
		} else {
			return nil, err
		}
	}

	return k8sEnvs, nil
}

func ValidateK8sEnvs(K8sEnvs []*K8sEnv, eventType string) error {
	if eventType == utils.EventTypePullRequest {
		if len(K8sEnvs) > 1 {
			return fmt.Errorf("multiple K8s environments on pull request events are not allowed")
		}
	}
	return nil
}
