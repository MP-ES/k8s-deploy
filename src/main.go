package main

import (
	"encoding/json"
	"k8s-deploy/entities"
	"os"

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

func main() {
	var deployenv entities.DeployEnv
	var deploymentResult []byte
	var err error

	// destroy the deployment directory to avoid keep sensitive data
	// defer infra.ClearDeploy()

	// set logging
	setLogging()

	if deployenv, err = entities.GetDeployEnvironment(); err != nil {
		githubactions.Fatalf(err.Error())
	}
	if err = deployenv.ValidateRules(); err != nil {
		githubactions.Fatalf(err.Error())
	}
	if deploymentResult, err = json.Marshal(deployenv.Apply()); err != nil {
		githubactions.Fatalf(err.Error())
	}

	githubactions.SetOutput("status", string(deploymentResult))
}
