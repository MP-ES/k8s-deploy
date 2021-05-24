package entities_test

import (
	"k8s-deploy/entities"
	"os"
	"reflect"
	"testing"
)

type k8sDeployEnvsTest struct {
	inputK8sEnvs    string
	expectedK8sEnvs []*entities.K8sEnv
	expectedError   string
}

var k8sDeployEnvsTests = [...]k8sDeployEnvsTest{
	{"", nil, "'k8s-env' is required"},
	{"wrong-env", nil, "kubernetes environment 'wrong-env' unknown"},
	{"test1\n", []*entities.K8sEnv{{Name: "test1"}}, ""},
	{"test1\ntest4", []*entities.K8sEnv{{Name: "test1"}}, "kubernetes environment 'test4' unknown"},
	{"test1\ntest3\ntest2", []*entities.K8sEnv{{Name: "test1"}, {Name: "test3"}, {Name: "test2"}}, ""},
}

func TestGetK8sDeployEnvironments(t *testing.T) {
	availableK8sEnvs := map[string]struct{}{
		"test1": {},
		"test2": {},
		"test3": {},
	}

	for _, test := range k8sDeployEnvsTests {
		orig := os.Getenv("INPUT_K8S-ENVS")
		os.Setenv("INPUT_K8S-ENVS", test.inputK8sEnvs)
		t.Cleanup(func() { os.Setenv("INPUT_K8S-ENVS", orig) })

		k8sEnvs, err := entities.GetK8sDeployEnvironments(&availableK8sEnvs)

		if err != nil {
			if test.expectedError == "" || err.Error() != test.expectedError {
				t.Errorf("k8s envs error %s not equal to expected %s", err, test.expectedError)
			}
		} else {
			if !reflect.DeepEqual(k8sEnvs, test.expectedK8sEnvs) {
				t.Errorf("k8s envs %v not equal to expected %v", k8sEnvs, test.expectedK8sEnvs)
			}
		}
	}
}
