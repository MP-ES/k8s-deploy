package entities

import (
	"bufio"
	"encoding/json"
	"fmt"
	"k8s-deploy/infra"
	"k8s-deploy/utils"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type Ingress struct {
	Name string
}

type ingressHosts struct {
	Name  string
	Hosts []string
	Tls   [][]string
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

func GetDeployedIngresses(finalDeployedPath string, repository *Repository, kEnv *K8sEnv) ([]string, error) {
	var globalErr *multierror.Error
	ingresses := []string{}

	if _, ok := (*repository.GitOpsRules.Ingresses)[kEnv.Name]; !ok {
		return ingresses, nil // don't have ingress
	}

	// get ingresses hosts
	search, err := infra.YqSearchQueryInFileWithJsonReturn(finalDeployedPath,
		"{\"name\": .metadata.name, \"kind\": .kind, \"hosts\": [.spec.rules[].host]} | select (.kind == \"Ingress\")")

	if err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		scanner := bufio.NewScanner(search)
		for scanner.Scan() {
			var ingHosts ingressHosts
			if err := json.Unmarshal(scanner.Bytes(), &ingHosts); err != nil {
				globalErr = multierror.Append(globalErr, err)
			} else {
				ingresses = append(ingresses, ingHosts.Hosts...)
			}
		}
	}

	return ingresses, globalErr.ErrorOrNil()

}

func GetIngressesHostReplace(appDeployPath string, repository *Repository, gitOpsRepo *GitOpsRepository,
	eventRef *eventRef, kEnv *K8sEnv) ([]*infra.IngressReplacement, error) {
	var globalErr *multierror.Error
	ingressReplacements := []*infra.IngressReplacement{}

	if _, ok := (*repository.GitOpsRules.Ingresses)[kEnv.Name]; eventRef.Type != utils.EventTypePullRequest || !ok {
		return nil, nil // don't have ingress or is not PR
	}

	// count ingresses to deploy
	qtyIngressHosts, err := countQtyIngressHosts(appDeployPath)
	if err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		// get ingresses hosts
		search, err := infra.YqSearchQueryInFileWithJsonReturn(appDeployPath,
			"{\"name\": .metadata.name, \"kind\": .kind, \"hosts\": [.spec.rules[].host]} | select (.kind == \"Ingress\")")

		if err != nil {
			globalErr = multierror.Append(globalErr, err)
		} else {
			scanner := bufio.NewScanner(search)
			for scanner.Scan() {
				var ingHosts ingressHosts
				if err := json.Unmarshal(scanner.Bytes(), &ingHosts); err != nil {
					globalErr = multierror.Append(globalErr, err)
				} else {
					for i, v := range ingHosts.Hosts {
						ingressReplacements = append(ingressReplacements, createIngressHostReplacement(ingHosts.Name, i,
							generatePrHostValue(v, qtyIngressHosts, repository.Name, eventRef.Identifier, gitOpsRepo.UrlPR)))
					}
				}
			}
		}

		// get ingresses tls
		search, err = infra.YqSearchQueryInFileWithJsonReturn(appDeployPath,
			"{\"name\": .metadata.name, \"kind\": .kind, \"tls\": [.spec.tls[].hosts]} | select (.kind == \"Ingress\")")

		if err != nil {
			globalErr = multierror.Append(globalErr, err)
		} else {
			scanner := bufio.NewScanner(search)
			for scanner.Scan() {
				var ingHosts ingressHosts
				if err := json.Unmarshal(scanner.Bytes(), &ingHosts); err != nil {
					globalErr = multierror.Append(globalErr, err)
				} else {
					for iTls, v := range ingHosts.Tls {
						for iHost, t := range v {
							ingressReplacements = append(ingressReplacements, createIngressTlsReplacement(ingHosts.Name, iHost, iTls,
								generatePrHostValue(t, qtyIngressHosts, repository.Name, eventRef.Identifier, gitOpsRepo.UrlPR)))
						}
					}
				}
			}
		}

	}

	return ingressReplacements, globalErr.ErrorOrNil()
}

func countQtyIngressHosts(appDeployPath string) (int, error) {
	ingresses, err := infra.YqSearchQueryInFileWithStringSliceReturn(appDeployPath,
		".spec.rules[].host")
	if err != nil {
		return 0, err
	}
	return len(ingresses), nil
}

func createIngressHostReplacement(ingressName string, hostIndex int, hostNewValue string) *infra.IngressReplacement {
	return &infra.IngressReplacement{IngressName: ingressName, HostIndex: hostIndex, HostNewValue: hostNewValue}
}

func createIngressTlsReplacement(ingressName string, hostIndex int, tlsIndex int, hostNewValue string) *infra.IngressReplacement {
	return &infra.IngressReplacement{IngressName: ingressName, IsTls: true, HostIndex: hostIndex, TlsIndex: tlsIndex, HostNewValue: hostNewValue}
}

func generatePrHostValue(currentValue string, qtyHosts int, repoName string, eventIdentifier string, urlPR string) string {
	var initialName = fmt.Sprintf("%s-pr%s", repoName, eventIdentifier)
	if qtyHosts == 1 {
		return fmt.Sprintf("%s%s", initialName, urlPR)
	}
	splitCurrent := strings.Split(currentValue, ".")
	return fmt.Sprintf("%s-%s%s", initialName, splitCurrent[0], urlPR)
}
