package utils_test

import (
	"k8s-deploy/utils"
	"testing"
)

type eventRefTest struct {
	githubRef, expectedType, expectedId string
	expectError                         string
}

var eventRefTests = [...]eventRefTest{
	{"refs/heads/main", "heads", "main", ""},
	{"refs/tags/v1.0.0", "tags", "v1.0.0", ""},
	{"refs/pull/1/merge", "pull", "1", ""},
	{"wrong/string", "", "", "unknown Github reference"},
}

func TestGetGithubEventRef(t *testing.T) {

	for _, test := range eventRefTests {
		eventType, eventId, err := utils.GetGithubEventRef(test.githubRef)

		if eventType != test.expectedType {
			t.Errorf("event type %s not equal to expected %s", eventType, test.expectedType)
		}

		if eventId != test.expectedId {
			t.Errorf("event identifier %s not equal to expected %s", eventId, test.expectedId)
		}

		if (err != nil && test.expectError == "") || (err != nil && err.Error() != test.expectError) {
			t.Errorf("event error %s not equal to expected %s", err, test.expectError)
		}
	}
}
