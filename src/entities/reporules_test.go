package entities_test

import (
	"k8s-deploy/entities"
	"testing"
)

type k8sEnvEnabledTest struct {
	kEnvs          []*entities.K8sEnv
	kEnv           *entities.K8sEnv
	expectedResult bool
}

var k8sEnvEnabledTests = [...]k8sEnvEnabledTest{
	{[]*entities.K8sEnv{}, &entities.K8sEnv{Name: "env1"}, false},
	{[]*entities.K8sEnv{{Name: "env1"}}, &entities.K8sEnv{Name: "env1"}, true},
	{[]*entities.K8sEnv{{Name: "env1"}}, &entities.K8sEnv{Name: "env2"}, false},
	{[]*entities.K8sEnv{{Name: "env1"}, {Name: "env2"}}, &entities.K8sEnv{Name: "env2"}, true},
}

func TestIsK8sEnvEnabled(t *testing.T) {
	for _, test := range k8sEnvEnabledTests {
		repoRules := entities.RepositoryRules{K8sEnvs: test.kEnvs}

		res := repoRules.IsK8sEnvEnabled(test.kEnv)

		if res != test.expectedResult {
			t.Errorf("enabled test for k8s environment '%s' not equal to expected '%t'", test.kEnv, res)
		}
	}
}
