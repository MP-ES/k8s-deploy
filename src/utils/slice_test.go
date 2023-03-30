package utils_test

import (
	"k8s-deploy/utils"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

type sliceFunctionsTest struct {
	inputSlice          []string
	inputElement        string
	functionMap         func(string) string
	expectedOutputSlice []string
}

var sliceFunctionMapTests = [...]sliceFunctionsTest{
	{nil, "", nil, nil},
	{nil, "", strings.TrimSpace, nil},
	{[]string{}, "", strings.TrimSpace, []string{}},
	{[]string{"string"}, "", nil, []string{"string"}},
	{
		[]string{"   initial spaces", "end spaces   ", "   a lot of spaces   ", ""},
		"",
		strings.TrimSpace,
		[]string{"initial spaces", "end spaces", "a lot of spaces", ""},
	},
	{
		[]string{"lower_string", "lowerstring"},
		"",
		strings.ToUpper,
		[]string{"LOWER_STRING", "LOWERSTRING"},
	},
}

func TestSliceMapStrFunction(t *testing.T) {
	for _, test := range sliceFunctionMapTests {
		slice := utils.SliceMapStrFunction(test.inputSlice, test.functionMap)

		if diff := deep.Equal(slice, test.expectedOutputSlice); diff != nil {
			t.Errorf("returned slice not equal to expected")
			t.Error(diff)
		}
	}
}

var sliceFunctionRemoveEmptyTests = [...]sliceFunctionsTest{
	{nil, "", nil, nil},
	{[]string{}, "", nil, []string{}},
	{[]string{"validContent"}, "", nil, []string{"validContent"}},
	{
		[]string{"", "validContent", "   ", " ", "", "null"},
		"",
		nil,
		[]string{"validContent"},
	},
}

func TestSliceRemoveEmptyElements(t *testing.T) {
	for _, test := range sliceFunctionRemoveEmptyTests {
		slice := utils.SliceRemoveEmptyElements(test.inputSlice)

		if diff := deep.Equal(slice, test.expectedOutputSlice); diff != nil {
			t.Errorf("returned slice not equal to expected")
			t.Error(diff)
		}
	}
}

var sliceFunctionRemoveElementTests = [...]sliceFunctionsTest{
	{nil, "", nil, nil},
	{[]string{"content"}, "", nil, []string{"content"}},
	{[]string{"content"}, "content", nil, []string{}},
	{[]string{"content", "content2"}, "content", nil, []string{"content2"}},
	{[]string{"content", "content2", "content2"}, "content2", nil, []string{"content"}},
	{[]string{"content", "content2"}, "", nil, []string{"content", "content2"}},
}

func TestSliceRemoveElement(t *testing.T) {
	for _, test := range sliceFunctionRemoveElementTests {
		slice := utils.SliceRemoveElement(test.inputSlice, test.inputElement)

		if diff := deep.Equal(slice, test.expectedOutputSlice); diff != nil {
			t.Errorf("returned slice not equal to expected")
			t.Error(diff)
		}
	}
}

var sliceFunctionRemoveDuplicateTests = [...]sliceFunctionsTest{
	{nil, "", nil, nil},
	{[]string{}, "", nil, []string{}},
	{[]string{"validContent"}, "", nil, []string{"validContent"}},
	{
		[]string{"repeated", "validContent", "repeated"},
		"",
		nil,
		[]string{"repeated", "validContent"},
	},
	{
		[]string{"", "validContent", "", "rep1", "rep1"},
		"",
		nil,
		[]string{"", "validContent", "rep1"},
	},
}

func TestSliceRemoveDuplicateElements(t *testing.T) {
	for _, test := range sliceFunctionRemoveDuplicateTests {
		slice := utils.SliceRemoveDuplicateElements(test.inputSlice)

		if diff := deep.Equal(slice, test.expectedOutputSlice); diff != nil {
			t.Errorf("returned slice not equal to expected")
			t.Error(diff)
		}
	}
}
