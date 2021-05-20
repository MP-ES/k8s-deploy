package main

import (
	"k8s-deploy/entities"

	"github.com/gdexlab/go-render/render"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	deployenv, err := entities.GetDeployEnvironment()
	if err != nil {
		githubactions.Fatalf(err.Error())
	}
	output := render.AsCode(deployenv)

	githubactions.SetOutput("test", output)
}
