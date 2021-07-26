package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
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
	// var deployenv entities.DeployEnv
	// var err error

	// // set logging
	// setLogging()

	// if deployenv, err = entities.GetDeployEnvironment(); err != nil {
	// 	githubactions.Fatalf(err.Error())
	// }
	// if err = deployenv.ValidateRules(); err != nil {
	// 	githubactions.Fatalf(err.Error())
	// }

	// output := render.Render(deployenv)
	// githubactions.SetOutput("status", output)

	for _, e := range os.Environ() {
		fmt.Println(e)
	}
}
