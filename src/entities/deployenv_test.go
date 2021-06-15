package entities_test

import (
	"k8s-deploy/entities"
	"testing"
)

func TestIfErrorOnGetGitOpsRepoCausePrematureFailOnGetDeployment(t *testing.T) {
	_, err := entities.GetDeployEnvironment()

	if err == nil {
		t.Errorf("error on getting gitOps repository doesn't cause premature failure on getting deployment")
	}
}
