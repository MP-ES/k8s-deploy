package main

import (
	"fmt"
	"k8s-deploy/entities"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	deployenv := entities.GetDeployEnvironment()
	output := fmt.Sprintf("%+v\n", deployenv)

	fmt.Println(output)
	githubactions.SetOutput("test", output)
}
