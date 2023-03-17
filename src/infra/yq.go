package infra

import (
	"bytes"
	"fmt"
	"k8s-deploy/utils"
	"os/exec"
	"strings"
)

func YqSearchQueryInFileWithStringSliceReturn(fileName string, query string) ([]string, error) {
	var out *bytes.Buffer
	var err error

	if out, err = runYq(fileName, query, false); err != nil {
		return nil, err
	}

	returnedSlice := strings.Split((*out).String(), "\n")
	// apply trim
	returnedSlice = utils.SliceMapStrFunction(returnedSlice, strings.TrimSpace)
	// remove empty
	returnedSlice = utils.SliceRemoveEmptyElements(returnedSlice)
	// remove duplicates
	returnedSlice = utils.SliceRemoveDuplicateElements(returnedSlice)

	return returnedSlice, nil
}

func YqSearchQueryInFileWithJsonReturn(fileName string, query string) (*bytes.Buffer, error) {
	var out *bytes.Buffer
	var err error

	if out, err = runYq(fileName, query, true); err != nil {
		return nil, err
	}
	return out, nil
}

func runYq(fileName string, query string, outputToJSON bool) (*bytes.Buffer, error) {
	var outputParam string

	if outputToJSON {
		outputParam = "json"
	} else {
		outputParam = "yaml"
	}

	cmdRes, err := exec.Command("yq", "-I0", "-o="+outputParam, query, fileName).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error on run yq: %s", cmdRes)
	}

	return bytes.NewBuffer(cmdRes), nil
}
