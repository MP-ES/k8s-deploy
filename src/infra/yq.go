package infra

import (
	"bytes"
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

	return strings.Split(out.String(), "\n"), nil
}
