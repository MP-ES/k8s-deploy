package main

import (
	"fmt"
	"io/ioutil"
	"k8s-deploy/entities"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	deployenv := entities.GetDeployEnvironment()
	output := fmt.Sprintf("%+v\n", deployenv)

	fmt.Println(output)
	githubactions.SetOutput("test", output)

	file, err := os.Open(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	fmt.Print(string(b))
}
