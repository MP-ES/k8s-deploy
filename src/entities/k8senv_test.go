package entities_test

import (
	"fmt"
	"k8s-deploy/entities"
	"k8s-deploy/utils"
	"os"
	"testing"

	"github.com/go-test/deep"
)

type k8sDeployEnvsTest struct {
	inputK8sEnvs    string
	expectedK8sEnvs []*entities.K8sEnv
	expectedError   string
}

var k8sDeployEnvsTests = [...]k8sDeployEnvsTest{
	{"", nil, "'k8s-env' is required"},
	{"wrong-env", nil, "kubernetes environment 'wrong-env' unknown"},
	{"test1\n", []*entities.K8sEnv{{Name: "test1", Kubeconfig: "test"}}, ""},
	{"test1\ntest4", []*entities.K8sEnv{{Name: "test1", Kubeconfig: "test"}}, "kubernetes environment 'test4' unknown"},
	{"test1\ntest3\ntest2", []*entities.K8sEnv{
		{Name: "test1", Kubeconfig: "test"},
		{Name: "test3", Kubeconfig: "test"},
		{Name: "test2", Kubeconfig: "test"}},
		""},
}

func TestGetK8sDeployEnvironments(t *testing.T) {
	availableK8sEnvs := map[string]struct{}{
		"test1": {},
		"test2": {},
		"test3": {},
	}

	for _, test := range k8sDeployEnvsTests {
		orig := os.Getenv("INPUT_K8S_ENVS")
		os.Setenv("INPUT_K8S_ENVS", test.inputK8sEnvs)
		for _, e := range test.expectedK8sEnvs {
			os.Setenv(fmt.Sprintf("base64_kubeconfig_%s", e.Name), "test")
		}
		t.Cleanup(func() { os.Setenv("INPUT_K8S_ENVS", orig) })

		k8sEnvs, err := entities.GetK8sDeployEnvironments(&availableK8sEnvs)

		if err != nil {
			if test.expectedError == "" || err.Error() != test.expectedError {
				t.Errorf("k8s envs error '%s' not equal to expected '%s'", err, test.expectedError)
			}
		} else {
			if diff := deep.Equal(k8sEnvs, test.expectedK8sEnvs); diff != nil {
				t.Errorf("k8s envs not equal to expected")
				t.Error(diff)
			}
		}
	}
}

type k8sEnvsTest struct {
	K8sEnvs       []*entities.K8sEnv
	eventType     string
	expectedError string
}

var k8sEnvsTests = [...]k8sEnvsTest{
	{[]*entities.K8sEnv{{Name: "test1"}}, utils.EventTypePullRequest, ""},
	{[]*entities.K8sEnv{{Name: "test1"}}, utils.EventTypeHead, ""},
	{[]*entities.K8sEnv{{Name: "test1"}, {Name: "test3"}, {Name: "test2"}}, utils.EventTypeHead, ""},
	{[]*entities.K8sEnv{{Name: "test1"}, {Name: "test3"}, {Name: "test2"}}, utils.EventTypePullRequest, "multiple K8s environments on pull request events are not allowed"},
}

func TestValidateK8sEnvs(t *testing.T) {
	for _, test := range k8sEnvsTests {
		err := entities.ValidateK8sEnvs(test.K8sEnvs, test.eventType)

		if (err == nil && test.expectedError != "") || (err != nil && test.expectedError == "") {
			t.Errorf("validate error '%v' not equal to expected '%s'", err, test.expectedError)
		}
	}
}
