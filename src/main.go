package main

import (
	"k8s-deploy/entities"

	"github.com/gdexlab/go-render/render"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	var deployenv entities.DeployEnv
	var err error

	if deployenv, err = entities.GetDeployEnvironment(); err != nil {
		githubactions.Fatalf(err.Error())
	}
	if err = deployenv.ValidateRules(); err != nil {
		githubactions.Fatalf(err.Error())
	}

	output := render.Render(deployenv)
	githubactions.SetOutput("status", output)
}
