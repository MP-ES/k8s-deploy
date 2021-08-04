package infra

import (
	"bytes"
	"k8s-deploy/utils"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
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
	var out bytes.Buffer

	printer := yqlib.NewPrinter(&out, outputToJSON, true, false, 0, false)
	streamEvaluator := yqlib.NewStreamEvaluator()

	if err := streamEvaluator.EvaluateFiles(query, []string{fileName}, printer, true); err != nil {
		return nil, err
	}
	return &out, nil
}
