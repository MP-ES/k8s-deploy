package infra_test

import (
	"k8s-deploy/infra"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

type searchInFileSliceTest struct {
	fileName      string
	query         string
	expectedSlice []string
	expectedError string
}

var searchInFileSliceTests = [...]searchInFileSliceTest{
	{"fileNotFound.yaml", "", nil, "no such file or directory"},
	{"../../testdata/repository-min.yaml", "wrong-query", nil, "error on run yq: Error: 1:1: invalid input text \"wrong-query\""},
	{"../../testdata/repository-min.yaml", "", []string{
		"name: repository-min",
		"k8s-envs:",
		"- env1",
		"images:",
		"- docker_image",
		"resources-quotas:",
		"limits.cpu: 100m",
		"limits.memory: 100Mi"},
		""},
	{"../../testdata/repository-min.yaml", ".images[]", []string{"docker_image"}, ""},
	{"../../testdata/repository-min.yaml", ".resources-quotas", []string{"limits.cpu: 100m", "limits.memory: 100Mi"}, ""},
}

func TestYqSearchQueryInFileWithStringSliceReturn(t *testing.T) {
	for _, test := range searchInFileSliceTests {
		slice, err := infra.YqSearchQueryInFileWithStringSliceReturn(test.fileName, test.query)

		if err != nil {
			if test.expectedError == "" || !strings.Contains(err.Error(), test.expectedError) {
				t.Errorf("Yq search with slice return error '%v' not equal to expected '%s'", err, test.expectedError)
			}
		} else {
			if diff := deep.Equal(slice, test.expectedSlice); diff != nil {
				t.Errorf("returned slice not equal to expected")
				t.Error(diff)
			}
		}
	}
}

type searchInFileJsonTest struct {
	fileName      string
	query         string
	expectedJson  string
	expectedError string
}

var searchInFileJsonTests = [...]searchInFileJsonTest{
	{"fileNotFound.yaml", "", "", "no such file or directory"},
	{"../../testdata/repository-min.yaml", "wrong-query", "", "error on run yq: Error: 1:1: invalid input text \"wrong-query\""},
	{"../../testdata/repository-min.yaml", "",
		"{\"name\":\"repository-min\",\"k8s-envs\":[\"env1\"],\"images\":[\"docker_image\"],\"resources-quotas\":{\"limits.cpu\":\"100m\",\"limits.memory\":\"100Mi\"}}\n",
		""},
	{"../../testdata/repository-min.yaml", ".images[]", "\"docker_image\"\n", ""},
	{"../../testdata/repository-min.yaml", ".resources-quotas", "{\"limits.cpu\":\"100m\",\"limits.memory\":\"100Mi\"}\n", ""},
}

func TestYqSearchQueryInFileWithJsonReturn(t *testing.T) {
	for _, test := range searchInFileJsonTests {
		bytes, err := infra.YqSearchQueryInFileWithJsonReturn(test.fileName, test.query)

		if err != nil {
			if test.expectedError == "" || !strings.Contains(err.Error(), test.expectedError) {
				t.Errorf("Yq search with slice return error '%v' not equal to expected '%s'", err, test.expectedError)
			}
		} else {
			if diff := deep.Equal(bytes.String(), test.expectedJson); diff != nil {
				t.Errorf("returned JSON not equal to expected")
				t.Error(diff)
			}
		}
	}
}
