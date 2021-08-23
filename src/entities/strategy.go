package entities

import (
	"fmt"
	"k8s-deploy/infra"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

const defaultStrategy = "none"

func getAvailableStrategies() *map[string]struct{} {
	return &map[string]struct{}{"none": {}, "canary": {}, "blue-green": {}}
}

type Strategy struct {
	Name string
}

func (s *Strategy) String() string {
	return s.Name
}

func (s *Strategy) Deploy(kEnv *K8sEnv, finalDeployedPath string) (string, error) {
	if s.Name == "none" {
		return infra.KubectlApply(kEnv.Name, finalDeployedPath)
	}
	return "", fmt.Errorf("Strategy '%s' doesn't have a deployment plan", s.Name)
}

func GetStrategy() (*Strategy, error) {
	strategyInput := githubactions.GetInput("strategy")
	if strategyInput == "" {
		strategyInput = defaultStrategy
	}

	if _, ok := (*getAvailableStrategies())[strategyInput]; !ok {
		return nil, fmt.Errorf("Strategy '%s' unknown. Available strategies: %s", strategyInput, getStrAvailableStrategies())
	}

	// not implemented yet
	if strategyInput != "none" {
		return nil, fmt.Errorf("Strategy '%s' not implemented yet", strategyInput)
	}

	return &Strategy{Name: strategyInput}, nil
}

func getStrAvailableStrategies() string {
	str := ""
	for s := range *getAvailableStrategies() {
		str += s + ", "
	}
	return strings.Trim(str, ", ")
}
