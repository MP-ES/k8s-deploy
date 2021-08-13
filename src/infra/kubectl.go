package infra

import "fmt"

func KubectlApply(finalDeployedPath string) error {
	fmt.Println(finalDeployedPath)
	return nil
}
