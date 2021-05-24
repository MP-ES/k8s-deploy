package utils_test

import (
	"encoding/base64"
	"k8s-deploy/utils"
	"strings"
	"testing"
)

type unmarshalSingleKeyMultifileTest struct {
	base64FileContent string
	dataStruct        interface{}
	expectedData      string
	expectedError     string
}

type structDataTest struct {
	Data string
}

var structData = structDataTest{}
var unmarshalSingleKeyMultifileTests = [...]unmarshalSingleKeyMultifileTest{
	{"wrong-base64", nil, "", "illegal base64 data"},
	{base64.StdEncoding.EncodeToString([]byte("wrong-yaml")), &structData, "", "cannot unmarshal !!str `wrong-yaml`"},
	{base64.StdEncoding.EncodeToString([]byte("data: test-single-file")), &structData, "test-single-file", ""},
	{base64.StdEncoding.EncodeToString([]byte("---\n---\ndata: test-multifile")), &structData, "test-multifile", ""},
	{base64.StdEncoding.EncodeToString([]byte("---\n---\nno-data: test-nodata")), &structData, "", ""},
}

func TestUnmarshalSingleYamlKeyFromMultifile(t *testing.T) {
	for _, test := range unmarshalSingleKeyMultifileTests {
		err := utils.UnmarshalSingleYamlKeyFromMultifile(&test.base64FileContent, test.dataStruct)

		if err != nil {
			if test.expectedError == "" || !strings.Contains(err.Error(), test.expectedError) {
				t.Errorf("unmarshal single key error %s not equal to expected %s", err, test.expectedError)
			}
		} else {
			if structData.Data != test.expectedData {
				t.Errorf("unmarshal single key data %s not equal to expected %s", structData.Data, test.expectedData)
			}
		}

	}
}
