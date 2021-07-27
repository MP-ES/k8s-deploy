package infra

import (
	"bytes"
	"k8s-deploy/utils"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

func YqSearchQueryInFileWithStringSliceReturn(fileName string, query string) ([]string, error) {
	var out bytes.Buffer

	printer := yqlib.NewPrinter(&out, false, true, false, 0, false)
	streamEvaluator := yqlib.NewStreamEvaluator()

	if err := streamEvaluator.EvaluateFiles(query, []string{fileName}, printer, true); err != nil {
		return nil, err
	}

	returnedSlice := strings.Split(out.String(), "\n")
	// apply trim
	returnedSlice = utils.SliceMapStrFunction(returnedSlice, strings.TrimSpace)
	// remove empty
	returnedSlice = utils.SliceRemoveEmptyElements(returnedSlice)
	// remove duplicates
	returnedSlice = utils.SliceRemoveDuplicateElements(returnedSlice)

	return returnedSlice, nil
}
