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

type imageEnabledTest struct {
	images         []*entities.Image
	imageName      string
	expectedResult bool
}

var imageEnabledTests = [...]imageEnabledTest{
	{[]*entities.Image{}, "image", false},
	{[]*entities.Image{{Name: "image"}}, "image", true},
	{[]*entities.Image{{Name: "image"}}, "image2", false},
	{[]*entities.Image{{Name: "image1"}, {Name: "image2"}}, "image2", true},
}

func TestIsImageEnabled(t *testing.T) {
	for _, test := range imageEnabledTests {
		repoRules := entities.RepositoryRules{Images: test.images}

		res := repoRules.IsImageEnabled(test.imageName)

		if res != test.expectedResult {
			t.Errorf("enabled test for image '%s' not equal to expected '%t'", test.imageName, res)
		}
	}
}

type secretEnabledTest struct {
	secrets        []*entities.Secret
	secretName     string
	expectedResult bool
}

var secretEnabledTests = [...]secretEnabledTest{
	{[]*entities.Secret{}, "secret", false},
	{[]*entities.Secret{{Name: "secret"}}, "secret", true},
	{[]*entities.Secret{{Name: "secret"}}, "secret2", false},
	{[]*entities.Secret{{Name: "secret1"}, {Name: "secret2"}}, "secret2", true},
}

func TestIsSecretEnabled(t *testing.T) {
	for _, test := range secretEnabledTests {
		repoRules := entities.RepositoryRules{Secrets: test.secrets}

		res := repoRules.IsSecretEnabled(test.secretName)

		if res != test.expectedResult {
			t.Errorf("enabled test for secret '%s' not equal to expected '%t'", test.secretName, res)
		}
	}
}
