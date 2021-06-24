package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

type K8sEnv struct {
	Name string
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

func getK8sEnv(availableK8sEnvs *map[string]struct{}, s string) (*K8sEnv, error) {
	if _, ok := (*availableK8sEnvs)[s]; !ok {
		return nil, fmt.Errorf("kubernetes environment '%s' unknown", s)
	}
	return &K8sEnv{Name: s}, nil
}

func GetK8sDeployEnvironments(availableK8sEnvs *map[string]struct{}) ([]*K8sEnv, error) {
	k8sEnvs := []*K8sEnv{}

	k8sEnvsInput := githubactions.GetInput("k8s-envs")
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
