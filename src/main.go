package main

import (
	"encoding/json"
	"k8s-deploy/entities"
	"k8s-deploy/infra"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-githubactions"
	"gopkg.in/op/go-logging.v1"
)

func setLogging() {
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
	)
	var backend = logging.AddModuleLevel(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		format))

	backend.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend)
}

func getStrErrorDeployment(results *[]entities.DeploymentResult) string {
	var sb strings.Builder

	for _, dr := range *results {
		if !dr.Deployed {
			sb.WriteString("Deployment error on k8s-env '")
			sb.WriteString(dr.K8sEnv)
			sb.WriteString("':\nError message:\n")
			sb.WriteString(dr.ErrMsg)
			sb.WriteString("\nDeployment Log:\n")
			sb.WriteString(dr.DeploymentLog)
		}
	}
	return sb.String()
}

func main() {
	var deployenv entities.DeployEnv
	var deploymentResultByte []byte
	var err error

	var dryRunMode, _ = strconv.ParseBool(os.Getenv("DRY_RUN_MODE"))

	if !dryRunMode {
		// destroy the deployment directory to avoid keep sensitive data
		defer infra.ClearDeploy()
	} else {
		println("Dry run mode enabled, skipping deployment directory cleanup.")
	}

	// set logging
	setLogging()

	if deployenv, err = entities.GetDeployEnvironment(); err != nil {
		githubactions.Fatalf("Error getting deploy environment: %v", err)
	}

	if err = deployenv.ValidateRules(); err != nil {
		githubactions.Fatalf("Error validating rules: %v", err)
	}

	deploymentResult := deployenv.Apply(dryRunMode)
	if deploymentResultByte, err = json.Marshal(deploymentResult); err != nil {
		githubactions.Fatalf("Error marshaling deployment result: %v", err)
	}

	if !dryRunMode {
		if err = deployenv.PostApplyActions(&deploymentResult); err != nil {
			githubactions.Warningf("Warning during post-apply actions: %v", err)
		}

		githubactions.SetOutput("status", string(deploymentResultByte))
	} else {
		println("Dry run mode enabled, skipping post-apply actions.")
	}

	errorsStr := getStrErrorDeployment(&deploymentResult)
	if errorsStr != "" {
		githubactions.Fatalf("Deployment errors: %s", errorsStr)
	}
}
